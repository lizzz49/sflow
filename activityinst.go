package sflow

type ActivityInstance struct {
	Process    *ProcessInstance
	Definition *ActivityDefinition
	Actions    []*ActionInstance
	status     int
}

func (ai *ActivityInstance) Start() error {
	ai.status = ProcessInstanceStatusStarted
	if ai.Definition.IsStart || ai.Definition.IsEnd {
		return ai.Process.FinishActivity(ai)
	}
	for _, ad := range ai.Definition.Actions {
		action := ai.CreateActionInstance(ad)
		action.Start()
	}
	return nil
}

func (ai *ActivityInstance) CreateActionInstance(ad *ActionDefinition) *ActionInstance {
	a := &ActionInstance{Activity: ai, Definition: ad, Status: ProcessInstanceStatusNew}
	ai.Actions = append(ai.Actions, a)
	return a
}

func (ai *ActivityInstance) FinishAction(a *ActionInstance) error {
	a.Status = ProcessInstaneStatusFinish
	allFinish := true
	for _, action := range ai.Actions {
		if action.Status != ProcessInstaneStatusFinish {
			allFinish = false
			break
		}
	}
	if allFinish {
		return ai.Process.FinishActivity(ai)
	}
	return nil
}
