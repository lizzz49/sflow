package main

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/definition"
	"github.com/lizzz49/sflow/instance"
	_ "github.com/lizzz49/sflow/samples/def/action"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

func main() {
	var db *gorm.DB
	dialet := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local",
		"root", "", "localhost", 3306, "sflow")
	db, err := gorm.Open(mysql.Open(dialet), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		AllowGlobalUpdate: false,
	})
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(
		&instance.ProcessInstance{},
		&instance.ActivityInstance{},
		&instance.ActionInstance{},
		&instance.ActionData{},
		&instance.TransitionInstance{})
	pdm := definition.NewProcessDefinitionManager(db)
	pds := pdm.List()
	if len(pds) == 0 {
		println("not process definition found.")
		return
	}
	pim := instance.NewProcessInstanceManager(db)
	runLineProcess(pdm, pim)
	//runSplitProcess(pdm, pim)
	//runExpressProcess(pdm, pim)
	ch := make(chan int, 1)
	<-ch
}

func runLineProcess(pdm *definition.ProcessDefinitionManager, pim *instance.ProcessInstanceManager) {
	log.Println("run line process.")
	process, has := pdm.GetProcessDefinitionById(1)
	if !has {
		println("not process definition found with id: 1")
		return
	}
	process.Status = definition.PDPublishStatus
	pi, err := pim.CreateProcessInstance(process.Name, process)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(pi.Start())
}

func runSplitProcess(pdm *definition.ProcessDefinitionManager, pim *instance.ProcessInstanceManager) {
	log.Println("run split process.")
	process, has := pdm.GetProcessDefinitionById(2)
	if !has {
		println("not process definition found with id: 2")
		return
	}
	process.Status = definition.PDPublishStatus
	pi, err := pim.CreateProcessInstance(process.Name, process)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(pi.Start())
}

func runExpressProcess(pdm *definition.ProcessDefinitionManager, pim *instance.ProcessInstanceManager) {
	log.Println("run express process.")
	process, has := pdm.GetProcessDefinitionById(3)
	if !has {
		println("not process definition found with id: 3")
		return
	}
	process.Status = definition.PDPublishStatus
	pi, err := pim.CreateProcessInstance(process.Name, process)
	if err != nil {
		log.Println(err.Error())
		return
	}
	ctx := make(sflow.ProcessContext)
	ctx["age"] = sflow.Value{Type: sflow.Int64Type, Data: "30"}
	pi.Init(ctx)
	log.Println(pi.Start())
}
