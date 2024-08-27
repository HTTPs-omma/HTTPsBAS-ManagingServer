package Core

import "github.com/your/repo/Model"

/*
	Agent Model = Managing + model
*/

type AgentManager struct {
}

func (self *AgentManager) createAgent(uuid string) error {

	return nil
}

/*
agent 의 상태 Status 는 총 3가지로 나뉜다.

Running : 동작중인 상태
waiting : 대기중인 상태
Stopping : 정지후 사라지기 전에 상태
*/
func (self *AgentManager) checkAgent(uuid string) bool {

	return false
}

func (self *AgentManager) updateAgentStatus(uuid string, status bool) bool {
	agtStat := &Model.NewAgentStatusDB()
	record := &Model.AgentStatusRecord{
		UUID:   uuid,
		Status: status,
	}
	agtStat.UpdateRecord()
}

func (self *AgentManager) deleteAgent(uuid string) bool {

}
