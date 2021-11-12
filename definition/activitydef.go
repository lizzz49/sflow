package definition

import "github.com/lizzz49/sflow"

type ActivityDefinition struct {
	Definition
	Type        int                 `json:"type"`
	AutoCommit  bool                `json:"auto_commit"`
	Actions     []*ActionDefinition `json:"actions"`
	Participant UserDefinition      `json:"participant"`
	//Effective when the process is set to flat
	PreActivity *ActivityDefinition `json:"preActivity,omitempty"`
	actionCount int
}

func newActivityDefinition(id int, name string, autoCommit bool, pre *ActivityDefinition) *ActivityDefinition {
	a := &ActivityDefinition{
		Definition:  Definition{Id: id, Name: name},
		Type:        sflow.BusinessActivity,
		AutoCommit:  autoCommit,
		Actions:     []*ActionDefinition{},
		Participant: UserDefinition{},
		PreActivity: pre,
	}
	return a
}
func (a *ActivityDefinition) AddActionDefinition(name, invokerName string, autoCommit bool, role *RoleDefinition, pre *ActionDefinition) *ActionDefinition {
	a.actionCount++
	action := newActionDefinition(a.actionCount, name, invokerName, autoCommit, role, pre)
	a.Actions = append(a.Actions, action)
	return action
}
