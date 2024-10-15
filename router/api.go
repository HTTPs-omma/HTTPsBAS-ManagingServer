package router

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/your/repo/Core"
	"github.com/your/repo/Model"
	_ "github.com/your/repo/docs"
)

type AgentAction string

const (
	ExecutePayLoad AgentAction = "ExecutePayLoad"
	ExecuteCleanUp AgentAction = "ExecuteCleanUp"
	GetSystemInfo  AgentAction = "GetSystemInfo"
	GetApplication AgentAction = "GetApplication"
	ChangeProtocol AgentAction = "ChangeProtocol"
)

// swagger:parameters Request
type InstructionD struct {
	ProcedureID string `json:"procedureID" default:"P_PrivilegeEscalation_0001"`
	AgentUUID   string `json:"agentUUID" default:"5610eb3154c742d4bc95ce9194166ac4"`
	Action      string `json:"action" default:"ExecutePayLoad"`
	//MessageUUID string `json:"messageUUID" default:"32a2833486414af9bc4596caef585538"`
}

// @title			ManagingServer API
// @version		1.0
// @description	This is a sample server for the ManagingServer project.
// @termsOfService	http://managingserver.io/terms/
// @contact.name	ManagingServer API Support
// @contact.url	http://managingserver.io/support
// @contact.email	support@managingserver.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			uskawjdu.iptime.org:80
// @BasePath		/
// @Path			/api
func SetupAPIRoutes(app *fiber.App) {

	app.Post("/api/checkInstReq", checkInstReq)
	app.Post("/api/postInst", postInst)
	app.Post("/ipinfo", updateIPInfo)

}

func checkInstReq(ctx fiber.Ctx) error {
	req := ctx.Body()
	HSMgr := HSProtocol.NewHSProtocolManager()
	//fmt.Println("debug")

	hs, err := HSMgr.Parsing(req)
	if err != nil {
		fmt.Println(err)
		ctx.Status(404)
		return fmt.Errorf("Error parsing:", err)
	}

	//fmt.Println()

	//fmt.Println("hs.uuid : ", hs.UUID)
	dipt := Core.CommandDispatcher{}
	ack, err := dipt.Action(hs)
	if err != nil {
		ack := &HSProtocol.HS{ // HSProtocol.ACK
			ProtocolID:     hs.ProtocolID,
			Command:        HSProtocol.ERROR_ACK,
			UUID:           hs.UUID,
			HealthStatus:   hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength:    hs.TotalLength,
			Data:           []byte{},
		}
		rstb, _ := HSMgr.ToBytes(ack)
		fmt.Println(err)
		return ctx.Send(rstb)
	}

	rstb, err := HSMgr.ToBytes(ack)
	return ctx.Send(rstb)
}


// postInst example
//
//	@Summary		ExecutePayLoad ExecuteCleanUp GetSystemInfo GetApplication StopAgent
//	@Description	이 API는 다양한 작업을 실행하는 명령어를 포함합니다.
//	@ID				get-struct-array2-by-string
//	@Accept			json
//	@Produce		json
//	@Param			loginUserRequest	body	InstructionD	true	"request job. 'Action' 필드에는 'ExecutePayLoad', 'ExecuteCleanUp', 'GetSystemInfo', 'GetApplication', 'StopAgent', 'ChangeProtocol' 값이 들어갈 수 있습니다."
//	@Router			/api/postInst [post]
func postInst(ctx fiber.Ctx) error {
	// https://github.com/gofiber/fiber/issues/2958
	InstD := new(InstructionD)
	err := ctx.Bind().JSON(InstD)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return ctx.Status(404).Send([]byte(err.Error()))
	}
	jobdb, err := Model.NewJobDB()
	if err != nil {
		return ctx.Status(404).Send([]byte(err.Error()))
	}
	//fmt.Println("test : ", InstD.ProcedureID, InstD.AgentUUID, InstD.MessageUUID)

	newUUID := uuid.New()
	MessageUUID := newUUID.String()
	MessageUUID = strings.Replace(MessageUUID, "-", "", -1)
	err = jobdb.InsertJobData(&Model.JobData{
		0,
		InstD.ProcedureID,
		InstD.AgentUUID,
		MessageUUID,
		InstD.Action,
		time.Now(),
	})
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return fmt.Errorf("Error inserting data into job manager: %v", err)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"MessageUUID": MessageUUID,
	})
}

func updateIPInfo(c fiber.Ctx) error {
	fmt.Println("================================================")
	// 클라이언트에서 보낸 JSON 데이터를 파싱
	var data struct {
		PrivateIP []string `json:"PrivateIP"`
		PublicIP  string   `json:"PublicIP"`
		agentUUID string   `json:"agentUUID"`
	}

	if err := json.Unmarshal(c.Body(), &data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request : Json Parsing Failed")
	}
	fmt.Println("================================================")
	fmt.Println(data)

	// DB에 기록하거나 업데이트
	agentDB, err := Model.NewSystemInfoDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to connect to the database")
	}

	err = agentDB.UpdateIP(data.PublicIP, data.agentUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update record")
	}

	return c.SendString("IP information updated successfully")
}
