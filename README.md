# sflow
A simple and flexible workflow framework, refer to WfMC XPDL model, but process definition with json.

### 1. Framework feature
* Simple definition (process,activity,action,transition)
* Support process,activity,action(work payload) auto commit (use for robot) 
* Condition transition support
* Process participant (user or role) support (TODO)

### 2. Installation
```bash
go get github.com/lizzz49/sflow
```

### 3. Process Definition
There are two ways to complete the process definition.
#### 1.Use process definition SDK
```
//sample code in samples/def folder
//new a process definition manager with save path
pdm := sflow.NewProcessDefinitionManager("../inst/pds")
//add a new process definition with process name
process := pdm.AddProcessDefinition("express")

//add a activity definition to the process
node1 := process.AddActivityDefinition("node1",true)
//add a action to the activity definition
node1.AddActionDefinition("print something","action1",true)

node2 := process.AddActivityDefinition("node2",true)
node2.AddActionDefinition("print something","action2",true)

//add a transition
process.AddTransmitDefinition("sTn1",process.StartActivity.Id,node1.Id,pdl.Express{})
//add a transition with condition
exp:= sflow.Express{Key:"age",OP:">",Value:pdl.Value{Type:pdl.Int64Type,Data:"18"}}
process.AddTransmitDefinition("n1tn2",node1.Id,node2.Id,exp)
process.AddTransmitDefinition("n2Te",node2.Id,process.EndActivity.Id,pdl.Express{})

//publish the process after process definition
process.Publish()
```
#### 2.Edit json directly
```json
{
    "id": "3",
    "name": "express",
    "start_activity": {
        "id": "start",
        "name": "start",
        "is_start": true
    },
    "end_activity": {
        "id": "end",
        "name": "end",
        "is_end": true
    },
    "activities": [
        {
            "id": "1",
            "name": "node1",
            "auto_commit": true,
            "actions": [
                {
                    "id": "1",
                    "name": "print something",
                    "auto_commit": true,
                    "invoker_name": "action1",
                    "config": []
                }
            ]
        },
        {
            "id": "2",
            "name": "node2",
            "auto_commit": true,
            "actions": [
                {
                    "id": "1",
                    "name": "print something",
                    "auto_commit": true,
                    "invoker_name": "action2"
                }
            ]       
        }
    ],
    "transitions": [
        {
            "id": "1",
            "name": "sTn1",
            "from": "start",
            "to": "1",
            "always_true": true
        },
        {
            "id": "2",
            "name": "n1tn2",
            "from": "1",
            "to": "2",
            "always_true": false,
            "express": {
                "key": "age",
                "op": "\u003e",
                "value": {
                    "type": 0,
                    "data": "18"
                }
            }
        },
        {
            "id": "3",
            "name": "n2Te",
            "from": "2",
            "to": "end",
            "always_true": true
        }
    ],
    "status": 12
}
```
### 4. Registry work payload
```go
//package github.com/lizzz49/sflow/samples/def/action
package action

import (
	"github.com/lizzz49/sflow"
	"fmt"
)

func init(){
	sflow.RegistryAction("action1",action1)
}
func action1(context *sflow.ProcessContext)bool{
	fmt.Println("1. hello word!")
	return true
}

```
### 5. Run process
#### Load work payload (action)
```
_ "github.com/lizzz49/sflow/samples/def/action"
```
#### Run
```
pdm := sflow.NewProcessDefinitionManager("pds")
pim := sflow.NewProcessInstanceManager()
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
//wait for  process finish
log.Println("express process exit with code:",pi.Waite())
```