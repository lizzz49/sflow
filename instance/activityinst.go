package instance

import (
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/definition"
	"time"
)

type ActivityInstance struct {
	Id                int        //`gorm:"primaryKey;autoIncrement;column:id;type:int(11);comment:实例Id"`
	ProcessId         int        `gorm:"column:process_id;type:int(11);comment:流程实例Id"`
	Name              string     `gorm:"column:name;type:varchar(255);comment:实例名称"`
	ProcessDefinition int        `gorm:"column:process_definition;type:int(11);comment:流程定义Id"`
	Definition        int        `gorm:"column:definition;type:int(11);comment:定义Id"`
	Type              int        `gorm:"column:type;type:int(3);comment:活动类型"`
	Status            int        `gorm:"column:Status;type:int(3);comment:活动状态"`
	AutoCommit        bool       `gorm:"column:auto_commit;comment:是否自动完成"`
	CreateTime        time.Time  `gorm:"column:create_time;comment:活动创建时间"`
	StartTime         *time.Time `gorm:"column:start_time;comment:活动启动时间"`
	FinishTime        *time.Time `gorm:"column:finish_time;comment:活动完成时间"`
	ModelTime
	ControlBy
}

func (ActivityInstance) TableName() string {
	return "sflow_activity"
}
func (ai *ActivityInstance) Start() error {
	ai.Status = sflow.ProcessInstanceStatusStarted
	if ai.Type == sflow.StartActivity || ai.Type == sflow.EndActivity {
		return manager.FinishActivity(ai.ProcessId, ai.Id)
	}
	xs, _ := manager.pdm.ListActions(ai.ProcessDefinition, ai.Definition)
	for _, xd := range xs {
		action, err := ai.CreateActionInstance(xd)
		if err != nil {

		} else {
			action.Start()
		}
	}
	t := time.Now()
	rs := manager.db.Model(&ActivityInstance{}).
		Where("id = ? and process_id = ?", ai.Id, ai.ProcessId).
		Updates(ActivityInstance{
			Status:    sflow.ProcessInstanceStatusStarted,
			StartTime: &t,
		})
	return rs.Error
}

func (ai *ActivityInstance) CreateActionInstance(ad *definition.ActionDefinition) (x *ActionInstance, err error) {
	t := time.Now()
	a := &ActionInstance{
		ProcessId:  ai.ProcessId,
		ActivityId: ai.Id,
		Name:       ad.Name,
		InvokeName: ad.InvokerName,
		Definition: ad.Id,
		AutoCommit: ad.AutoCommit,
		Status:     sflow.ProcessInstanceStatusNew,
		CreateTime: t,
		StartTime:  &t,
	}
	rs := manager.db.Model(&ActionInstance{}).Create(a)
	return a, rs.Error
}
