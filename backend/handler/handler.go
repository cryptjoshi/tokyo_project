package handler

import (
	"context"
	"fmt"
	"github.com/amalfra/etag"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"
	"github.com/valyala/fasthttp"
	_ "github.com/go-sql-driver/mysql"
	"hanoi/models"
	"gorm.io/gorm"
	"math/rand"
	//"github.com/golang-jwt/jwt"
	jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	"hanoi/repository"
	"hanoi/database"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	//"strings"
)

var Words = []string{
	"apple", "banana", "cherry", "date", "elderberry",
	"fig", "grape", "honeydew", "kiwi", "lemon",
	"mango", "orange", "papaya", "quince", "raspberry",
	"strawberry", "tangerine", "watermelon", "blueberry", "blackberry",
	"apricot", "cranberry", "pineapple", "pomegranate", "grapefruit",
	"avocado", "coconut", "guava", "lime", "passionfruit",
	"lychee", "nectarine", "plum", "apricot", "kiwifruit",
	"boysenberry", "cantaloupe", "rambutan", "starfruit", "persimmon",
	"currant", "dragonfruit", "gooseberry", "papaya", "ugli fruit",
	"quince", "ackee", "durian", "jackfruit", "kumquat",
	"litchi", "mulberry", "olive", "rhubarb", "tamarind",
	"tomato", "coconut", "cucumber", "eggplant", "zucchini",
	"potato", "carrot", "onion", "garlic", "broccoli",
	"cauliflower", "spinach", "kale", "lettuce", "cabbage",
	"brussels sprouts", "artichoke", "asparagus", "celery", "green bean",
	"peas", "corn", "radish", "beet", "turnip",
	"rutabaga", "pars"}


var ctx = context.Background()
var amqp_connection *amqp.Connection
var amqp_channel *amqp.Channel
var queue amqp.Queue = amqp.Queue{}
var is_connection = false
var has_channel = false
var has_queue = false

var redis_master_host = os.Getenv("REDIS_HOST")
var redis_master_port = os.Getenv("REDIS_PORT")
var redis_master_password = os.Getenv("REDIS_PASSWORD")
var redis_slave_host = os.Getenv("REDIS_SLAVE_HOST")
var redis_slave_port = os.Getenv("REDIS_SLAVE_PORT")
var redis_slave_password = os.Getenv("REDIS_SLAVE_PASSWORD")
var queue_name = "wallet" //os.Getenv("QUEUE_NAME")
var exchange_name = "wallet" //os.Getenv("EXCHANGE_NAME")
var rabbit_mq = "amqp://128.199.92.45:5672" //os.Getenv("RABBIT_MQ") @rabbitmq:5672/wallet
var connection_timeout = os.Getenv("CONNECTION_TIMEOUT")
var redis_database = getEnv("REDIS_DATABASE", "0")
var go_pixel_log = os.Getenv("GO_PIXEL_LOG")

var jwtSecret = os.Getenv("PASSWORD_SECRET")


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func InitAMQP() {
	fmt.Println("Init AMQP RABBIT")
	fmt.Println("channel")
	fmt.Println(amqp_channel)
	fmt.Println(connection_timeout)

	conn, err := amqp.DialConfig(rabbit_mq, amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			conn_timeout, _ := strconv.ParseInt(connection_timeout, 10, 32)
			conn_timeout = conn_timeout * 365 * 24 * 60
			return net.DialTimeout(network, addr, time.Duration(conn_timeout)*time.Second)
		},
	})
	failOnError(err, "Failed to connect to RabbitMQ")
	defer fmt.Println("conenction close")
	// defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer fmt.Println("channel close")
	// defer ch.Close()

	if err == nil {
		amqp_channel = ch
		amqp_connection = conn
		is_connection = true
	}
}

func get_channel() *amqp.Channel {
	return amqp_channel
}

