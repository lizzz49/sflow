package main

import (
	"github.com/lizzz49/sflow"
	_"github.com/lizzz49/sflow/samples/def/action"
	"log"
)

func main() {

	pdm := sflow.NewProcessDefinitionManager("pds")
	pds := pdm.List()
	if len(pds) == 0 {
		println("not process definition found.")
		return
	}
	pim := sflow.NewProcessInstanceManager()
	runLineProcess(pdm,pim)
	runSplitProcess(pdm,pim)
	runExpressProcess(pdm,pim)
}

func runLineProcess(pdm *sflow.ProcessDefinitionManager,pim *sflow.ProcessInstanceManager){
	log.Println("run line process.")
	process,has := pdm.GetProcessDefinitionById("1")
	if !has{
		println("not process definition found with id: 1")
		return
	}

	pi := pim.CreateProcessInstance(process)
	pi.Start()
	log.Println("line process exit with code:",pi.Wait())
}

func runSplitProcess(pdm *sflow.ProcessDefinitionManager,pim *sflow.ProcessInstanceManager){
	log.Println("run split process.")
	process,has := pdm.GetProcessDefinitionById("2")
	if !has{
		println("not process definition found with id: 2")
		return
	}

	pi := pim.CreateProcessInstance(process)
	pi.Start()
	log.Println("split process exit with code:",pi.Wait())
}

func runExpressProcess(pdm *sflow.ProcessDefinitionManager,pim *sflow.ProcessInstanceManager){
	log.Println("run express process.")
	process,has := pdm.GetProcessDefinitionById("3")
	if !has{
		println("not process definition found with id: 3")
		return
	}

	pi := pim.CreateProcessInstance(process)
	ctx := make(sflow.ProcessContext)
	ctx["age"] = sflow.Value{Type:sflow.Int64Type,Data:"30"}
	pi.Init(ctx)
	pi.Start()
	log.Println("express process exit with code:",pi.Wait())
}