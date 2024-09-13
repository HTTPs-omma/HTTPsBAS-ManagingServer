package router

import (
	"fmt"
	"github.com/HTTPs-omma/HSProtocol/HSProtocol"
	"github.com/gofiber/fiber/v3"
	"github.com/your/repo/Core"
	"github.com/your/repo/Model"
	_ "github.com/your/repo/docs"
	"time"
)

type InstructionD struct {
	ProcedureID     string `json:"procedureID"`
	AgentUUID       string `json:"agentUUID"`
	InstructionUUID string `json:"instructionUUID"`
}

//	@title			ManagingServer API
//	@version		1.0
//	@description	This is a sample server for the ManagingServer project.
//	@termsOfService	http://managingserver.io/terms/
//	@contact.name	ManagingServer API Support
//	@contact.url	http://managingserver.io/support
//	@contact.email	support@managingserver.io
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//	@host			localhost:80
//	@BasePath		/
//	@Path			/api

func SetupAPIRoutes(app *fiber.App) {

	app.Post("/api/checkInstReq", checkInstReq)
	app.Post("/api/postInst", postInst)
}

func checkInstReq(ctx fiber.Ctx) error {
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
}

// postInst example
//
//	@Description	get struct array by ID
//	@ID				get-struct-array2-by-string
//	@Accept			json
//	@Produce		json
//	@Param			procedureID		query		string	true	"Procedure ID (default: 'P_Collection_Kimsuky_001')"
//	@Param			agentUUID		query		string	true	"Agent UUID (default: '123e4567e89b12d3a456426614174000')"
//	@Param			instructionUUID	query		string	true	"Instruction UUID (default: '550e8400e29b41d4a716446655440000')"
//	@Success		200				{string}	string	"ok"
//	@Default		{"procedureID": "How to implement a Library", "agentUUID": "2023-09-15T10:00:00Z", "instructionUUID" : :"550e8400-e29b-41d4-a716-446655440000"}
//	@Router			/api/postInst [post]
func postInst(ctx fiber.Ctx) error {
	// https://github.com/gofiber/fiber/issues/2958
	InstD := new(InstructionD)
	err := ctx.Bind().JSON(InstD)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		ctx.Status(404)
		return ctx.Send([]byte(err.Error()))
	}
	jobMgr, err := Core.NewJobManager()
	if err != nil {
		return ctx.Send([]byte(err.Error()))
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
}
