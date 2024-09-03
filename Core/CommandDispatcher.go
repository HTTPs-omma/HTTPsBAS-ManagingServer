package Core

import (
	"fmt"
	"github.com/HTTPs-omma/HSProtocol/HSProtocol"
	"github.com/your/repo/Model"
)
// https://github.com/HTTPs-omma/HSProtocol

type CommandDispatcher struct {

}


func (cd *CommandDispatcher)action(hs *HSProtocol.HS) (*HSProtocol.HS, error) {
	//hsMgr := HSProtocol.NewHSProtocolManager()

	switch hs.Command {
	case 0b0000000001 :
		return updateHealth(hs)
	case 0b0000000010 :
		updateProtocol(hs)
	case 0b0000000011 :
		postSystemInfo(hs)
	case 0b0000000100 :
		break; // 예약
	case 0b0000000101 :
		postApplicationInfo(hs)
	case 0b0000000110 :
		getProcedure(hs)
	case 0b0000000111 :
		postLogOfProcedure(hs)
	}

	return nil, fmt.Errorf("Invalid Command")
}

// Command: 1 (0b0000000001)
func updateHealth(hs *HSProtocol.HS)(*HSProtocol.HS, error){
	agsmd := Model.NewAgentStatusDB()
	rst, err := agsmd.ExistRecord()
	if err != nil {
		return nil, err
	}
	if( rst ) {
		return nil, fmt.Errorf("Agent Status DB : no Records")
	}



	records, err := agsmd.SelectRecords()
	if err != nil {
		return nil, err
	}

	hs_uuid := string(hs.UUID[:])


	flag := false
	for _, record  := range records {
		if record.UUID == hs_uuid {
			flag = true
		}
	}

	if flag == true {
		agsmd.UpdateRecord(&Model.AgentStatusRecord{
			UUID : string(hs.UUID[:]),
			Status: Model.BinaryToAgentStatus(hs.HealthStatus),
		})
		return &HSProtocol.HS{ // ACK
			ProtocolID: hs.ProtocolID,
			Command: 0b0000000000,
			UUID: hs.UUID,
			HealthStatus: hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength: hs.TotalLength,
			Data: []byte{},
		}, nil
	} else if (flag == false) && ( hs.HealthStatus == uint8(Model.Waiting) )  {
		agsmd.InsertRecord(&Model.AgentStatusRecord{
			UUID : string(hs.UUID[:]),
			Status: Model.BinaryToAgentStatus(hs.HealthStatus),
		})

		return &HSProtocol.HS{ // ACK
			ProtocolID : hs.ProtocolID,
			Command: 0b0000000000,
			UUID: hs.UUID,
			HealthStatus: hs.HealthStatus,
			Identification: hs.Identification,
			TotalLength: hs.TotalLength,
			Data: []byte{},
		}, nil
	}


	return nil, fmt.Errorf("incorrect AgentStatusRecords")
}

// Command: 2 (0b0000000010)
func updateProtocol(hs *HSProtocol.HS)(*HSProtocol.HS, error){

	// protocolID := binary.BigEndian.Uint32(hs.Data)
	agsmd := Model.NewAgentStatusDB()
	rst, err := agsmd.ExistRecord()
	if err != nil {
		return nil, err
	}
	if( rst ) {
		return nil, fmt.Errorf("Agent Status DB : no Records")
	}


	records, err := agsmd.SelectRecords()
	if err != nil {
		return nil, err
	}
	hs_uuid := string(hs.UUID[:])

	for _, record  := range records {
		if record.UUID == hs_uuid {
			record.Protocol = Model.BinaryToProtocol( hs.ProtocolID )
			err = agsmd.UpdateRecord(&record)
			if err != nil {
				return nil, err
			}
		}
	}


	return &HSProtocol.HS{ // ACK
		ProtocolID: hs.ProtocolID,
		Command: 0b0000000000,
		UUID: hs.UUID,
		HealthStatus: hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength: hs.TotalLength,
		Data: []byte{},
	}, nil
}


// Command: 3 (0b0000000011)
func postSystemInfo (hs *HSProtocol.HS)(*HSProtocol.HS, error){


	return &HSProtocol.HS{ // ACK
		ProtocolID: hs.ProtocolID,
		Command: 0b0000000000,
		UUID: hs.UUID,
		HealthStatus: hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength: hs.TotalLength,
		Data: []byte{},
	}, nil
}


// Command: 5 (0b0000000101)
func postApplicationInfo (hs *HSProtocol.HS)(*HSProtocol.HS, error){




	return &HSProtocol.HS{ // ACK
		ProtocolID: hs.ProtocolID,
		Command: 0b0000000000,
		UUID: hs.UUID,
		HealthStatus: hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength: hs.TotalLength,
		Data: []byte{},
	}, nil
}

// Command: 5 (0b0000000110)
func getProcedure (hs *HSProtocol.HS)(*HSProtocol.HS, error){




	return &HSProtocol.HS{ // ACK
		ProtocolID: hs.ProtocolID,
		Command: 0b0000000000,
		UUID: hs.UUID,
		HealthStatus: hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength: hs.TotalLength,
		Data: []byte{},
	}, nil
}

// Command: 5 (0b0000000110)
func postLogOfProcedure (hs *HSProtocol.HS)(*HSProtocol.HS, error){




	return &HSProtocol.HS{ // ACK
		ProtocolID : hs.ProtocolID,
		Command: 0b0000000000,
		UUID: hs.UUID,
		HealthStatus: hs.HealthStatus,
		Identification: hs.Identification,
		TotalLength: hs.TotalLength,
		Data: []byte{},
	}, nil
}


