package Core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

/**
ChatGpt 로 생성한 코드임
*/

// InstructionData는 주어진 YAML 데이터를 저장할 구조체입니다.
type InstructionData struct {
	ID               string `yaml:"id"`
	MITREID          string `yaml:"MITRE_ID"`
	Description      string `yaml:"Description"`
	Escalation       bool   `yaml:"Escalation"` // 새로운 필드 추가
	Tool             string `yaml:"tool"`
	RequisiteCommand string `yaml:"requisite_command"`
	Command          string `yaml:"command"`
	Cleanup          string `yaml:"cleanup"`
}

// ToBytes는 InstructionData 구조체를 YAML 바이트 슬라이스로 변환하는 함수입니다.
func (cd *InstructionData) ToBytes() ([]byte, error) {
	// YAML로 직렬화
	data, err := yaml.Marshal(cd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// InstructionManager는 모든 InstructionData를 관리하는 구조체입니다.
type InstructionManager struct {
	commands map[string]InstructionData
}

// NewInstructionManager는 InstructionManager를 초기화하고 YAML 파일들을 읽어들입니다.
func NewInstructionManager() (*InstructionManager, error) {
	cm := &InstructionManager{commands: make(map[string]InstructionData)}

	err := cm.loadCommands()
	if err != nil {
		return nil, err
	}

	return cm, nil
}

// loadCommands는 주어진 경로에서 모든 YAML 파일을 읽어들여 InstructionData로 변환합니다.
func (cm *InstructionManager) loadCommands() error {
	// 디렉토리 내의 모든 YAML 파일을 찾습니다.
	files, err := filepath.Glob(filepath.Join("./HTTPsBAS-Procedures/", "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// 각 파일을 읽고 InstructionData로 변환
	for _, file := range files {
		err := cm.loadCommandFile(file)
		if err != nil {
			log.Printf("failed to load file %s: %v\n", file, err)
		}
	}

	return nil
}

// loadCommandFile은 하나의 YAML 파일을 읽어 InstructionData로 변환하고 저장합니다.
func (cm *InstructionManager) loadCommandFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 파일 내용 읽기
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// YAML 데이터를 InstructionData 구조체로 변환
	var command InstructionData
	err = yaml.Unmarshal(data, &command)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	// ID를 키로 맵에 저장
	cm.commands[command.ID] = command

	return nil
}

// GetByID는 주어진 ID에 해당하는 InstructionData를 반환합니다.
func (cm *InstructionManager) GetByID(id string) (*InstructionData, bool) {
	command, exists := cm.commands[id]
	if !exists {
		return nil, false
	}
	return &command, true
}

//// Insert는 새로운 InstructionData를 삽입하는 함수입니다.
//// 이미 동일한 ID가 존재하면 false를 반환하고, 성공적으로 삽입되면 true를 반환합니다.
//func (cm *InstructionManager) Insert(command InstructionData) bool {
//	// ID가 이미 존재하는지 확인
//	if _, exists := cm.commands[command.ID]; exists {
//		return false // 이미 존재하면 삽입하지 않고 false 반환
//	}
//
//	// 맵에 새 InstructionData 삽입
//	cm.commands[command.ID] = command
//	return true
//}

//func main() {
//	// ../CommandDB/ 디렉토리에 있는 모든 YAML 파일을 읽어 InstructionManager를 초기화합니다.
//	InstructionManager, err := NewInstructionManager("../CommandDB/")
//	if err != nil {
//		log.Fatalf("Failed to initialize command manager: %v", err)
//	}
//
//	// ID로 데이터를 가져오기 (예시)
//	id := "P_Collection_Kimsuky_001"
//	command, exists := InstructionManager.GetByID(id)
//	if exists {
//		fmt.Printf("Command found: %+v\n", command)
//	} else {
//		fmt.Printf("Command with ID %s not found\n", id)
//	}
//}
