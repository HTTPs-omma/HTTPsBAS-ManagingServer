package Model

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"testing"
	"time"
)

func TestInsertDocument(t *testing.T) {
	tests := []struct {
		name string
		log  OperationLogDocument
	}{
		{
			name: "Insert valid document",
			log: OperationLogDocument{
				AgentUUID:       "agent-123",
				ProcedureID:     "tech-456",
				InstructionUUID: "msg-789",
				ConductAt:       time.Now(),
				ExitCode:        0,
				Log:             "Test log message",
				Command:         "Test command",
			},
		},
	}

	for _, tt := range tests {

		// MongoDB 클라이언트 옵션 설정
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		t.Run(tt.name, func(t *testing.T) {
			db, err := NewOperationLogDB()
			if err != nil {
				log.Fatal("에러: ", err)
			}

			// insertDocument 호출
			result, err := db.InsertDocument(tt.log)
			if err != nil {
				t.Errorf("insertDocument() 에러: %v", err)
			}

			// 결과 ID가 있는지 확인
			if result == nil || result.InsertedID == nil {
				t.Errorf("insertDocument() 결과가 유효하지 않습니다.")
			} else {
				t.Logf("삽입된 문서 ID: %v", result.InsertedID)
			}

			log, err := db.SelectDocumentById(tt.log.InstructionUUID)
			if err != nil {
				t.Fatalf("select 오류")
			}
			fmt.Printf("성공 : " + log.InstructionUUID)

			//rst2, err := db.UpdateDocumentByInstID(tt.log.InstructionUUID,
			//&OperationLogDocument{
			//	AgentUUID:   "agent-123",
			//	ProcedureID: "tech-456",
			//	InstructionUUID: "msg-789",
			//	ConductAt:   time.Now(),
			//	ExitCode:    0,
			//	Log:         "Test log message",
			//	Command:     "Test command",
			//})
			//if err == nil {
			//	t.Fatalf("update 오류")
			//}

			rst, err := db.DeleteDocumentByInstID(tt.log.InstructionUUID)
			if err != nil {
				t.Errorf("DeleteDocumentByInstID() 에러: %v", err)
			}
			// 결과 ID가 있는지 확인
			if rst == nil || rst.DeletedCount == 0 {
				t.Errorf("DeleteDocument() 결과가 유효하지 않습니다.")
			} else {
				t.Logf("정상적으로 삭제됨")
				//t.Logf("삭제된 문서 ID: %v", rst.)
			}

		})
	}
}
