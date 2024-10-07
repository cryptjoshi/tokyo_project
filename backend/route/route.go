package route

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/etag"
	
	"hanoi/handler"
	"hanoi/handler/users"
	"hanoi/handler/wallet"
	//"hanoi/middlewares"
	"hanoi/handler/ef" 
	"hanoi/handler/gc" 
	"hanoi/handler/pg" 
	"hanoi/handler/jwtn"
	//"github.com/swaggo/fiber-swagger"
	"os"
)
var jwtSecret = os.Getenv("PASSWORD_SECRET")

func ProviderMiddleware(c *fiber.Ctx) error {
	username := c.FormValue("username") // Assuming username is part of the request
	if len(username) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid username",
		})
	}
	
	// Extract prefix (first 3 characters)
	prefix := username[:3]

	// Determine provider based on prefix or logic
	switch prefix {
	case "EFI":
		return ef.GetBalance(c) // EFinity
	case "PGS":
		return pg.GetBalance(c) // PGSoft
	case "GCL":
		return gc.GetBalance(c) // GClub
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid provider or prefix",
		})
	}
}

func SetupRoutes(app fiber.Router) {
	// app.Use(etag.New())
	

	//jwtm := middlewares.NewAuthMiddleware(jwtSecret)

	// app.Static("/css", "./css")
	// app.Static("/js/libraries", "./js")
	//app.Get("/",handler.GetRoot)

	// เส้นทาง Swagger
	//app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// เส้นทาง API สำหรับดึงข้อมูลผู้ใช้งาน
	// user
	app.Post("/users/all",jwt.JwtMiddleware, users.GetUsers)
	app.Post("/users/login", users.Login)
	app.Post("/users/register",users.Register)
	app.Post("/users/balance",jwt.JwtMiddleware,users.GetBalance)
	app.Post("/users/sum/balance",jwt.JwtMiddleware,users.GetBalanceSum)
	app.Post("/users/info",jwt.JwtMiddleware,users.GetUser)
	app.Post("/users/statement",jwt.JwtMiddleware,users.GetUserStatement)
	app.Post("/users/transactions",jwt.JwtMiddleware,users.GetUserTransaction)
    app.Delete("/users/logout",jwt.JwtMiddleware,users.Logout)

	app.Get("/protected", jwt.JwtMiddleware, handler.Protected)
	//app.Post("/api/gateway/getBalance",jwt.JwtMiddleware, users.GetBalance)

	// Define individual routes for each provider (if needed)
	app.Post("/callback/Seamless/GetBalance", ef.GetBalance)
	app.Post("/callback/pgsoft/checkBalance", pg.GetBalance)
	app.Post("/api/Wallet/Balance", gc.GetBalance)
 


	// wallet
	// app.Post("/wallet/withdraw",wallet.WithDraw)
	// app.Post("/wallet/deposit",wallet.AddStatement)
	app.Post("/statement/update",wallet.UpdateStatement)
	app.Post("/statement/add",wallet.AddStatement)
	// app.Post("/transaction/add",handler.AddTransactions)
	// app.Post("/transaction/update",handler.UpdateTransactions)
	

	// dashboard
	// app.Post("/api/bank/statement",handler.GetBankStatement)
	// app.Post("/api/first/statement",handler.GetFirstUsers)
	// app.Post("/api/user/all/statement",handler.GetUserAllStatement)
	// app.Post("/api/user/statement",jwt,handler.GetUserStatement)
	// app.Post("/api/user/statement/id",handler.GetIdStatement)
	// app.Post("/api/user/sum/statement",handler.GetUserSumStatement)


	// app.Post("api/user/token",jwt,handler.UpdateToken)
	// app.Post("api/user/byid",handler.GetBalanceFromID)

	 
}