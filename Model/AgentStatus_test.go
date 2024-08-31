package Model

import (
	"testing"
	"time"
)

func TestAgentStatusDB(t *testing.T) {
	// 새로운 DB 인스턴스 생성
	db := NewAgentStatusDB()

	// 테이블 생성 테스트
	if err := db.CreateTable(); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// 테스트 데이터 생성
	record := &AgentStatusRecord{
		UUID:      "test-uuid",
		Status:    Running,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// InsertRecord 테스트
	if err := db.InsertRecord(record); err != nil {
		t.Fatalf("Failed to insert record: %v", err)
	}

	// SelectRecords 테스트 - 삽입 후 조회
	records, err := db.SelectRecords()
	if err != nil {
		t.Fatalf("Failed to select records: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(records))
	}

	if records[0].UUID != record.UUID || records[0].Status != record.Status {
		t.Fatalf("Record mismatch: expected %v, got %v", record, records[0])
	}

	// UpdateRecord 테스트 - 상태 변경
	record.Status = Stopping
	if err := db.UpdateRecord(record); err != nil {
		t.Fatalf("Failed to update record: %v", err)
	}

	// 업데이트 후 다시 조회
	updatedRecords, err := db.SelectRecords()
	if err != nil {
		t.Fatalf("Failed to select records after update: %v", err)
	}

	if updatedRecords[0].Status != Stopping {
		t.Fatalf("Expected status 'inactive', got '%s'", updatedRecords[0].Status)
	}

	// DeleteRecord 테스트
	if err := db.DeleteRecord(record.UUID); err != nil {
		t.Fatalf("Failed to delete record: %v", err)
	}

	// 삭제 후 조회하여 레코드가 없는지 확인
	finalRecords, err := db.SelectRecords()
	if err != nil {
		t.Fatalf("Failed to select records after delete: %v", err)
	}

	if len(finalRecords) != 0 {
		t.Fatalf("Expected 0 records, got %d", len(finalRecords))
	}

	// ExistRecord 테스트 - 데이터가 없는 상태에서 확인
	exists, err := db.ExistRecord()
	if err != nil {
		t.Fatalf("Failed to check if record exists: %v", err)
	}

	if exists {
		t.Fatalf("Expected no records to exist, but some do")
	}
}
