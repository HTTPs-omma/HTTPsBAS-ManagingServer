package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HTTPs-omma/HSProtocol/HSProtocol"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/your/repo/Core"
	"github.com/your/repo/Model"
	"net"
	"time"
)

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

	// udp
	////go UDPServer()

	// HTTP
	HTTPServer()

}

func TCPServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 8080")

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

type InstructionD struct {
	agentUUID       string `json:"agentUUID"`
	procedureID     string `json:"procedureID"`
	instructionUUID string `json:"instructionUUID"`
}

// HTTP 서버 함수 (Fiber 사용)
func HTTPServer() {
	app := fiber.New()

	app.Post("/getPacket", func(ctx fiber.Ctx) error {
		req := ctx.Body()
		HSMgr := HSProtocol.NewHSProtocolManager()
		hs, err := HSMgr.Parsing(req)
		if err != nil {
			ctx.Status(404)
			return fmt.Errorf("Error parsing:", err)
		}

		fmt.Println("hs.uuid : ", hs.UUID)

		return nil
	})

	// POST 요청을 처리하는 핸들러
	app.Post("/postInstruction", func(ctx fiber.Ctx) error {
		data := ctx.Body()
		InstD := &InstructionD{}
		err := json.Unmarshal(data, &InstD)
		if err != nil {
			fmt.Println("Error marshaling to JSON:", err)
			ctx.Status(404)
			return err
		}
		jobMgr, err := Core.NewJobManager()
		if err != nil {
			return err
		}

		jobMgr.InsertData(&Model.JobData{
			InstD.procedureID,
			InstD.agentUUID,
			InstD.instructionUUID,
			time.Now(),
		})

		ctx.Status(200)
		return ctx.JSON(fiber.Map{
			"status": true,
		})
	})

	fmt.Println("HTTP server listening on port 80")
	err := app.Listen(":80")
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}
