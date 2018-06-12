package sflow

import "fmt"

type ActivityDefinition struct {
	Definition
	IsStart    bool `json:"is_start"`
	IsEnd      bool `json:"is_end"`
	AutoCommit bool `json:"auto_commit"`
	//From        []string           `json:"from"` //input activities id
	//To          []string           `json:"to"`   //output activities id
	Actions     []*ActionDefinition `json:"actions"`
	Participant UserDefinition      `json:"participant"`
	actionCount int
}

func newActivityDefinition(id,name string,autoCommit bool)*ActivityDefinition{
	a := &ActivityDefinition{}
	a.Id = id
	a.Name = name
	a.AutoCommit = autoCommit
	a.Actions = []*ActionDefinition{}
	return a
}
func (a *ActivityDefinition)AddActionDefinition(name,invokerName string,autoCommit bool)*ActionDefinition{
	a.actionCount++
	action := newActionDefinition(fmt.Sprintf("%d",a.actionCount),name,invokerName,autoCommit)
	a.Actions = append(a.Actions,action)
	return action
}