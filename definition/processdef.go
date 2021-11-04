package definition

import (
	"encoding/json"
	"fmt"
	"github.com/lizzz49/sflow"
	"gorm.io/gorm"
)

const (
	NotAnyActivityNode = "not any activity node."
	NotAnyTransmit     = " not any transmit."
	ActivityHasSameId  = "activity has same id: [%d]"
	ActivityLoop       = "activity lop: [%d]"
	ActivityNotFrom    = "activity [%d] not from."
	ActivityNotTo      = "activity [%d] not to."
	ActivityNotAction  = "activity [%d] not action definition."
	ActionNotFound     = "action [%s] in activity [%d] not found."
)
const (
	PDNewStatus = iota
	PDEditingStatus
	PDDiscardStatus
	PDPublishStatus
)
const (
	ActivityStartId = iota
	ActivityEndId
)

type ProcessDefinition struct {
	Definition
	StartActivity   ActivityDefinition      `json:"startActivity,omitempty" gorm:"-"`
	EndActivity     ActivityDefinition      `json:"endActivity,omitempty" gorm:"-"`
	Activities      []*ActivityDefinition   `json:"activities,omitempty" gorm:"-"`
	Transitions     []*TransitionDefinition `json:"transitions,omitempty" gorm:"-"`
	Status          int                     `json:"-" gorm:"column:status;type:int(3);comment:流程状态"`
	maxActivityId   int                     `gorm:"-"`
	maxTransitionId int                     `gorm:"-"`
	JSON            string                  `json:"-" gorm:"column:json;type:text;comment:JSON格式流程定义"`
	Version         string                  `json:"-" gorm:"column:version;type:varchar(32);comment:版本号"`
}

func (ProcessDefinition) TableName() string {
	return "sflow_process_def"
}
func NewProcessDefinition(id int, name string) *ProcessDefinition {
	pd := &ProcessDefinition{Status: PDNewStatus}
	pd.Id = id
	pd.Name = name
	pd.StartActivity = ActivityDefinition{
		Definition: Definition{ActivityStartId, "start"},
		Type:       sflow.StartActivity,
	}
	pd.EndActivity = ActivityDefinition{
		Definition: Definition{ActivityEndId, "end"},
		Type:       sflow.EndActivity,
	}
	pd.maxActivityId = ActivityEndId
	return pd
}

func (p *ProcessDefinition) AddActivityDefinition(name string, autoCommit bool) *ActivityDefinition {
	p.maxActivityId++
	a := newActivityDefinition(p.maxActivityId, name, autoCommit)
	p.Activities = append(p.Activities, a)
	return a
}
func (p *ProcessDefinition) AddTransmitDefinition(name string, from, to int, exp Express) *TransitionDefinition {
	p.maxTransitionId++
	t := newTransitionDefinition(p.maxTransitionId, name, from, to, exp)
	p.Transitions = append(p.Transitions, t)
	return t
}

type DefError struct {
	ErrMsg string
}

func (de DefError) Error() string {
	return de.ErrMsg
}
func (p *ProcessDefinition) Check() ([]DefError, bool) {
	var des []DefError
	if p.Activities == nil || len(p.Activities) == 0 {
		des = append(des, DefError{NotAnyActivityNode})
	}
	if p.Transitions == nil || len(p.Transitions) == 0 {
		des = append(des, DefError{NotAnyTransmit})
	}
	aids := make(map[int]bool)
	if p.Activities != nil {
		for _, node := range p.Activities {
			if aids[node.Id] {
				des = append(des, DefError{fmt.Sprintf(ActivityHasSameId, node.Id)})
			} else {
				aids[node.Id] = true
			}
			if node.Actions == nil || len(node.Actions) == 0 {
				des = append(des, DefError{fmt.Sprintf(ActivityNotAction, node.Id)})
			}
			if p.Transitions != nil && len(p.Transitions) > 0 {
				var hasFrom, hasTo bool
				for _, transmit := range p.Transitions {
					if transmit.From == transmit.To {
						des = append(des, DefError{fmt.Sprintf(ActivityLoop, node.Id)})
					} else {
						hasTo = hasTo || transmit.From == node.Id
						hasFrom = hasFrom || transmit.To == node.Id
					}
					if hasTo && hasFrom {
						break
					}
				}
				if !hasFrom {
					des = append(des, DefError{fmt.Sprintf(ActivityNotFrom, node.Id)})
				}
				if !hasTo {
					des = append(des, DefError{fmt.Sprintf(ActivityNotTo, node.Id)})
				}
			}
		}
	}

	if p.Transitions != nil {
		for i := 0; i < len(p.Transitions); i++ {
			for j := i + 1; j < len(p.Transitions); j++ {

			}
		}
	}
	return des, des == nil || len(des) == 0
}

func (p *ProcessDefinition) parseMaxId() {
	for _, a := range p.Activities {
		aid := a.Id
		if aid > p.maxActivityId {
			p.maxActivityId = aid
		}
		for _, ac := range a.Actions {
			acid := ac.Id
			if acid > a.actionCount {
				a.actionCount = acid
			}
		}
	}
	for _, t := range p.Transitions {
		aid := t.Id
		if aid > p.maxActivityId {
			p.maxActivityId = aid
		}
	}
}

func (p *ProcessDefinition) Save(db *gorm.DB) error {
	bs, _ := json.Marshal(p)
	p.JSON = string(bs)
	p.Version = "0"
	rs := db.Model(&ProcessDefinition{}).Create(p)
	return rs.Error
}

func (p *ProcessDefinition) Publish() ([]DefError, bool) {
	errs, ok := p.Check()
	if ok {
		p.Status = PDPublishStatus
	}
	return errs, ok
}
