package jwt

import (
	//"bytes"
	//"crypto/cipher"
	//"crypto/des"
	//"encoding/base64"
	"hanoi/models"
	"hanoi/database"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	//"github.com/golang-jwt/jwt"
	//jtoken "github.com/golang-jwt/jwt/v4"
	//"github.com/solrac97gr/basic-jwt-auth/config"
	//"github.com/solrac97gr/basic-jwt-auth/models"
	//"github.com/solrac97gr/basic-jwt-auth/repository"
	// "pkd/repository"
	//"log"
	// "net"
	// "net/http"
	"encoding/json"
	"os"
	//"strconv"
	"time"
	//"strings"
	"fmt"
	"errors"
	"gorm.io/gorm"
	"regexp"
	
)
var jwtKey  = []byte(os.Getenv("PASSWORD_SECRET"))
//var CLIENT_ID = "6342e1be-fa03-456f-8d2d-8e1c9513c351" //[]byte(os.Getenv("CLIENT_ID"))
//var CLIENT_SECRET = "6d83ac42" //[]byte(os.Getenv("CLIENT_SECRET"))
//var DESKEY = "9c62a148"
//var DESIV =	"8e014099"
//var SYSTEMCODE = "tsxthb"
//var WEBID = "tsxthb"



// Struct สำหรับ JWT Claims

// type ECResult struct {
// 	Enc string `json:"enc"`
// 	Unx int64 `json:"unx"` // ค่า unx คุณสามารถกำหนดเอง
// 	Des string `json:"des"` // ค่า dex คุณสามารถกำหนดเอง
// }




type Claims struct {
    Username string `json:"username"`
	Id int `json:"id"`
	Role string `json:"role"`
	Prefix string `json:"prefix"`
	Walletid int `json:"walletid"`
	Checker string `json:"checker"`
    jwt.RegisteredClaims
}


func GetPrefix(input string) (string, error) {
	// ใช้ regexp เพื่อจับเฉพาะตัวอักษรก่อนตัวเลข
	re := regexp.MustCompile(`^[a-zA-Z]+`)
	matches := re.FindString(input)
	if matches == "" {
		return "", fmt.Errorf("No prefix found")
	}
	return matches, nil
}

