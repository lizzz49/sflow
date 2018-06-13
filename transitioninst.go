package sflow

import (
	"strconv"
)

type TransitionInstance struct {
	Process    *ProcessInstance
	Definition *TransitionDefinition
	Done       bool
}

func (t *TransitionInstance) Check() bool {
	if t.Definition.AlwaysTrue {
		return true
	} else {
		context := t.Process.context
		exp := t.Definition.Express
		key := exp.Key
		data := context[key]
		value := exp.Value
		if data.Type != value.Type {
			return false
		}
		op := exp.OP
		switch op {
		case ">":
			switch data.Type {
			case Int64Type:
				dv, _ := strconv.ParseInt(data.Data, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case Float64Type:
				dv, _ := strconv.ParseFloat(data.Data, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Data > value.Data
			}
		case "=":
			switch data.Type {
			case Int64Type:
				dv, _ := strconv.ParseInt(data.Data, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case Float64Type:
				dv, _ := strconv.ParseFloat(data.Data, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Data > value.Data
			}
		case ">=":
			switch data.Type {
			case Int64Type:
				dv, _ := strconv.ParseInt(data.Data, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case Float64Type:
				dv, _ := strconv.ParseFloat(data.Data, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Data > value.Data
			}
		case "<":
			switch data.Type {
			case Int64Type:
				dv, _ := strconv.ParseInt(data.Data, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case Float64Type:
				dv, _ := strconv.ParseFloat(data.Data, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Data > value.Data
			}
		case "<=":
			switch data.Type {
			case Int64Type:
				dv, _ := strconv.ParseInt(data.Data, 10, 64)
				vv, _ := strconv.ParseInt(value.Data, 10, 64)
				return dv > vv
			case Float64Type:
				dv, _ := strconv.ParseFloat(data.Data, 64)
				vv, _ := strconv.ParseFloat(value.Data, 64)
				return dv > vv
			default:
				return data.Data > value.Data
			}
		case "IsNull":
			switch data.Type {
			case Int64Type:
				dv, _ := strconv.ParseInt(data.Data, 10, 64)
				return dv == 0
			case Float64Type:
				dv, _ := strconv.ParseFloat(data.Data, 64)
				return dv == 0
			default:
				return data.Data == ""
			}
		}
	}
	return false
}
