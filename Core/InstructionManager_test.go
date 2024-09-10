package Core

import (
	"fmt"
	"testing"
)

// InstructionManager의 기본 동작을 테스트
func TestInstructionManager_LoadAndGetByID(t *testing.T) {

	cm, err := NewInstructionManager()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	// ID로 InstructionData 가져오기
	id := "P_Collection_Kimsuky_001"
	command, exists := cm.GetByID(id)
	if !exists {
		t.Fatalf("Expected command with ID %s not found", id)
	}

	fmt.Println(command)
}
