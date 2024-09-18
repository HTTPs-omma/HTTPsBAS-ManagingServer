package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
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
	defer conn.Close() // 함수 호출 종료 후 Close

	buffer := make([]byte, 1024*1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			break
		}
		if n < 1 {
			continue
		}

		HSMgr := HSProtocol.NewHSProtocolManager()
		hs, err := HSMgr.Parsing(buffer)
		if err != nil {
			fmt.Println("Error parsing:", err)
			continue
		}

		if hs.Command == 0b0000000110 { // payload 를 받아옴
			conn.Write([]byte(testCommand))
		}

		if hs.Command == 0b0000000111 { // 실행 결과를 작성함.
			msg := bytes.ReplaceAll(hs.Data, []byte{0x00}, []byte{})
			fmt.Println("Log : ", string(msg))
		}

	}
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
