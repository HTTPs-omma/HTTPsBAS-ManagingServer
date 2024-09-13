package router

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/your/repo/Model"
)

func SetupViewRoutes(app *fiber.App) {
	// /view/OperationLogDB 라우트 정의
	app.Get("/view/OperationLogDB", func(ctx fiber.Ctx) error {
		db, err := Model.NewOperationLogDB()
		if err != nil {
			ctx.Status(500)
			return ctx.SendString("Failed to connect to the database")
		}

		// 모든 문서 조회
		datas, err := db.SelectAllDocuments()
		if err != nil {
			ctx.Status(404)
			return ctx.SendString("Documents not found")
		}

		// 성공 시 JSON으로 응답
		ctx.Status(200)
		return ctx.JSON(datas)
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

	app.Get("/view/sdsa", func(ctx fiber.Ctx) error {
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
}
