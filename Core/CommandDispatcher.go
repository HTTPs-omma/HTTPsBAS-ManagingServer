package Core

import (
	"encoding/json"
	"fmt"

	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/your/repo/Model"
)

// https://github.com/HTTPs-omma/HSProtocol
type CommandDispatcher struct {
}

// Command 상수를 정의

func (cd *CommandDispatcher) Action(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	// hsMgr := HSProtocol.NewHSProtocolManager()

	switch hs.Command {
	case HSProtocol.UPDATE_AGENT_PROTOCOL:
		return UPDATE_AGENT_PROTOCOL(hs)
	case HSProtocol.UPDATE_AGENT_STATUS:
		return UPDATE_AGENT_STATUS(hs)
	case HSProtocol.SEND_AGENT_SYS_INFO:
		return SEND_AGENT_SYS_INFO(hs)
	case HSProtocol.ERROR_ACK:
		break // 예약
	case HSProtocol.SEND_AGENT_APP_INFO:
		return SEND_AGENT_APP_INFO(hs)
	case HSProtocol.FETCH_INSTRUCTION:
		return FETCH_INSTRUCTION(hs)
	case HSProtocol.SEND_PROCEDURE_LOG:
		return SEND_PROCEDURE_LOG(hs)
	}

	return nil, fmt.Errorf("Invalid Command")
}

