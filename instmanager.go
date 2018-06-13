package sflow

import "fmt"

type ProcessInstanceManager struct {
	process      map[string]*ProcessInstance
	maxProcessId int
}

func NewProcessInstanceManager() *ProcessInstanceManager {
	return &ProcessInstanceManager{process: make(map[string]*ProcessInstance)}
}

func (pim *ProcessInstanceManager) CreateProcessInstance(pd *ProcessDefinition) *ProcessInstance {
	pim.maxProcessId++
	pi := newProcessInstance(fmt.Sprintf("%d", pim.maxProcessId), pd)
	pim.process[pi.id] = pi
	return pi
}

func (pim *ProcessInstanceManager) TerminateProcessInstance(id string) bool {
	pi := pim.process[id]
	if pi == nil {
		return true
	}
	status := pi.status
	if status == ProcessInstanceStatusNew {
		status = ProcessInstanceStatusTerminated
	} else if status == ProcessInstanceStatusStarted || status == ProcessInstanceStatusSuspended {
		pi.status = ProcessInstanceStatusTerminated
		pi.finish <- ProcessInstanceStatusTerminated
	} else {
		return true
	}

	return true
}
