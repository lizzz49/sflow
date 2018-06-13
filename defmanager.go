package sflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcessDefinitionManager struct {
	definitions  []*ProcessDefinition
	savePath     string
	maxProcessId int
}

func NewProcessDefinitionManager(path string) *ProcessDefinitionManager {
	pdm := &ProcessDefinitionManager{definitions: []*ProcessDefinition{}, savePath: path}
	fi, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0700)
		if err != nil {
			panic("can not make pds directory: " + err.Error())
		}
		fi, err = os.Stat(path)
	}
	if fi != nil && !fi.IsDir() {
		panic("can not make pds directory, because a pds file exist.")
	}
	if fi == nil {
		err = os.Mkdir(path, 0700)
		if err != nil {
			panic("can not make pds directory: " + err.Error())
		}
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".pdl") {
			b, e := ioutil.ReadFile(path)
			if e != nil {
				return e
			}
			var pd ProcessDefinition
			e = json.Unmarshal(b, &pd)
			if e != nil {
				return e
			} else {
				pd.parseMaxId()
				id, _ := strconv.Atoi(pd.Id)
				if id > pdm.maxProcessId {
					pdm.maxProcessId = id
				}
				pdm.definitions = append(pdm.definitions, &pd)
			}
		}
		return nil
	})
	return pdm
}

func (pdm *ProcessDefinitionManager) AddProcessDefinition(name string) *ProcessDefinition {
	pdm.maxProcessId++
	pd := newProcessDefinition(fmt.Sprintf("%d", pdm.maxProcessId), name)
	pdm.definitions = append(pdm.definitions, pd)
	return pd
}

func (pdm *ProcessDefinitionManager) List() []*ProcessDefinition {
	return pdm.definitions
}

func (pdm *ProcessDefinitionManager) Save() (errs []error) {
	for _, pd := range pdm.definitions {
		err := pd.Save(pdm.savePath)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}

func (pdm *ProcessDefinitionManager) GetProcessDefinitionById(id string) (*ProcessDefinition, bool) {
	for _, pd := range pdm.definitions {
		if pd.Id == id {
			return pd, true
		}
	}
	return nil, false
}
