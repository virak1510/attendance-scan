package routers

import (
	"attendance/internal/middlewares"
	"attendance/internal/modules/attendances"
	"attendance/internal/modules/users"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.0.242:3000", "https://notlocal.local, https://attendance-scan-production.up.railway.app"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowWebSockets:  true,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	var (
		userService       = users.NewService(db)
		attendanceService = attendances.NewService(db)
	)

	var (
		userController       = users.NewController(userService)
		attendanceController = attendances.NewController(attendanceService)
	)

	public := router.Group("/api")
	// Public routes
	{
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// will change on production
				// var origin = r.Header.Get("Origin")
				// if origin == "http://localhost:3000" || origin == "http://192.168.0.242:3000" {
				// 	return true
				// }
				return true
			},
		}

		public.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
		public.GET("/ws", func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				log.Print("Error during connection upgradation:", err)
				return
			}
			defer conn.Close()

			// The event loop
			for {
				messageType, message, err := conn.ReadMessage()
				if err != nil {
					log.Println("Error during message reading:", err)
					break
				}
				log.Printf("Received: %s", message)
				err = conn.WriteMessage(messageType, message)
				if err != nil {
					log.Println("Error during message writing:", err)
					break
				}
			}
		})
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/users", userController.GetAllUsers)
		// protected.POST("/users", userController.Register)

		// Attendance routes
		protected.GET("/attendances", attendanceController.GetAttendances)
		protected.POST("/attendances", attendanceController.CheckIn)
		protected.PUT("/attendances", attendanceController.CheckOut)
	}

	return router
}
