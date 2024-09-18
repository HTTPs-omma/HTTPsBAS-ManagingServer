package router

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/your/repo/Model"
)

// 통합 구조체 정의
type CombinedData struct {
	OperationLog []Model.OperationLogDocument `json:"operation_log"`
	AgentStatus  []Model.AgentStatusRecord    `json:"agent_status"`
	Application  []Model.DapplicationDB       `json:"application"`
	JobData      []Model.JobData              `json:"job_data"`
	SystemInfo   []Model.DsystemInfoDB        `json:"system_info"`
}

func SetupViewRoutes(app *fiber.App) {
	app.Get("/combined-data", func(ctx fiber.Ctx) error {
		dbOperationLog, err := Model.NewOperationLogDB()
		if err != nil {
			return ctx.Status(404).SendString("Error : " + err.Error())
		}
		// 모든 문서 조회
		var dataOL []Model.OperationLogDocument
		dataOL, err = dbOperationLog.SelectAllDocuments()
		if err != nil {
			//dataOL = make([]Model.OperationLogDocument, 0)
			//fmt.Println("SelectAllDocuments Error : " + err.Error())
			//return ctx.Status(404).SendString("Error : " + err.Error())
		}

		dbAgentStatus, err := Model.NewAgentStatusDB()
		if err != nil {
			return ctx.Status(404).SendString("Error : " + err.Error())
		}
		dataAS, err := dbAgentStatus.SelectAllRecords()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			return ctx.Status(404).SendString("Error : " + err.Error())
		}

		appdb, err := Model.NewApplicationDB()
		if err != nil {
			return err
		}
		dataAPP, err := appdb.SelectAllRecords()
		if len(dataAPP) > 200 {
			dataAPP = dataAPP[len(dataAPP)-200:]
		}

		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return ctx.SendString("Error : " + err.Error())
		}

		dbJob, err := Model.NewJobDB()
		if err != nil {
			return err
		}
		dataJD, err := dbJob.SelectAllJobData()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return ctx.SendString("Error : " + err.Error())
		}

		dbSysInfo, err := Model.NewSystemInfoDB()
		if err != nil {
			return err
		}
		dataSys, err := dbSysInfo.SelectAllRecords()
		if err != nil {
			fmt.Println("Error selecting records:", err)
			ctx.Status(404)
			return ctx.SendString("Error : " + err.Error())
		}

		combined := CombinedData{
			OperationLog: dataOL,
			AgentStatus:  dataAS,
			Application:  dataAPP,
			JobData:      dataJD,
			SystemInfo:   dataSys,
		}

		return ctx.JSON(combined)
	})

	app.Get("/view/OperationLogDB", GetOperationLogDB)
	app.Get("/view/agentStatus", GetAgentStatus)
	app.Get("/view/ApplicationDB", GetApplicationDB)
	app.Get("/view/SystemInfoDB", GetSystemInfoDB)
	app.Get("/view/JobDataDB", GetJobDataDB)

	app.Get("/deleted/JobDataDB", func(ctx fiber.Ctx) error {
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

	app.Get("/deleted/OperationLog", func(ctx fiber.Ctx) error {
		db, err := Model.NewOperationLogDB()
		if err != nil {
			return err
		}
		_, err = db.DeleteAllDocument()
		if err != nil {
			fmt.Println("Error Deleted records:", err)
			ctx.Status(404)
			return nil
		}
		return nil
	})

	app.Get("/deleted/SystemInfoDB", func(ctx fiber.Ctx) error {
		db, err := Model.NewSystemInfoDB()
		if err != nil {
			return err
		}
		err = db.DeleteAllRecord()
		if err != nil {
			fmt.Println("Error Deleted records:", err)
			ctx.Status(404)
			return nil
		}
		return nil
	})

	app.Get("/deleted/ApplicationDB", func(ctx fiber.Ctx) error {
		db, err := Model.NewApplicationDB()
		if err != nil {
			return err
		}
		err = db.DeleteAllRecords()
		if err != nil {
			fmt.Println("Error Deleted records:", err)
			ctx.Status(404)
			return nil
		}
		return nil
	})

	app.Get("/deleted/AgentStatusDB", func(ctx fiber.Ctx) error {
		db, err := Model.NewAgentStatusDB()
		if err != nil {
			return err
		}
		err = db.DeleteAllRecord()
		if err != nil {
			fmt.Println("Error Deleted records:", err)
			ctx.Status(404)
			return nil
		}
		return nil
	})
}

// GetOperationLogDB retrieves all operation logs from the database
// @Path			/api/view/OperationLogDB
// @Summary		Get all operation logs
// @Description	Retrieves all operation logs from the database
// @Tags			OperationLogDB
// @Produce		json
func GetOperationLogDB(ctx fiber.Ctx) error {
	db, err := Model.NewOperationLogDB()
	if err != nil {
		ctx.Status(500)
		return ctx.SendString("Failed to connect to the database")
	}

	// 모든 문서 조회
	datas, err := db.SelectAllDocuments()
	if err != nil {
		return ctx.Status(404).JSON(datas)
	}

	return ctx.Status(200).JSON(datas)
}

// GetAgentStatus retrieves agent status by UUID, or all agents if no UUID is provided
// @Path			/api/view/agentStatus
// @Summary		Get agent status
// @Description	Retrieves agent status by UUID, or all agents if no UUID is provided
// @Tags			AgentStatus
// @Produce		json
// @Param			uuid query string false "Agent UUID"
// @Router			/view/agentStatus [get]
func GetAgentStatus(ctx fiber.Ctx) error {
	db, err := Model.NewAgentStatusDB()
	uuid := ctx.Query("uuid")
	var datas []Model.AgentStatusRecord
	if uuid != "" {
		datas, err = db.SelectRecordByUUID(uuid)
	} else {
		datas, err = db.SelectAllRecords()
	}
	if err != nil {
		fmt.Println("Error selecting records:", err)
		ctx.Status(404)
		return ctx.Status(404).SendString("Error : " + err.Error())
	}

	return ctx.Status(200).JSON(datas)
}

// GetApplicationDB retrieves application database records by UUID, or all records if no UUID is provided
// @Path			/api/view/ApplicationDB
// @Summary		Get application database records
// @Description	Retrieves application database records by UUID, or all records if no UUID is provided
// @Tags			ApplicationDB
// @Produce		json
// @Param			uuid query string false "Application UUID"
// @Router			/view/ApplicationDB [get]
func GetApplicationDB(ctx fiber.Ctx) error {
	fmt.Println("ApplicationDB logging")
	db, err := Model.NewApplicationDB()
	if err != nil {
		return err
	}

	uuid := ctx.Query("uuid")
	var datas []Model.DapplicationDB
	fmt.Println(uuid)
	if uuid != "" {
		datas, err = db.SelectRecordByUUID(uuid)
	} else {
		datas, err = db.SelectAllRecords()
	}
	if err != nil {
		fmt.Println("Error selecting records:", err)
		return ctx.Status(404).JSON(datas)
	}
	return ctx.Status(200).JSON(datas)
}

// GetSystemInfoDB retrieves system information records by UUID, or all records if no UUID is provided
// @Path			/api/view/SystemInfoDB
// @Summary		Get system information records
// @Description	Retrieves system information records by UUID, or all records if no UUID is provided
// @Tags			SystemInfoDB
// @Produce		json
// @Param			uuid query string false "System UUID"
// @Router			/view/SystemInfoDB [get]
func GetSystemInfoDB(ctx fiber.Ctx) error {
	db, err := Model.NewSystemInfoDB()
	if err != nil {
		return err
	}

	uuid := ctx.Query("uuid")
	var datas []Model.DsystemInfoDB
	if uuid != "" {
		datas, err = db.SelectRecordByUUID(uuid)
	} else {
		datas, err = db.SelectAllRecords()
	}
	if err != nil {
		fmt.Println("Error selecting records:", err)
		return ctx.Status(404).JSON(datas)
	}
	return ctx.Status(200).JSON(datas)
}

// GetJobDataDB retrieves all job data from the database
// @Path			/api/view/JobDataDB
// @Summary		Get all job data
// @Description	Retrieves all job data from the database
// @Tags			JobDataDB
// @Produce		json
// @Router			/view/JobDataDB [get]
func GetJobDataDB(ctx fiber.Ctx) error {
	db, err := Model.NewJobDB()
	if err != nil {
		return err
	}
	datas, err := db.SelectAllJobData()
	if err != nil {
		fmt.Println("Error selecting records:", err)
		return ctx.Status(404).JSON(datas)
	}
	ctx.Status(200)
	return ctx.JSON(datas)
}
