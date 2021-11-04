package instance

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/definition"
	"gorm.io/gorm"
	"log"
	"time"
)

type ProcessInstanceManager struct {
	pdm *definition.ProcessDefinitionManager
	db  *gorm.DB
}

var manager *ProcessInstanceManager

func NewProcessInstanceManager(db *gorm.DB) *ProcessInstanceManager {
	if manager == nil {
		manager = &ProcessInstanceManager{pdm: definition.NewProcessDefinitionManager(db), db: db}
	}
	return manager
}

func (pim *ProcessInstanceManager) CreateProcessInstance(name string, def *definition.ProcessDefinition) (pi *ProcessInstance, err error) {
	des, ok := def.Check()
	if ok {
		des, ok = checkProcessDefinition(def)
	}
	if ok {
		pi = &ProcessInstance{
			Name:       name,
			Definition: def.Id,
			Status:     sflow.ProcessInstanceStatusNew,
			CreateTime: time.Now(),
			StartTime:  nil,
			FinishTime: nil,
			context:    nil,
		}
		rs := manager.db.Model(&ProcessInstance{}).Create(&pi)
		if rs.Error != nil {
			return pi, rs.Error
		}
		tds, _ := pim.pdm.ListTransitions(pi.Definition)
		for _, transmit := range tds {
			_, er := pi.createTransition(&transmit)
			if er != nil {
				log.Println(er)
			}
		}
		return pi, nil
	} else {
		var errs string
		for i, er := range des {
			errs += fmt.Sprintf("%d\t%s\n", i, er.Error())
		}
		return nil, fmt.Errorf(errs)
	}
}

func (pim *ProcessInstanceManager) TerminateProcessInstance(id int) bool {
	pi, err := pim.GetProcess(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return true
		}
		return false
	}
	if pi.Status == sflow.ProcessInstanceStatusTerminated {
		return true
	}
	if pi.Status == sflow.ProcessInstanceStatusNew {
		pi.Status = sflow.ProcessInstanceStatusTerminated
	} else if pi.Status == sflow.ProcessInstanceStatusStarted || pi.Status == sflow.ProcessInstanceStatusSuspended {
		pi.Status = sflow.ProcessInstanceStatusTerminated
	}
	rs := pim.db.Model(&ProcessInstance{}).Update("status", pi.Status)
	return rs.Error == nil
}

func (pim *ProcessInstanceManager) ListProcess(did int, status int, name string, page, size int) (total int64, ps []ProcessInstance, err error) {
	if size == 0 {
		size = 10
	}
	if page < 1 {
		size = 1
	}
	rs := pim.db.Model(&ProcessInstance{})
	if did != 0 {
		rs = rs.Where("definition = ?", did)
	}
	if status != 0 {
		rs = rs.Where("status = ?", status)
	}
	if name != "" {
		rs = rs.Where("name like ?", fmt.Sprintf("%s%s%s", "%", name, "%"))
	}
	rs = rs.Limit(size).Offset((page - 1) * size).Find(&ps)
	if rs.Error != nil {
		return total, ps, rs.Error
	}
	rs = rs.Limit(-1).Offset(-1).Count(&total)
	return total, ps, rs.Error
}

