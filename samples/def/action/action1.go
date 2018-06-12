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
