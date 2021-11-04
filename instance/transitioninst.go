package instance

import (
	"github.com/lizzz49/sflow"
	"strconv"
)

type TransitionInstance struct {
	Id                int    //`gorm:"primaryKey;autoIncrement;column:id;type:int(11);comment:跳转实例Id"`
	Name              string `gorm:"column:name;type:varchar(255);comment:跳转实例名称"`
	ProcessId         int    `gorm:"column:process_id;type:int(11);comment:流程实例Id"`
	ProcessDefinition int    `gorm:"column:process_definition;type:int(11);comment:流程定义Id"`
	Definition        int    `gorm:"column:definition;type:int(11);comment:定义Id"`
	Done              bool   `gorm:"column:done;comment:是否满足条件"`
}

func (t *TransitionInstance) Check() bool {
	def, err := manager.pdm.GetTransition(t.ProcessDefinition, t.Definition)
	if err != nil {
		return false
	}
	if def.AlwaysTrue {
		return true
	} else {
		exp := def.Express
		key := exp.Key
		data, _ := manager.GetProcessData(t.ProcessId, key)
		value := exp.Value
		if data.Type != value.Type {
			return false
		}
		op := exp.OP
		switch op {
		case ">":
			switch data.Type {
			case sflow.Int64Type:
				dv, _ := strconv.ParseInt(data.Value, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case sflow.Float64Type:
				dv, _ := strconv.ParseFloat(data.Value, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Value > value.Data
			}
		case "=":
			switch data.Type {
			case sflow.Int64Type:
				dv, _ := strconv.ParseInt(data.Value, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case sflow.Float64Type:
				dv, _ := strconv.ParseFloat(data.Value, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Value > value.Data
			}
		case ">=":
			switch data.Type {
			case sflow.Int64Type:
				dv, _ := strconv.ParseInt(data.Value, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case sflow.Float64Type:
				dv, _ := strconv.ParseFloat(data.Value, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Value > value.Data
			}
		case "<":
			switch data.Type {
			case sflow.Int64Type:
				dv, _ := strconv.ParseInt(data.Value, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case sflow.Float64Type:
				dv, _ := strconv.ParseFloat(data.Value, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Value > value.Data
			}
		case "<=":
			switch data.Type {
			case sflow.Int64Type:
				dv, _ := strconv.ParseInt(data.Value, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case sflow.Float64Type:
				dv, _ := strconv.ParseFloat(data.Value, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Value > value.Data
			}
		case "IsNull":
			switch data.Type {
			case sflow.Int64Type:
				dv, _ := strconv.ParseInt(data.Value, 10, 64)
				return dv == 0
			case sflow.Float64Type:
				dv, _ := strconv.ParseFloat(data.Value, 64)
				return dv == 0
			default:
				return data.Value == ""
			}
		}
	}
	return false
}