func CheckDBConnection(db interface{},prefix string) (*gorm.DB, error) {
	// ตรวจสอบว่า db ไม่เป็น nil
	if db == nil {
		db, _ := database.ConnectToDB(prefix)
		return db,nil
		//return nil, fmt.Errorf("database connection is nil")
	}

	// พยายามแปลงค่าเป็น *gorm.DB
	dbConnection, ok := db.(*gorm.DB)
	if !ok {
		db, _ := database.ConnectToDB(prefix)
		return db,nil
		//return nil, fmt.Errorf("interface conversion failed: interface {} is not *gorm.DB")
	}

	return dbConnection, nil
}
func MapToJSONString(data fiber.Map) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
func ValidateJWTReturn(tokenString string) models.Users {
	claims := &Claims{}
	//dbClaims := &Claims{}
	//tokenString := c.Get("Authorization")[7:]
	// token, claims, err := ValidateJWT(tokenString) // เรียกใช้ฟังก์ชันจาก utils
	// if err != nil || !token.Valid {
	// 	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	// }
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	username := claims.Username
	 
	prefix,_ := GetPrefix(username) //username[:3] // Extract prefix

	// Connect to the database based on the prefix
	db, err := database.ConnectToDB(prefix)
	//checkerFromRequest := claims.Checker
	var user models.Users
	//fmt.Println(err)
	if err==nil {
		db.Select("id,username,balance,Token").Where("username = ?", username).First(&user)
	}

	 
	return user
	 
}
func ValidateJWT(tokenString string) (error) {
	claims := &Claims{}
	//dbClaims := &Claims{}
	//tokenString := c.Get("Authorization")[7:]
	// token, claims, err := ValidateJWT(tokenString) // เรียกใช้ฟังก์ชันจาก utils
	// if err != nil || !token.Valid {
	// 	return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	// }
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	username := claims.Username
	//checkerFromRequest := claims.Checker

	

	// ดึง JWT Token ที่เก็บไว้ในฐานข้อมูลสำหรับผู้ใช้ที่เกี่ยวข้อง
	var user models.Users
	//prefix := username[:3] // Extract prefix
	prefix,_ := GetPrefix(username)
	// Connect to the database based on the prefix
	db, err := database.ConnectToDB(prefix)
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return result.Error
	}
	
	utoken := user.Token
	fmt.Println(utoken)
	//fmt.Println(tokenString)
	// ตรวจสอบและเปรียบเทียบค่า checker
	// _,err_ := jwt.ParseWithClaims(utoken, dbClaims, func(token *jwt.Token) (interface{}, error) {
    //     return jwtKey, nil
    // })
	
	
	// if err_!= nil {
	// 	fmt.Println("77")
	// 	fmt.Println(err_)
	// 	return err_
	// }
	// checkerFromDB := dbClaims.Checker
	// fmt.Println(&dbClaims)
	// fmt.Println(&claims)
	// แสดงค่า checker จาก request และจากฐานข้อมูล
	//fmt.Printf("Checker from request token: %s\n", checkerFromRequest)
	//fmt.Printf("Checker from DB token: %s\n", checkerFromDB)

	// เปรียบเทียบค่า checker
	if utoken != tokenString {
		//return c.Status(fiber.StatusUnauthorized).SendString("Checker mismatch")
		return errors.New("checker ไม่ตรง!")
	}
	return err
}
// ฟังก์ชันสำหรับตรวจสอบและแยก JWT Token
func CheckedJWT(db *gorm.DB,c *fiber.Ctx) (error) {
	
	tokenString := c.Get("Authorization")[7:] 
    claims := &Claims{}
    
    _, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
	
	var user models.Users
	//prefix := claims.Username[:3] // Extract prefix
	
	// Connect to the database based on the prefix
	//db, err := database.ConnectToDB(prefix)
	result := db.Debug().Where("username = ?", claims.Username).First(&user)

	if result.Error != nil {
		//http.Error(w, "User not found", http.StatusUnauthorized)
		return errors.New("มีข้อผิดพลาด")
	}
	fmt.Println("------------")
	fmt.Println(claims.Checker)
	fmt.Println("------------")
	fmt.Println(tokenString)
	fmt.Println("------------")

	// ตรวจสอบว่า token ที่ส่งมาไม่ตรงกับ token ที่เก็บในฐานข้อมูล
	if user.Token != tokenString {
		//http.Error(w, "Token ไม่ตรง", http.StatusUnauthorized)
		return errors.New("มีข้อผิดพลาด")
	}

	// หาก token ถูกต้องและตรงกัน
	//fmt.Fprintf(w, "Hello, %s", claims.Username)

    return err
}
// ฟังก์ชันสำหรับสร้าง JWT Token (เพื่อใช้ทดสอบ)
func createJWT(username string) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}
func ExtractPrefixFromToken(c *fiber.Ctx) (string, error) {
	// ดึง token จาก header Authorization
	tokenString := c.Get("Authorization")[7:] 
	//tokenString := c.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("ไม่พบโทเคน!")
	}

	// ถอดรหัส token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// นี่คือที่สำหรับการ validate signature key, เช่นการใช้ secret
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	// ตรวจสอบ claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ดึง prefix จาก claims
		if prefix, ok := claims["prefix"].(string); ok {
			return prefix, nil
		}
		return "", fmt.Errorf("ไม่พบ prefix ใน token")
	}

	return "", fmt.Errorf("โทเคน ผิดผลาด!")
}
func JwtMiddleware(c *fiber.Ctx) error {
	
	claims := &Claims{}
	tokenString := c.Get("Authorization")[7:]
 	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
    })
 	if err==nil {
		db, _ := database.ConnectToDB(claims.Prefix)
        
		fmt.Println("claims",claims.Walletid)

		c.Locals("Walletid", claims.Walletid)
        c.Locals("ID", claims.ID)
        c.Locals("username", claims.Username)
       // c.Locals("PartnersKey",claims.PartnersKey)
        c.Locals("role", claims.Role)
        c.Locals("prefix", claims.Prefix)
	 	c.Locals("db", db)
        return c.Next()
    } else {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "โทเคน ผิดผลาด!"})
    } 
	 
}
func jwtMiddleware(c *fiber.Ctx) error {
    tokenString := c.Get("Authorization")
	fmt.Println(tokenString)
    if tokenString == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
    }

    token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtKey), nil
    })


    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
        // ดึงข้อมูล claims เช่น id, role หรือ prefix
        c.Locals("Walletid", claims["walletid"])
        c.Locals("ID", claims["id"])
        c.Locals("username", claims["username"])
        c.Locals("PartnersKey",claims["partnersKey"])
        c.Locals("role", claims["role"])
        c.Locals("prefix", claims["prefix"])
        return c.Next()
    } else {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }
}