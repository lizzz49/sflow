package sflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

const (
	NotAnyActivityNode = "not any activity node."
	NotAnyTransmit     = " not any transmit."
	ActivityHasSameId  = "activity has same id: [%s]"
	ActivityLoop       = "activity lop: [%s]"
	ActivityNotFrom    = "activity [%s] not from."
	ActivityNotTo      = "activity [%s] not to."
	ActivityNotAction  = "activity [%s] not action definition."
	ActionNotFound     = "action [%s] in activity [%s] not found."
)
const (
	PDNewStatus = iota
	PDEditingStatus
	PDDiscardStatus
	PDPublishStatus
)

type ProcessDefinition struct {
	Definition
	StartActivity   ActivityDefinition      `json:"start_activity"`
	EndActivity     ActivityDefinition      `json:"end_activity"`
	Activities      []*ActivityDefinition   `json:"activities"`
	Transitions     []*TransitionDefinition `json:"transitions"`
	Status          int                     `json:"status"`
	maxActivityId   int
	maxTransitionId int
}

func newProcessDefinition(id, name string) *ProcessDefinition {
	pd := &ProcessDefinition{}
	pd.Id = id
	pd.Name = name
	pd.StartActivity = ActivityDefinition{
		Definition: Definition{"start", "start"},
		IsStart:    true,
	}
	pd.EndActivity = ActivityDefinition{
		Definition: Definition{"end", "end"},
		IsEnd:      true,
	}
	return pd
}

func (p *ProcessDefinition) AddActivityDefinition(name string, autoCommit bool) *ActivityDefinition {
	p.maxActivityId++
	a := newActivityDefinition(fmt.Sprintf("%d", p.maxActivityId), name, autoCommit)
	p.Activities = append(p.Activities, a)
	return a
}
func (p *ProcessDefinition) AddTransmitDefinition(name, from, to string, exp Express) *TransitionDefinition {
	p.maxTransitionId++
	t := newTransitionDefinition(fmt.Sprintf("%d", p.maxTransitionId), name, from, to, exp)
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
	aids := make(map[string]bool)
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
		aid, _ := strconv.Atoi(a.Id)
		if aid > p.maxActivityId {
			p.maxActivityId = aid
		}
		for _, ac := range a.Actions {
			acid, _ := strconv.Atoi(ac.Id)
			if acid > a.actionCount {
				a.actionCount = acid
			}
		}
	}
	for _, t := range p.Transitions {
		aid, _ := strconv.Atoi(t.Id)
		if aid > p.maxActivityId {
			p.maxActivityId = aid
		}
	}
}

func (p *ProcessDefinition) Save(path string) error {
	bs, _ := json.MarshalIndent(p, "", "    ")
	err := ioutil.WriteFile(fmt.Sprintf("%s/%s.pdl", path, p.Id), bs, 0600)
	if err == nil {
		log.Println("save prcess [" + p.Id + "] success.")
	} else {
		log.Println("save prcess [" + p.Id + "] fail:" + err.Error())
	}
	return err
}

func (p *ProcessDefinition) Publish() ([]DefError, bool) {
	errs, ok := p.Check()
	if ok {
		p.Status = PDPublishStatus
	}
	return errs, ok
}
