package main

import (
	"fmt"
	"net"
	"os"

	cors2 "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/your/repo/docs"
	"github.com/your/repo/router"
)

// @title			ManagingServer API
// @version			1.0
// @description		parameter 에 아무런 값을 넣지 않으면, 모든 값을 불러옵니다.
// @contact.name	ManagingServer API Support
// @contact.email	uskawjdu@gmail.com
// @license.name	Apache 2.0
// @license.url		http://www.apache.org/licenses/LICENSE-2.0.html
// @host			uskawjdu.iptime.org
// @BasePath		/
// @Path api/
var testCommand string = "dir /"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	if err != nil {
		panic("큐 생성 에러")
	}

	// tcp
	go TCPServer()

	// Swagger
	go Swagger()

	// HTTP
	HTTPServer()

}

// https://zzihyeon.tistory.com/76
func Swagger() {
	r := gin.Default()

	// Swagger 엔드포인트
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Use(cors2.New(cors2.Config{
		AllowOrigins:     []string{"*"},                            // 모든 도메인 허용, 보안 상 필요한 경우 특정 도메인만 허용해야 함
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 허용할 HTTP 메서드
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Access-Control-Allow-Origin", "Connection", "Accept-Encoding"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Run("0.0.0.0:8001")
}

func TCPServer() {
	listener, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleTCPConnection(conn)
	}

}

func handleTCPConnection(conn net.Conn) {
	// defer conn.Close() // 함수 호출 종료 후 Close

	// HSMgr := HSProtocol.NewHSProtocolManager()

	// for {
	// 	bData := []byte{}
	// 	n, err := conn.Read(bData)

	// 	hs, err := HSMgr.Parsing(bData)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	if n == 0 || err != nil {
	// 		ack := &HSProtocol.HS{ // HSProtocol.ACK
	// 			ProtocolID:     hs.ProtocolID,
	// 			Command:        HSProtocol.ERROR_ACK,
	// 			UUID:           hs.UUID,
	// 			HealthStatus:   hs.HealthStatus,
	// 			Identification: hs.Identification,
	// 			TotalLength:    hs.TotalLength,
	// 			Data:           []byte{},
	// 		}
	// 		rstb, _ := HSMgr.ToBytes(ack)
	// 	}

	// 	fmt.Println("request by hs.uuid : ", hs.UUID)
	// 	dipt := Core.CommandDispatcher{}
	// 	ack, err := dipt.Action(hs)
	// 	if err != nil {
	// 		rstb, _ := HSMgr.ToBytes(ack)
	// 		fmt.Println(err)
	// 		// return ctx.Send(rstb)
	// 		continue
	// 	}

	// 	rstb, err := HSMgr.ToBytes(ack)
	// }
}

func HTTPServer() {
	app := fiber.New()
	app.Get("/view/db", func(c fiber.Ctx) error {
		// HTML 파일을 읽어서 응답으로 반환
		htmlData, err := os.ReadFile("./view/html/viewdata.html")
		if err != nil {
			return c.Status(500).SendString("Error loading page")
		}
		c.Set("Content-Type", "text/html")
		return c.Send(htmlData)
	})

	// 효과적인 Cors 에러 해결
	//app.Use(cors.New(cors.Config{
	//	AllowCredentials: true,
	//	AllowOriginsFunc: func(origin string) bool { return true },
	//}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "*"},
	}))

	router.SetupAPIRoutes(app)
	router.SetupViewRoutes(app)

	fmt.Println("HTTP server listening on port 80")
	err := app.Listen("0.0.0.0:80")
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}

}
