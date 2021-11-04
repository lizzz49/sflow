package instance

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/definition"
	"log"
	"time"
)

type ProcessInstance struct {
	Id         int        //`gorm:"primaryKey;autoIncrement;column:id;type:int(11);comment:实例Id"`
	Name       string     `gorm:"column:name;type:varchar(255);comment:流程实例名称"`
	Definition int        `gorm:"column:definition;type:int(11);comment:定义Id"`
	Status     int        `gorm:"column:Status;type:int(3);comment:流程状态"`
	CreateTime time.Time  `gorm:"column:create_time;comment:流程创建时间"`
	StartTime  *time.Time `gorm:"column:start_time;comment:流程启动时间"`
	FinishTime *time.Time `gorm:"column:finish_time;comment:流程完成时间"`
	//definition            *definition.ProcessDefinition
	context sflow.ProcessContext
	//activityInstances     []*ActivityInstance
	//activityDefinitions   map[int]*definition.ActivityDefinition
	//transitionDefinitions map[int]*definition.TransitionDefinition
	//transitions           map[int]*TransitionInstance
	//finish                chan int
	ModelTime
	ControlBy
}

func (ProcessInstance) TableName() string {
	return "sflow_process"
}
func checkProcessDefinition(p *definition.ProcessDefinition) ([]definition.DefError, bool) {
	var des []definition.DefError
	if p.Status != definition.PDPublishStatus {
		des = append(des, definition.DefError{ErrMsg: fmt.Sprintf("process [%d] not in published [%d] Status,Status: %d", p.Id, definition.PDPublishStatus, p.Status)})
	}
	for _, node := range p.Activities {
		for _, action := range node.Actions {
			_, has := GetActionInvoker(action.InvokerName)
			if !has {
				des = append(des, definition.DefError{ErrMsg: fmt.Sprintf(definition.ActionNotFound, action.InvokerName, node.Id)})
			}
		}

	}
	return des, des == nil || len(des) == 0
}
func (pi *ProcessInstance) Init(ctx sflow.ProcessContext) {
	pi.context = ctx
}
func (pi *ProcessInstance) createActivityInstance(def *definition.ActivityDefinition) (ai *ActivityInstance, err error) {
	ai = &ActivityInstance{
		ProcessId:         pi.Id,
		Name:              def.Name,
		ProcessDefinition: pi.Definition,
		Definition:        def.Id,
		Type:              def.Type,
		AutoCommit:        def.AutoCommit,
		Status:            sflow.ProcessInstanceStatusNew,
		CreateTime:        time.Now(),
	}
	//pi.activityInstances = append(pi.activityInstances, st)
	rs := manager.db.Model(&ActivityInstance{}).Create(ai)
	return ai, rs.Error
}
func (pi *ProcessInstance) createTransition(def *definition.TransitionDefinition) (t *TransitionInstance, err error) {
	t = &TransitionInstance{
		Name:              def.Name,
		ProcessId:         pi.Id,
		ProcessDefinition: pi.Definition,
		Definition:        def.Id,
		Done:              false,
	}
	rs := manager.db.Model(&TransitionInstance{}).Create(t)
	if rs.Error != nil {
		return nil, rs.Error
	}
	//pi.transitions[def.Id] = t
	return t, nil
}
func (pi *ProcessInstance) Start() error {
	pi.Status = sflow.ProcessInstanceStatusStarted
	def, has := manager.pdm.GetProcessDefinitionById(pi.Definition)
	if !has {
		return fmt.Errorf("process definition %d not found", pi.Definition)
	}
	st, err := pi.createActivityInstance(&def.StartActivity)
	if err != nil {
		return err
	}
	t := time.Now()
	rs := manager.db.Model(&ProcessInstance{}).Where("id = ?", pi.Id).Updates(ProcessInstance{
		Status:    sflow.ProcessInstaneStatusFinish,
		StartTime: &t,
	})
	if rs.Error != nil {
		return rs.Error
	}
	return st.Start()
}
func (pi *ProcessInstance) FinishActivity(ai *ActivityInstance) error {
	ai.Status = sflow.ProcessInstaneStatusFinish
	if ai.Type == sflow.EndActivity {
		pi.Status = sflow.ProcessInstaneStatusFinish
		t := time.Now()
		rs := manager.db.Model(&ProcessInstance{}).Where("id = ?", pi.Id).Updates(ProcessInstance{
			Status:     sflow.ProcessInstaneStatusFinish,
			FinishTime: &t,
		})
		return rs.Error
	} else {
		t := time.Now()
		rs := manager.db.Model(&ActivityInstance{}).Where("id = ?", ai.Id).Updates(ActivityInstance{
			Status:     sflow.ProcessInstaneStatusFinish,
			FinishTime: &t,
		})
		if rs.Error != nil {
			return rs.Error
		}
	}

	tos := pi.findNext(ai)
	if len(tos) == 0 && pi.allActivityFinish() {
		pi.Status = sflow.ProcessInstaneStatusFinish
		t := time.Now()
		rs := manager.db.Model(&ProcessInstance{}).Where("id = ?", pi.Id).Updates(ProcessInstance{
			Status:     sflow.ProcessInstaneStatusFinish,
			FinishTime: &t,
		})
		return rs.Error
	}
	var allNext []*ActivityInstance
	for _, to := range tos {
		adef, has := manager.pdm.GetActivity(pi.Definition, to)
		if has {
			next, err := pi.createActivityInstance(&adef)
			if err != nil {
				log.Println(err.Error())
			} else {
				allNext = append(allNext, next)
			}
		}
	}
	for _, next := range allNext {
		go next.Start()
	}
	return nil
}
func (pi *ProcessInstance) allActivityFinish() bool {
	acts, err := manager.ListActivities(pi.Id)
	if err != nil {
		return false
	}
	for _, a := range acts {
		if a.Status != sflow.ProcessInstaneStatusFinish {
			return false
		}
	}
	return true
}
func (pi *ProcessInstance) findNext(ai *ActivityInstance) (ads []int) {
	ads = []int{}
	var trans []int
	tds, _ := manager.pdm.ListTransitions(pi.Definition)
	for _, td := range tds {
		if td.From == ai.Definition {
			trans = append(trans, td.Id)
		}
	}
	var next []int
	for _, t := range trans {
		ti, err := manager.GetTransitions(pi.Id, t)
		if err != nil {
			continue
		}
		td, _ := manager.pdm.GetTransition(pi.Definition, ti.Definition)
		if !ti.Done {
			if ti.Check() {
				ti.Done = true
				manager.db.Model(&TransitionInstance{}).Where("id = ? and process_id = ?", ti.Id, ti.ProcessId).Update("done", true)
				next = append(next, td.To)
			}
		} else {
			next = append(next, td.To)
		}
	}
	if len(next) == 0 {
		return
	}
	trs, _ := manager.ListTransitions(pi.Id)
	for _, nid := range next {
		// all pre step done
		allFromOk := true
		for _, ti := range trs {
			td, _ := manager.pdm.GetTransition(pi.Definition, ti.Definition)
			if td.To == nid && !ti.Done {
				allFromOk = false
				break
			}
		}
		if allFromOk {
			ads = append(ads, nid)
		}
	}
	return
}
func (pi *ProcessInstance) suspend() {
	pi.Status = sflow.ProcessInstanceStatusSuspended
	//TODO deal suspend event
}
