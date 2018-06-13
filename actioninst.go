package sflow

import (
	"fmt"
)

var (
	actions = make(map[string]ActionInvoker)
)

type ActionInstance struct {
	Activity   *ActivityInstance
	Definition *ActionDefinition
	Status     int
}
type ActionInvoker struct {
	Name       string
	InvokeFunc InvokeFunc
}
type InvokeFunc func(*ProcessContext) bool

func RegistryAction(name string, ai InvokeFunc) {
	if actions == nil {
		actions = make(map[string]ActionInvoker)
	}
	if actions[name].InvokeFunc != nil {
		panic(fmt.Sprintf("the action allready exist with the name %s", name))
	}
	actions[name] = ActionInvoker{name, ai}
}

func GetActionInvoker(name string) (ai ActionInvoker, has bool) {
	ai = actions[name]
	return ai, ai.InvokeFunc != nil
}

func ListActionInvoker() (ais []ActionInvoker) {
	for _, ai := range actions {
		ais = append(ais, ai)
	}
	return
}
func (ai *ActionInstance) Start() error {
	ai.Status = ProcessInstanceStatusStarted
	invoker, has := GetActionInvoker(ai.Definition.InvokerName)
	if has {
		invokeFunc := invoker.InvokeFunc
		process := ai.Activity.Process
		invokeFunc(&process.context)
		if ai.Definition.AutoCommit {
			ai.Activity.FinishAction(ai)
		}
	} else {
		ai.Activity.FinishAction(ai)
	}
	return nil
}
