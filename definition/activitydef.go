package definition

import "github.com/lizzz49/sflow"

type ActivityDefinition struct {
	Definition
	Type        int                 `json:"type"`
	AutoCommit  bool                `json:"auto_commit"`
	Actions     []*ActionDefinition `json:"actions"`
	Participant UserDefinition      `json:"participant"`
	actionCount int
}

func newActivityDefinition(id int, name string, autoCommit bool) *ActivityDefinition {
	a := &ActivityDefinition{
		Definition:  Definition{Id: id, Name: name},
		Type:        sflow.BusinessActivity,
		AutoCommit:  autoCommit,
		Actions:     []*ActionDefinition{},
		Participant: UserDefinition{},
	}
	return a
}
func (a *ActivityDefinition) AddActionDefinition(name, invokerName string, autoCommit bool) *ActionDefinition {
	a.actionCount++
	action := newActionDefinition(a.actionCount, name, invokerName, autoCommit)
	a.Actions = append(a.Actions, action)
	return action
}
