package definition

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type ProcessDefinitionManager struct {
	definitions  []*ProcessDefinition
	db           *gorm.DB
	maxProcessId int
}

var manager *ProcessDefinitionManager

func NewProcessDefinitionManager(db *gorm.DB) *ProcessDefinitionManager {
	if manager == nil {
		pdm := &ProcessDefinitionManager{definitions: []*ProcessDefinition{}, db: db}
		var pds []ProcessDefinition
		rs := db.Model(&ProcessDefinition{}).Find(&pds)
		if rs.Error != nil {
			panic(rs.Error.Error())
		}
		for _, pd := range pds {
			var p ProcessDefinition
			_ = json.Unmarshal([]byte(pd.JSON), &p)
			p.Id = pd.Id
			p.Name = pd.Name
			p.parseMaxId()
			pdm.definitions = append(pdm.definitions, &p)
			if pdm.maxProcessId < pd.Id {
				pdm.maxProcessId = pd.Id
			}
		}
		manager = pdm
	}
	return manager
}

func (pdm *ProcessDefinitionManager) AddProcessDefinition(name string) *ProcessDefinition {
	pdm.maxProcessId++
	pd := NewProcessDefinition(pdm.maxProcessId, name)
	pdm.definitions = append(pdm.definitions, pd)
	return pd
}

func (pdm *ProcessDefinitionManager) GetActivity(pid, aid int) (act ActivityDefinition, has bool) {
	pd, has := pdm.GetProcessDefinitionById(pid)
	if !has {
		return act, has
	}
	if aid == 0 {
		return pd.StartActivity, has
	}
	if aid == 1 {
		return pd.EndActivity, has
	}
	has = false
	for _, a := range pd.Activities {
		if a.Id == aid {
			act = *a
			has = true
			break
		}
	}
	return
}

func (pdm *ProcessDefinitionManager) ListActions(pid, aid int) (xs []*ActionDefinition, has bool) {
	ad, has := pdm.GetActivity(pid, aid)
	if !has {
		return
	}
	return ad.Actions, true
}

func (pdm *ProcessDefinitionManager) ListTransitions(pid int) (trs []TransitionDefinition, err error) {
	pd, has := pdm.GetProcessDefinitionById(pid)
	if !has {
		return trs, fmt.Errorf("process definition %d not found", pid)
	}
	for _, tr := range pd.Transitions {
		trs = append(trs, *tr)
	}
	return
}

func (pdm *ProcessDefinitionManager) GetTransition(pid, tid int) (td TransitionDefinition, err error) {
	pd, has := pdm.GetProcessDefinitionById(pid)
	if !has {
		return td, fmt.Errorf("process definition %d not found", pid)
	}
	for _, tr := range pd.Transitions {
		if tr.Id == tid {
			return *tr, nil
		}
	}
	return td, fmt.Errorf("process defnition %d transition %d not found", pid, tid)
}

func (pdm *ProcessDefinitionManager) List() []*ProcessDefinition {
	return pdm.definitions
}

func (pdm *ProcessDefinitionManager) Save() (errs []error) {
	for _, pd := range pdm.definitions {
		err := pd.Save(pdm.db)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func (pdm *ProcessDefinitionManager) GetProcessDefinitionById(id int) (*ProcessDefinition, bool) {
	for _, pd := range pdm.definitions {
		if pd.Id == id {
			return pd, true
		}
	}
	return nil, false
}
