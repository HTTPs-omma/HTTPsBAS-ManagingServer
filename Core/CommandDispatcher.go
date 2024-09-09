package Core

import (
	"fmt"
	"github.com/HTTPs-omma/HSProtocol/HSProtocol"
	"github.com/your/repo/Model"
)

// https://github.com/HTTPs-omma/HSProtocol
type CommandDispatcher struct {
}

// Command 상수를 정의
const (
	ACK                   = 0b0000000000
	UPDATE_HEALTH         = 0b0000000001
	UPDATE_PROTOCOL       = 0b0000000010
	POST_SYSTEM_INFO      = 0b0000000011
	RESERVED              = 0b0000000100
	POST_APPLICATION_INFO = 0b0000000101
	GET_PROCEDURE         = 0b0000000110
	POST_LOG_OF_PROCEDURE = 0b0000000111
)

func (cd *CommandDispatcher) action(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	// hsMgr := HSProtocol.NewHSProtocolManager()

	switch hs.Command {
	case UPDATE_HEALTH:
		return updateHealth(hs)
	case UPDATE_PROTOCOL:
		return updateProtocol(hs)
	case POST_SYSTEM_INFO:
		return postSystemInfo(hs)
	case RESERVED:
		break // 예약
	case POST_APPLICATION_INFO:
		return postApplicationInfo(hs)
	case GET_PROCEDURE:
		return getProcedure(hs)
	case POST_LOG_OF_PROCEDURE:
		return postLogOfProcedure(hs)
	}

	return nil, fmt.Errorf("Invalid Command")
}

// Command: 1 (0b0000000001)
func updateHealth(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	agsmd := Model.NewAgentStatusDB()
	rst, err := agsmd.ExistRecord()
	if err != nil {
		return nil, err
	}
	if rst {
		return nil, fmt.Errorf("Agent Status DB : no Records")
	}

	records, err := agsmd.SelectRecords()
	if err != nil {
		return nil, err
	}

	hs_uuid := string(hs.UUID[:])

	flag := false
	for _, record := range records {
		if record.UUID == hs_uuid {
			flag = true
		}
	}

	if flag == true {
		agsmd.UpdateRecord(&Model.AgentStatusRecord{
			UUID:   string(hs.UUID[:]),
			Status: Model.BinaryToAgentStatus(hs.HealthStatus),
		})
		return &HSProtocol.HS{ // ACK
			ProtocolID:     hs.ProtocolID,
			Command:        ACK,
			UUID:           hs.UUID,
			HealthStatus:   hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength:    hs.TotalLength,
			Data:           []byte{},
		}, nil
	} else if (flag == false) && (hs.HealthStatus == uint8(Model.Waiting)) {
		agsmd.InsertRecord(&Model.AgentStatusRecord{
			UUID:   string(hs.UUID[:]),
			Status: Model.BinaryToAgentStatus(hs.HealthStatus),
		})

		return &HSProtocol.HS{ // ACK
			ProtocolID:     hs.ProtocolID,
			Command:        ACK,
			UUID:           hs.UUID,
			HealthStatus:   hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength:    hs.TotalLength,
			Data:           []byte{},
		}, nil
	}

	return nil, fmt.Errorf("incorrect AgentStatusRecords")
}

// Command: 2 (0b0000000010)
func updateProtocol(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	// protocolID := binary.BigEndian.Uint32(hs.Data)
	agsmd := Model.NewAgentStatusDB()
	rst, err := agsmd.ExistRecord()
	if err != nil {
		return nil, err
	}
	if rst {
		return nil, fmt.Errorf("Agent Status DB : no Records")
	}

	records, err := agsmd.SelectRecords()
	if err != nil {
		return nil, err
	}
	hs_uuid := string(hs.UUID[:])

	for _, record := range records {
		if record.UUID == hs_uuid {
			record.Protocol = Model.BinaryToProtocol(hs.ProtocolID)
			err = agsmd.UpdateRecord(&record)
			if err != nil {
				return nil, err
			}
		}
	}

	return &HSProtocol.HS{ // ACK
		ProtocolID:     hs.ProtocolID,
		Command:        ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 3 (0b0000000011)
func postSystemInfo(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	sysDB := Model.NewSystemInfoDB()
	Dsys, err := sysDB.Unmarshal(hs.Data)
	if err != nil {
		return nil, err
	}

	err = sysDB.InsertRecord(Dsys)
	if err != nil {
		return nil, err
	}

	return &HSProtocol.HS{ // ACK
		ProtocolID:     hs.ProtocolID,
		Command:        ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 5 (0b0000000101)
func postApplicationInfo(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	appDB := Model.ApplicationDB{}
	Dapp, err := appDB.Unmarshal(hs.Data)
	if err != nil {
		return nil, err
	}
	err = appDB.InsertRecord(Dapp)
	if err != nil {
		return nil, err
	}

	return &HSProtocol.HS{ // ACK
		ProtocolID:     hs.ProtocolID,
		Command:        ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 6 (0b0000000110)
func getProcedure(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	appDB := Model.ApplicationDB{}
	Dapp, err := appDB.Unmarshal(hs.Data)
	if err != nil {
		return nil, err
	}
	err = appDB.InsertRecord(Dapp)
	if err != nil {
		return nil, err
	}

	return &HSProtocol.HS{ // ACK
		ProtocolID:     hs.ProtocolID,
		Command:        ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 7 (0b0000000111)
func postLogOfProcedure(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	jbMgr := NewJobManager() // stack 에 있어야겠지?

	var agentUuid string
	copy(hs.UUID[:], agentUuid)

	job, exist := jbMgr.GetData(agentUuid) // agentuuid 에 해당하는 값 찾아괴

	if exist == true { // job 이 있다면
		// procedureID 에 맵핑하여 yaml 파일을 직렬화하고 불러와서

		return &HSProtocol.HS{ // ACK
			ProtocolID:     hs.ProtocolID,
			Command:        ACK,
			UUID:           hs.UUID,
			HealthStatus:   hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength:    hs.TotalLength,
			Data:           []byte{},
		}, nil
	}

	// false
	return &HSProtocol.HS{ // ACK
		ProtocolID:     hs.ProtocolID,
		Command:        ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}
