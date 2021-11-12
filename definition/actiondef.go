package definition

type ActionDefinition struct {
	Definition
	AutoCommit  bool            `json:"auto_commit"`
	InvokerName string          `json:"invoker_name"` //action invoker func name
	Role        *RoleDefinition `json:"role"`
	PreAction   *ActionDefinition
}

func newActionDefinition(id int, name, invokerName string, autoCommit bool, role *RoleDefinition, pre *ActionDefinition) *ActionDefinition {
	a := &ActionDefinition{
		Definition:  Definition{Id: id, Name: name},
		AutoCommit:  autoCommit,
		InvokerName: invokerName,
		Role:        role,
		PreAction:   pre,
	}
	return a
}
