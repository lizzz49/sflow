package instance

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"time"
)

var (
	actions = make(map[string]ActionInvoker)
)

type ActionInstance struct {
	Id         int        //`gorm:"primaryKey;autoIncrement;column:id;type:int(11);comment:实例Id"`
	ProcessId  int        `gorm:"column:process_id;type:int(11);comment:流程Id"`
	ActivityId int        `gorm:"column:activity_id;type:int(11);comment:活动Id"`
	Name       string     `gorm:"column:name;type:varchar(255);comment:实例名称"`
	Definition int        `gorm:"column:definition;type:int(11);comment:定义Id"`
	AutoCommit bool       `gorm:"column:auto_commit;type:int(1);comment:自动提交"`
	Status     int        `gorm:"column:Status;type:int(3);comment:活动状态"`
	CreateTime time.Time  `gorm:"column:create_time;comment:活动创建时间"`
	StartTime  *time.Time `gorm:"column:start_time;comment:活动启动时间"`
	FinishTime *time.Time `gorm:"column:finish_time;comment:活动完成时间"`
	//activity         *ActivityInstance
	//actionDefinition *definition.ActionDefinition
	formData []sflow.Value
	ModelTime
	ControlBy
}

type ActionData struct {
	Id         int    //`gorm:"primaryKey;autoIncrement;column:id;type:int(11);comment:数据Id"`
	ProcessId  int    `gorm:"column:process_id;type:int(11);comment:流程Id"`
	ActivityId int    `gorm:"column:activity_id;type:int(11);comment:活动Id"`
	ActionId   int    `gorm:"column:action_id;type:int(11);comment:动作Id"`
	Key        string `gorm:"column:key;type:varchar(255);comment:数据名称"`
	Type       int    `gorm:"column:type;type:int(3);comment:数据类型"`
	Value      string `gorm:"column:value;type:varchar(255);comment:数据值"`
	ModelTime
	ControlBy
}

type ActionInvoker struct {
	Name       string
	InvokeFunc InvokeFunc
}

type InvokeFunc func(*sflow.ProcessContext, *ActionInstance, []sflow.Value) bool

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
	ai.Status = sflow.ProcessInstanceStatusStarted
	invoker, has := GetActionInvoker(ai.Name)
	if has {
		invokeFunc := invoker.InvokeFunc
		context := manager.GetProcessContext(ai.ProcessId)
		invokeFunc(&context, ai, []sflow.Value{})
		if ai.AutoCommit {
			return manager.FinishAction(ai)
		}
	} else {
		return manager.FinishAction(ai)
	}
	return nil
}
