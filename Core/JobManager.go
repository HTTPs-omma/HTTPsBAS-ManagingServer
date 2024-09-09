package Core

import (
	"time"
)

type JobData struct {
	TechnicalID string `json:"technicalID"`
	AgentUUID   string `json:"agentUUID"`
	MessageUUID string `json:"messageUUID"`
	createAt    time.Time
}

type JobManager struct {
	dataMap map[string][]JobData
}

func NewJobManager() *JobManager {
	return &JobManager{
		dataMap: make(map[string][]JobData),
	}
}

func popFront(slice []JobData) (JobData, []JobData) {
	if len(slice) > 0 {
		job := slice[0]
		return job, slice[1:]
	}
	return slice[0], slice
}

func (jm *JobManager) GetData(agentUUID string) (*JobData, bool) {
	jobs, exists := jm.dataMap[agentUUID]

	if len(jobs) > 0 {
		var job JobData
		job, jm.dataMap[agentUUID] = popFront(jobs)
		return &job, exists
	}

	return &JobData{}, exists
}

func (jm *JobManager) insertData(jobData JobData) bool {
	jm.dataMap[jobData.AgentUUID] = append(jm.dataMap[jobData.AgentUUID], jobData)

	return true
}
