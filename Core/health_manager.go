package Core

// 이거 보세요 -> https://d2.naver.com/helloworld/8588537
import (
	"encoding/json"
	"github.com/your/repo/Model"
	"net"
	"os"
	"strconv"
)

type HealthMsg struct {
	UUID   string `json:"uuid"`
	Status int    `json:"status"`
}

/*
Let Check Example
refer : https://gist.github.com/miguelmota/01ba5131838ae31947ac9b03e57f3773
*/
func Heartbeat() {

	// UDP 열기 : HealthCheck
	port, err := strconv.Atoi(os.Getenv("PORT"))
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			message := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(message[:])
			if err != nil {
				panic(err)
			}
			if n == 0 {
				continue
			}

			hmsg := &HealthMsg{}
			// json 역질렬화, Unmarshal
			err = json.Unmarshal(message[:n], &hmsg)
			if err != nil {
				panic(err)
			}

			agtstatdb := Model.NewAgentStatusDB()

			agtstatRcrd := Model.AgentStatusRecord{
				UUID:   hmsg.UUID,
				Status: getStatusType(hmsg.Status),
			}

			err = agtstatdb.UpdateRecord(&agtstatRcrd)
			if err != nil {
				panic(err) // 에러처리를 바꿔야함
			}
		}

		defer conn.Close()
	}()

}

func getStatusType(status int) Model.AgentStatus {

	switch status {
	case 0:
		return Model.Running
	case 1:
		return Model.Waiting
	case 2:
		return Model.Stopping
	}
	return Model.Stopping
}
