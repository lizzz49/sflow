package sflow

type ActionDefinition struct {
	Definition
	AutoCommit  bool       `json:"auto_commit"`
	InvokerName string     `json:"invoker_name"` //action invoker func name
	Config      []KeyValue `json:"config"`
}

func newActionDefinition(id, name, invokerName string, autoCommit bool) *ActionDefinition {
	a := &ActionDefinition{}
	a.Id = id
	a.Name = name
	a.AutoCommit = autoCommit
	a.InvokerName = invokerName
	a.Config = []KeyValue{}
	return a
}
