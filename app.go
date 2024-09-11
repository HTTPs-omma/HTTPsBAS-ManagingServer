package main

import (
	"bytes"
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
	ProcedureID     string `json:"procedureID"`
	AgentUUID       string `json:"agentUUID"`
	InstructionUUID string `json:"instructionUUID"`
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
		dipt := Core.CommandDispatcher{}
		dipt.Action(hs)

		return nil
	})

	app.Post("/postInstruction", func(ctx fiber.Ctx) error {
		//https://github.com/gofiber/fiber/issues/2958
		InstD := new(InstructionD)
		err := ctx.Bind().JSON(InstD)
		if err != nil {
			fmt.Println("Error marshaling to JSON:", err)
			ctx.Status(404)
			return err
		}
		jobMgr, err := Core.NewJobManager()
		if err != nil {
			return err
		}
		fmt.Println("test : ", InstD.ProcedureID, InstD.AgentUUID, InstD.InstructionUUID)

		err = jobMgr.InsertData(&Model.JobData{
			InstD.ProcedureID,
			InstD.AgentUUID,
			InstD.InstructionUUID,
			time.Now(),
		})
		if err != nil {
			fmt.Println("Error inserting data:", err)
			return fmt.Errorf("Error inserting data into job manager: %v", err)
		}

		ctx.Status(200)
		return ctx.JSON(fiber.Map{
			"status": true,
		})
	})

	app.Get("/view/agentStatus", func(ctx fiber.Ctx) error {
		//data := ctx.Body()
		db := Model.NewAgentStatusDB()
		datas, err := db.SelectRecords()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return nil
		}
		ctx.Status(200)
		return ctx.JSON(datas)
	})

	app.Get("/view/ApplicationDB", func(ctx fiber.Ctx) error {
		fmt.Println("ApplicationDB loging")
		db := Model.NewApplicationDB()
		datas, err := db.SelectAllRecords()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return nil
		}
		ctx.Status(200)
		return ctx.JSON(datas)
	})

	app.Get("/view/OperationLogDB", func(ctx fiber.Ctx) error {
		db, _ := Model.NewOperationLogDB()
		datas, err := db.SelectAllDocuments()
		if err != nil {
			ctx.Status(404)
			return nil
		}
		ctx.Status(200)
		return ctx.JSON(datas)
	})

	app.Get("/view/SystemInfoDB", func(ctx fiber.Ctx) error {
		//data := ctx.Body()
		db := Model.NewSystemInfoDB()
		datas, err := db.SelectRecords()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return nil
		}
		ctx.Status(200)
		return ctx.JSON(datas)
	})

	app.Get("/view/JobDataDB", func(ctx fiber.Ctx) error {
		db, err := Model.NewJobDB()
		if err != nil {
			return err
		}
		datas, err := db.SelectAllJobData()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return nil
		}
		ctx.Status(200)
		return ctx.JSON(datas)
	})

	app.Get("/deleted/DeletedJobDataDB", func(ctx fiber.Ctx) error {
		db, err := Model.NewJobDB()
		if err != nil {
			return err
		}
		err = db.DeleteAllJobData()
		if err != nil {
			fmt.Println("Error Deleted records:", err)
			ctx.Status(404)
			return nil
		}
		return nil
	})

	fmt.Println("HTTP server listening on port 80")
	err := app.Listen(":80")
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}