// Command: 1 (0b0000000001)
func UPDATE_AGENT_PROTOCOL(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	agsmd, err := Model.NewAgentStatusDB()
	if err != nil {
		return nil, err
	}
	rst, err := agsmd.ExistRecord()
	if err != nil {
		return nil, err
	}
	if rst {
		return nil, fmt.Errorf("Agent Status DB : no Records")
	}

	records, err := agsmd.SelectAllRecords()
	if err != nil {
		return nil, err
	}

	hs_uuid := HSProtocol.ByteArrayToHexString(hs.UUID)

	flag := false
	for _, record := range records {
		if record.UUID == hs_uuid {
			flag = true
		}
	}

	if flag == true {
		agsmd.UpdateRecord(&Model.AgentStatusRecord{
			UUID:   hs_uuid,
			Status: hs.HealthStatus,
		})
		return &HSProtocol.HS{ // HSProtocol.ACK
			ProtocolID:     hs.ProtocolID,
			Command:        HSProtocol.ACK,
			UUID:           hs.UUID,
			HealthStatus:   hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength:    hs.TotalLength,
			Data:           []byte{},
		}, nil
	} else if (flag == false) && (hs.HealthStatus == HSProtocol.WAIT) {
		agsmd.InsertRecord(&Model.AgentStatusRecord{
			UUID:   hs_uuid,
			Status: hs.HealthStatus,
		})

		return &HSProtocol.HS{ // HSProtocol.ACK
			ProtocolID:     hs.ProtocolID,
			Command:        HSProtocol.ACK,
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
func UPDATE_AGENT_STATUS(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	// protocolID := binary.BigEndian.Uint32(hs.Data)
	agsDb, err := Model.NewAgentStatusDB()
	if err != nil {
		return nil, err
	}

	err = agsDb.InsertRecord(&Model.AgentStatusRecord{
		ID:       0,
		UUID:     HSProtocol.ByteArrayToHexString(hs.UUID),
		Protocol: hs.ProtocolID,
		Status:   hs.HealthStatus,
	})
	// err = agsDb.InsertRecord(&Model.AgentStatusRecord{})
	if err != nil {
		return nil, err
	}

	return &HSProtocol.HS{ // HSProtocol.ACK
		ProtocolID:     hs.ProtocolID,
		Command:        HSProtocol.ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 3 (0b0000000011)
func SEND_AGENT_SYS_INFO(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	sysDB, err := Model.NewSystemInfoDB()
	if err != nil {
		return nil, err
	}
	sysinfo := Model.DsystemInfoDB{}
	err = json.Unmarshal(hs.Data, &sysinfo)
	// fmt.Println("SEND_AGENT_SYS_INFO UUID : " + sysinfo.Uuid)
	if err != nil {
		return nil, err
	}
	if err = sysDB.DeleteRecordByUUID(sysinfo.Uuid); err != nil {
		return nil, err
	}
	err = sysDB.InsertRecord(&sysinfo)
	if err != nil {
		return nil, err
	}

	return &HSProtocol.HS{ // HSProtocol.ACK
		ProtocolID:     hs.ProtocolID,
		Command:        HSProtocol.ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 5 (0b0000000101)
//func SEND_AGENT_APP_INFO(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
//
//	appDB, err := Model.NewApplicationDB()
//	if err != nil {
//		return nil, err
//	}
//	applist, err := appDB.FromJSON(hs.Data)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, Dapp := range applist {
//		err = appDB.InsertRecord(&Dapp)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &HSProtocol.HS{ // HSProtocol.ACK
//		ProtocolID:     hs.ProtocolID,
//		Command:        HSProtocol.ACK,
//		UUID:           hs.UUID,
//		HealthStatus:   hs.HealthStatus,
//		Identification: hs.Identification,
//		TotalLength:    hs.TotalLength,
//		Data:           []byte{},
//	}, nil
//}

func SEND_AGENT_APP_INFO(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	appDB, err := Model.NewProgramsDB()
	if err != nil {
		return nil, err
	}
	applist, err := appDB.FromJSON(hs.Data)
	if err != nil {
		return nil, err
	}

	for _, Dapp := range applist {
		err = appDB.InsertRecord(Dapp.AgentUUID, Dapp.FileName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return &HSProtocol.HS{ // HSProtocol.ACK
		ProtocolID:     hs.ProtocolID,
		Command:        HSProtocol.ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}

// Command: 6 (0b0000000110)
func FETCH_INSTRUCTION(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	agentUuid := HSProtocol.ByteArrayToHexString(hs.UUID)
	//fmt.Println("agent uuid : " + agentUuid)
	jobdb, err := Model.NewJobDB()
	if err != nil {
		return nil, err
	}
	job, err, exist := jobdb.PopbyAgentUUID(agentUuid)

	if err != nil {
		return nil, err
	}

	//fmt.Println("debug === : " + job.ProcedureID)
	if exist == true { // job 이 있다면
		cmdMgr, err := NewInstructionManager()
		if err != nil {
			return nil, err
		}
		cmdData, issuccess := cmdMgr.GetByID(job.ProcedureID) // 프로시저를 불러와야함.
		if issuccess != true {
			if job.Action == "GetSystemInfo" || job.Action == "GetApplication" || job.Action == "StopAgent" || job.Action == "ChangeProtocolToTCP" || job.Action == "ChangeProtocolToHTTP" {
				cmdData, issuccess = cmdMgr.GetByID("P_Collection_0001")
			} else {
				return nil, fmt.Errorf("job procedure not found")
			}

		}

		extendedcmdData := cmdData.ConvertToExtended(job.MessageUUID, job.Action, job.Files)
		bData, err := extendedcmdData.ToBytes()
		if err != nil {
			return nil, err
		}
		return &HSProtocol.HS{ // HSProtocol.ACK
			ProtocolID:     hs.ProtocolID,
			Command:        HSProtocol.ACK,
			UUID:           hs.UUID,
			HealthStatus:   hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength:    hs.TotalLength,
			Data:           bData,
		}, nil
	}

	// false
	return &HSProtocol.HS{ // HSProtocol.ACK
		ProtocolID:     hs.ProtocolID,
		Command:        HSProtocol.ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil

}

// Command: 7 (0b0000000111)
func SEND_PROCEDURE_LOG(hs *HSProtocol.HS) (*HSProtocol.HS, error) {

	//hs_uuid := HSProtocol.ByteArrayToHexString(hs.UUID)

	logdb, err := Model.NewOperationLogDB()
	if err != nil {
		return nil, err
	}
	log := &Model.OperationLogDocument{}
	err = json.Unmarshal(hs.Data, &log)
	if err != nil {
		return nil, err
	}

	_, err = logdb.InsertDocument(log)
	if err != nil {
		return nil, err
	}

	return &HSProtocol.HS{ // HSProtocol.ACK
		ProtocolID:     hs.ProtocolID,
		Command:        HSProtocol.ACK,
		UUID:           hs.UUID,
		HealthStatus:   hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength:    hs.TotalLength,
		Data:           []byte{},
	}, nil
}