func (pim *ProcessInstanceManager) GetProcess(pid int) (p ProcessInstance, err error) {
	rs := pim.db.Model(&ProcessInstance{}).First(&p, pid)
	return p, rs.Error
}
func (pim *ProcessInstanceManager) StartProcess(pid int) error {
	p, err := pim.GetProcess(pid)
	if err != nil {
		return err
	}
	return p.Start()
}
func (pim *ProcessInstanceManager) ListActivities(pid int) (acts []ActivityInstance, err error) {
	rs := pim.db.Model(&ActivityInstance{}).Where("process_id = ?", pid).Find(&acts)
	return acts, rs.Error
}
func (pim *ProcessInstanceManager) GetActivity(pid, aid int) (act ActivityInstance, err error) {
	rs := pim.db.Model(&ActivityInstance{}).Where("id = ? and process_id = ?", aid, pid).Take(&act)
	return act, rs.Error
}
func (pim *ProcessInstanceManager) FinishActivity(pid, aid int) error {
	p, err := pim.GetProcess(pid)
	if err != nil {
		return err
	}
	act, err := pim.GetActivity(pid, aid)
	if err != nil {
		return err
	}
	return p.FinishActivity(&act)
}
func (pim *ProcessInstanceManager) ListActions(pid, aid int) (xs []ActionInstance, err error) {
	rs := pim.db.Model(&ActionInstance{}).Where("process_id = ? and activity_id = ?", pid, aid).Find(&xs)
	return xs, rs.Error
}
func (pim *ProcessInstanceManager) GetAction(pid, aid, xid int) (x ActionInstance, err error) {
	rs := pim.db.Model(&ActionInstance{}).Where("process_id=? and activity_id=? and id=?", pid, aid, xid).Take(&x)
	return x, rs.Error
}
func (pim *ProcessInstanceManager) InvokeAction(pid, aid, xid int, form []sflow.Value) error {
	if len(form) == 0 {
		return nil
	}
	tx := pim.db.Begin()
	rs := tx.Model(&ActionData{}).Where("process_id = ? and activity_id = ? and action_id = ?", pid, aid, xid).Delete(ActionData{})
	if rs.Error != nil {
		tx.Rollback()
		return rs.Error
	}
	var datas []ActionData
	for _, data := range form {
		datas = append(datas, ActionData{
			ProcessId:  pid,
			ActivityId: aid,
			ActionId:   xid,
			Key:        data.Key,
			Type:       data.Type,
			Value:      data.Data,
		})
	}
	rs = tx.Model(&ActionData{}).Create(&datas)
	if rs.Error != nil {
		tx.Rollback()
		return rs.Error
	}
	return tx.Commit().Error
}
func (pim *ProcessInstanceManager) ListTransitions(pid int) (trs []TransitionInstance, err error) {
	rs := pim.db.Model(&TransitionInstance{}).Where("process_id = ?", pid).Find(&trs)
	return trs, rs.Error
}
func (pim *ProcessInstanceManager) GetTransitions(pid, td int) (tr TransitionInstance, err error) {
	rs := pim.db.Model(&TransitionInstance{}).Where("process_id = ? and definition = ?", pid, td).Take(&tr)
	return tr, rs.Error
}
func (pim *ProcessInstanceManager) GetProcessData(pid int, key string) (data ActionData, err error) {
	rs := pim.db.Model(&ActionData{}).Where("process_id = ? and key = ?", pid, key).Take(&data)
	return data, rs.Error
}
func (pim *ProcessInstanceManager) GetProcessContext(pid int) (ctx sflow.ProcessContext) {
	var datas []ActionData
	pim.db.Model(&ActionData{}).Where("process_id = ?", pid).Find(&datas)
	for _, data := range datas {
		ctx[data.Key] = sflow.Value{
			Key:  data.Key,
			Type: data.Type,
			Data: data.Value,
		}
	}
	return
}

func (pim *ProcessInstanceManager) FinishAction(x *ActionInstance) error {
	x.Status = sflow.ProcessInstaneStatusFinish
	t := time.Now()
	rs := manager.db.Model(&ActionInstance{}).Where("id = ? and process_id=? and activity_id = ?", x.Id, x.ProcessId, x.ActivityId).Updates(ActionInstance{
		Status:     sflow.ProcessInstaneStatusFinish,
		FinishTime: &t,
	})
	if rs.Error != nil {
		return rs.Error
	}
	ai, err := pim.GetActivity(x.ProcessId, x.ActivityId)
	if rs.Error != nil {
		return rs.Error
	}
	if !ai.AutoCommit {
		return nil
	}
	allFinish := true
	xs, err := pim.ListActions(x.ProcessId, x.ActivityId)
	if err != nil {
		return err
	}
	for _, action := range xs {
		if action.Status != sflow.ProcessInstaneStatusFinish {
			allFinish = false
			break
		}
	}
	if allFinish {
		return pim.FinishActivity(x.ProcessId, x.ActivityId)
	}
	return nil
}
