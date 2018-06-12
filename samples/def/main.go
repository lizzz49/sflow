package main

import (
	pdl "github.com/lizzz49/sflow"
)

func main() {
	pdm := pdl.NewProcessDefinitionManager("../inst/pds")
	createLineProcess(pdm)
	createSplitProcess(pdm)
	createExpressProcess(pdm)
	pdm.Save()

}
func createExpressProcess(pdm *pdl.ProcessDefinitionManager){
	process := pdm.AddProcessDefinition("express")

	node1 := process.AddActivityDefinition("node1",true)
	node1.AddActionDefinition("print something","action1",true)

	node2 := process.AddActivityDefinition("node2",true)
	node2.AddActionDefinition("print something","action2",true)

	process.AddTransmitDefinition("sTn1",process.StartActivity.Id,node1.Id,pdl.Express{})

	exp:= pdl.Express{Key:"age",OP:">",Value:pdl.Value{Type:pdl.Int64Type,Data:"18"}}
	process.AddTransmitDefinition("n1tn2",node1.Id,node2.Id,exp)
	process.AddTransmitDefinition("n2Te",node2.Id,process.EndActivity.Id,pdl.Express{})

	process.Publish()
}

func createLineProcess(pdm *pdl.ProcessDefinitionManager){
	process := pdm.AddProcessDefinition("line")

	node1 := process.AddActivityDefinition("node1",true)
	node1.AddActionDefinition("print something","action1",true)

	node2 := process.AddActivityDefinition("node2",true)
	node2.AddActionDefinition("print something","action2",true)

	process.AddTransmitDefinition("sTn1",process.StartActivity.Id,node1.Id,pdl.Express{})
	process.AddTransmitDefinition("n1tn2",node1.Id,node2.Id,pdl.Express{})
	process.AddTransmitDefinition("n2Te",node2.Id,process.EndActivity.Id,pdl.Express{})

	process.Publish()
}
func createSplitProcess(pdm *pdl.ProcessDefinitionManager){
	process := pdm.AddProcessDefinition("express")

	node1 := process.AddActivityDefinition("node1",true)
	node1.AddActionDefinition("print something","action1",true)

	node2 := process.AddActivityDefinition("node2",true)
	node2.AddActionDefinition("print something","action2",true)

	node3 := process.AddActivityDefinition("node3",true)
	node3.AddActionDefinition("print something","action3",true)

	node4 := process.AddActivityDefinition("node4",true)
	node4.AddActionDefinition("print something","action4",true)

	process.AddTransmitDefinition("sTn1",process.StartActivity.Id,node1.Id,pdl.Express{})
	process.AddTransmitDefinition("n1tn2",node1.Id,node2.Id,pdl.Express{})
	process.AddTransmitDefinition("n1tn3",node1.Id,node3.Id,pdl.Express{})
	process.AddTransmitDefinition("n2tn4",node2.Id,node4.Id,pdl.Express{})
	process.AddTransmitDefinition("n3tn4",node3.Id,node4.Id,pdl.Express{})
	process.AddTransmitDefinition("n4Te",node4.Id,process.EndActivity.Id,pdl.Express{})

	process.Publish()
}