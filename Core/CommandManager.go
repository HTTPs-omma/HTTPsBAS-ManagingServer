package Core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/**
ChatGpt 로 생성한 코드임
*/

// CommandData는 주어진 YAML 데이터를 저장할 구조체입니다.
type CommandData struct {
	ID               string `yaml:"id"`
	MITREID          string `yaml:"MITRE_ID"`
	Description      string `yaml:"Description"`
	Tool             string `yaml:"tool"`
	RequisiteCommand string `yaml:"requisite_command"`
	Command          string `yaml:"command"`
	Cleanup          string `yaml:"cleanup"`
}

// ToBytes는 CommandData 구조체를 YAML 바이트 슬라이스로 변환하는 함수입니다.
func (cd *CommandData) ToBytes() ([]byte, error) {
	// YAML로 직렬화
	data, err := yaml.Marshal(cd)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CommandManager는 모든 CommandData를 관리하는 구조체입니다.
type CommandManager struct {
	commands map[string]CommandData
}

// NewCommandManager는 CommandManager를 초기화하고 YAML 파일들을 읽어들입니다.
func NewCommandManager() (*CommandManager, error) {
	cm := &CommandManager{commands: make(map[string]CommandData)}

	err := cm.loadCommands()
	if err != nil {
		return nil, err
	}

	return cm, nil
}

// loadCommands는 주어진 경로에서 모든 YAML 파일을 읽어들여 CommandData로 변환합니다.
func (cm *CommandManager) loadCommands() error {
	// 디렉토리 내의 모든 YAML 파일을 찾습니다.
	files, err := filepath.Glob(filepath.Join("../CommandDB/", "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// 각 파일을 읽고 CommandData로 변환
	for _, file := range files {
		err := cm.loadCommandFile(file)
		if err != nil {
			log.Printf("failed to load file %s: %v\n", file, err)
		}
	}

	return nil
}

// loadCommandFile은 하나의 YAML 파일을 읽어 CommandData로 변환하고 저장합니다.
func (cm *CommandManager) loadCommandFile(filepath string) error {
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

	// YAML 데이터를 CommandData 구조체로 변환
	var command CommandData
	err = yaml.Unmarshal(data, &command)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	// ID를 키로 맵에 저장
	cm.commands[command.ID] = command

	return nil
}

// GetByID는 주어진 ID에 해당하는 CommandData를 반환합니다.
func (cm *CommandManager) GetByID(id string) (*CommandData, bool) {
	command, exists := cm.commands[id]
	if !exists {
		return nil, false
	}
	return &command, true
}

//func main() {
//	// ../CommandDB/ 디렉토리에 있는 모든 YAML 파일을 읽어 CommandManager를 초기화합니다.
//	commandManager, err := NewCommandManager("../CommandDB/")
//	if err != nil {
//		log.Fatalf("Failed to initialize command manager: %v", err)
//	}
//
//	// ID로 데이터를 가져오기 (예시)
//	id := "P_Collection_Kimsuky_001"
//	command, exists := commandManager.GetByID(id)
//	if exists {
//		fmt.Printf("Command found: %+v\n", command)
//	} else {
//		fmt.Printf("Command with ID %s not found\n", id)
//	}
//}
