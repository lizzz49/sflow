package sflow

import (
	"fmt"
)

const (
	StatusNew = iota
	StatusStarted
	StatusSuspended
	StatusError
	StatusFinish
)

type ProcessInstance struct {
	id string
	definition            *ProcessDefinition
	context               ProcessContext
	activityInstances     []*ActivityInstance
	activityDefinitions   map[string]*ActivityDefinition
	transitionDefinitions map[string]*TransitionDefinition
	transmits             map[string]*TransitionInstance
	status                int
	finish                chan int
}

func newProcessInstance(id string,def *ProcessDefinition) *ProcessInstance {
	des,ok:=def.Check()
	if ok{
		des, ok = checkProcessDefinition(def)
	}

	if ok {
		pi := &ProcessInstance{
			id:id,
			definition:        def,
			context:           make(map[string]Value),
			status:            StatusNew,
			activityInstances: []*ActivityInstance{},
			transmits:make(map[string]*TransitionInstance),
			finish:make(chan int),
		}
		as := map[string]*ActivityDefinition{}
		for _,ad:=range def.Activities{
			as[ad.Id] = ad
		}
		as[def.StartActivity.Id] = &def.StartActivity
		as[def.EndActivity.Id] = &def.EndActivity
		pi.activityDefinitions = as

		ts := map[string]*TransitionDefinition{}
		for _,td:=range def.Transitions{
			ts[td.Id] = td
		}
		pi.transitionDefinitions = ts
		for _,transmit:=range def.Transitions {
			pi.createTransition(transmit)
		}
		return pi
	} else {
		var errs string
		for i, err := range des {
			errs += fmt.Sprintf("%d\t%s\n", i, err.Error())
		}
		panic(errs)
	}
}
func checkProcessDefinition(p *ProcessDefinition)([]DefError,bool){
	var des []DefError
	if p.Status != PDPublishStatus{
		des = append(des, DefError{ErrMsg: fmt.Sprintf("process [%s] not in published [%d] status,status: %d",p.Id,PDPublishStatus,p.Status)})
	}
	for _,node:=range p.Activities{
			for _,action:=range node.Actions {
				_, has := GetActionInvoker(action.InvokerName)
				if !has {
					des = append(des, DefError{ErrMsg: fmt.Sprintf(ActionNotFound, action.InvokerName, node.Id)})
				}
			}

	}
	return des,des==nil||len(des)==0
}
func(pi *ProcessInstance)Init(ctx ProcessContext){
	pi.context = ctx
}
func(pi *ProcessInstance) createActivityInstance(def *ActivityDefinition)*ActivityInstance{
	st := &ActivityInstance{Process:pi,Definition:def, status:StatusNew}
	pi.activityInstances = append(pi.activityInstances, st)
	return st
}
func (pi *ProcessInstance) createTransition(def *TransitionDefinition)*TransitionInstance{
	t:=&TransitionInstance{pi,def,false}
	pi.transmits[def.Id] = t
	return t
}
func (pi *ProcessInstance) Start() error {
	pi.status = StatusStarted
	st := pi.createActivityInstance(&pi.definition.StartActivity)
	st.Start()
	return nil
}
func (pi *ProcessInstance)FinishActivity(ai *ActivityInstance)error{
	ai.status = StatusFinish
	if ai.Definition.IsEnd{
		pi.status = StatusFinish
		pi.finish<-0
	}
	tos := pi.findNext(ai)
	if len(tos)==0&&pi.allActivityFinish(){
		pi.finish<-1
	}
	for _,to:=range tos{
		adef:=pi.activityDefinitions[to]
		next := pi.createActivityInstance(adef)
		go next.Start()
	}
	return nil
}
func (pi *ProcessInstance)allActivityFinish()bool{
	for _,a:=range pi.activityInstances{
		if a.status !=StatusFinish{
			return false
		}
	}
	return true
}
func (pi *ProcessInstance)findNext(ai *ActivityInstance)(ads []string){
	ads = []string{}
	var trans []string
	for _,td:=range pi.definition.Transitions{
		if td.From == ai.Definition.Id{
			trans = append(trans,td.Id)
		}
	}
	var next []string
	for _,t:=range trans{
		ti := pi.transmits[t]
		if ti.Check(){
			ti.Done = true
			next = append(next,ti.Definition.To)
		}
	}
	if next == nil || len(next)==0{
		return
	}
	for _,nid:=range next{
		allFromOk := true
		for _,ti:=range pi.transmits{
			if ti.Definition.To==nid && !ti.Done{
				allFromOk = false
				break
			}
		}
		if allFromOk{
			ads = append(ads,nid)
		}
	}
	return
}
func (pi *ProcessInstance) Suspend() {
	pi.status = StatusSuspended
}

func (pi *ProcessInstance) Wait()int{
	return <-pi.finish
}