func FastPost(url string, referrer string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	req.SetRequestURI(url)
	req.Header.Add("Referer", referrer)
	req.Header.SetMethod("POST")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GetPixelTrack(c *fiber.Ctx) error {
	if !is_connection {
		InitAMQP()
	}

	fmt.Println("is close")
	fmt.Println(amqp_connection.IsClosed())

	if amqp_connection.IsClosed() {
		is_connection = false
		InitAMQP()
	}

	if !amqp_connection.IsClosed() {
		body := c.OriginalURL() + "&HTTP_REFERER=" + string(c.Request().Header.Referer())
		err := amqp_channel.Publish(
			exchange_name, // exchange
			"",            // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		is_log, _ := strconv.ParseBool(go_pixel_log)
		if is_log {
			log.Printf(" [x] Sent %s", body)
			log.Printf("Exchange name %s", exchange_name)
		} else {
			log.Printf(" [x] Sent Data")
		}

		failOnError(err, "Failed to publish a message")

		message := fmt.Sprintf("Success")
		return c.SendString(message)
	}

	return c.SendString("404 not found")
}

func GetPixelPath(c *fiber.Ctx) error {
	log.Println("==============================================")
	log.Println("Get Pixel Path")
	// InitAMQP()
	c.Append("Cache-Control", "public, max-age=300")
	c.Append("content-type", "text/javascript")
	c.Append("Accept-Encoding", "gzip, deflate, brotli")
	c.Append("Expires", time.Now().AddDate(0, 0, 1).Format(http.TimeFormat))

	key := c.Params("key")
	etagKey := key + "_etag"
	db, _ := strconv.Atoi(redis_database)
	rdbWrite := redis.NewClient(&redis.Options{
		Addr:     redis_master_host + ":" + redis_master_port,
		Password: redis_master_password,
		DB:       db,
	})
	rdbRead := redis.NewClient(&redis.Options{
		Addr:     redis_slave_host + ":" + redis_slave_port,
		Password: redis_slave_password,
		DB:       db,
	})
	pixel, err := rdbRead.Get(ctx, key).Result()

	etagPixel, errEtag := rdbRead.Get(ctx, etagKey).Result()
	reqEtag := string(c.Request().Header.Peek("If-None-Match"))

	fmt.Println("If non match : " + reqEtag)

	if errEtag != redis.Nil && err != redis.Nil {
		if etagPixel == reqEtag {
			fmt.Println("is etag match")
			return c.SendStatus(304)
		}
	}

	if err == redis.Nil || len(pixel) == 0 {
		fmt.Println("use rest client")
		url := os.Getenv("PHP_URL") + "/pixel/" + key
		fmt.Println(url)
		resp, resp_err := FastPost(url, string(c.Request().Header.Referer()))

		if resp_err == nil {
			m := minify.New()
			m.AddFunc("text/javascript", js.Minify)
			// pixel_resp, _ := m.String("text/javascript", string(resp.Body()))
			pixel_resp := string(resp.Body())

			fmt.Println(pixel_resp)

			eTag := etag.Generate(pixel_resp, false)
			redis_err := rdbWrite.Set(ctx, key, pixel_resp, 0).Err()

			if redis_err != nil {
				fmt.Println(redis_err)
			}

			etag_err := rdbWrite.Set(ctx, etagKey, eTag, 0).Err()

			if etag_err != nil {
				fmt.Println(etag_err)
			}

			c.Append("ETag", eTag)

			return c.SendString(pixel_resp)
		}

		return c.SendString("404 not found")

	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("use redis")
		eTag := etag.Generate(pixel, false)
		redis_err := rdbWrite.Set(ctx, etagKey, eTag, 0).Err()

		if redis_err != nil {
			fmt.Println(redis_err)
		}

		c.Append("ETag", eTag)
		return c.SendString(pixel)
	}

	return c.SendString("404 not found")
}

func GetRoot(c *fiber.Ctx) error {
	// if !is_connection {
	// 	InitAMQP()
	// }

	// fmt.Println("is close")
	// fmt.Println(amqp_connection.IsClosed())

	// if amqp_connection.IsClosed() {
	// 	is_connection = false
	// 	InitAMQP()
	// }
	// tokenString := getToken(c.Request())
	// if tokenString == "" {
	// 	return c.SendStatus(http.StatusUnauthorized)
	// }
	// token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(jwtSecret), nil
	//   })
  
	//   if err != nil {
	// 	return c.SendStatus(http.StatusUnauthorized)
	//   }
	
	//   claims := token.Claims.(*MyCustomClaims)
	//   query := "SELECT * FROM USERS WHERE username = ?"
	//   row := db.QueryRow(query, claims.username)
	//   var user User
	//   err2 := row.Scan(&user.username, &user.id, &user.role, &user.active)
	//   if err2 != nil {
	// 	return c.SendStatus(http.StatusNotFound)
	//   }
  
	//   return c.JSON(user)

	message := fmt.Sprintf("I am GOPRO!")
	return c.SendString(message)
}


//  authorized handler



func Signup(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	message := fmt.Sprintf("I am GROOT!")
	return c.SendString(message)
}

func Login(c *fiber.Ctx) error {
	 



	loginRequest := new(models.Users)

	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	db, _ := database.ConnectToDB(loginRequest.Prefix)

	user,err := repository.FindUser(db,loginRequest.Preferredname,loginRequest.Password)
	if err != nil {
		response := fiber.Map{
			"Message": "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			"Status":  false,
			"Data":    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
	
		return c.JSON(response)
	}

	
	//day := time.Hour * 24

	claims := jtoken.MapClaims{
		"ID": user.ID,
		"Walletid": user.Walletid,
		"Username": user.Username,
		"Role": user.Role,
		"PartnersKey": user.PartnersKey,
		//"exp":   time.Now().Add(day * 1).Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256,claims)

	t,err_ := token.SignedString([]byte(jwtSecret))
	
	
	if err_ != nil {
		response := fiber.Map{
			"Message": "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			"Status":  false,
			"Data":    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	}
	updates := map[string]interface{}{
		"Token": t,
			}

	// อัปเดตข้อมูลยูสเซอร์
	_err := repository.UpdateFieldsUserString(db,user.Username, updates) // อัปเดตยูสเซอร์ที่มี ID = 1
	if _err != nil {
		response := fiber.Map{
			"Message": "กรุณาตรวจสอบ รหัสผู้ใช้ หรือ รหัสผ่านอีกครั้ง!",
			"Status":  false,
			"Data":    fiber.Map{ 
				"id": -1,
			}, // ตัวอย่างข้อมูลใน data สามารถเป็นโครงสร้างอื่นได้
		}
		return c.JSON(response)
	} else {
		response := fiber.Map{
			"Message": "เข้าระบบสำเร็จ!",
			"Status":  true,
			"Data": fiber.Map{  
				"Token": t, 
				},
		}
		return c.JSON(response)
	}

	
	

	

}

// Protected route
func Protected(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	
	username := claims["Username"].(string)
	favPhrase := claims["PartnersKey"].(string)

	return c.SendString("Bye Bye 👋" + username + " " + favPhrase)
   	
}


func Logout(c *fiber.Ctx) error {
	// var data = formData
	// c.Bind(&data)
	// message := fmt.Sprintf("I am GROOT!")
	// return c.SendString(message)
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	
	username := claims["Username"].(string)
	//favPhrase := claims["PartnersKey"].(string)
	//id := claims["ID"]
	prefix := username[:3] 
	db, _ := database.ConnectToDB(prefix)
	//if claims != nil {
		updates := map[string]interface{}{
			"Token": "",
				}
	
		// อัปเดตข้อมูลยูสเซอร์
		 repository.UpdateFieldsUserString(db,username, updates) 

		response := fiber.Map{
			"Message": "ออกจากระบบสำเร็จ!",
			"Status":  true,
			"Data": fiber.Map{ 
				"id": -1,
			},
		}
		return c.JSON(response)
	// } else {
	// 	response := fiber.Map{
	// 		"Message": "มีข้อผิดพลาด กรุณาออกจากระบบอีกครั้ง!",
	// 		"Status":  false,
	// 		"Data": fiber.Map{ 
	// 			"id": -1,
	// 		},
	// 	}
	// 	return c.JSON(response)
	// }
	 
}

func GenerateSeedPhrase(length int) string {
	seedPhrase := ""
	rand.Seed(time.Now().UnixNano())

	for i:= 0; i < length; i++{
		randomInt := rand.Intn(40)
		seedPhrase = fmt.Sprintf("%s %s",seedPhrase,Words[randomInt])
	}

	return seedPhrase

}


func GetDBFromContext(c *fiber.Ctx) (*gorm.DB, error) {
	dbInterface := c.Locals("db")
	if dbInterface == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "No database connection found")
	}

	// แปลง interface{} ให้เป็น *gorm.DB
	db, ok := dbInterface.(*gorm.DB)
	if !ok {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Invalid database connection")
	}

	return db, nil
}