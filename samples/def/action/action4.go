package action

import (
	"github.com/lizzz49/sflow"
	"fmt"
)

func init(){
	sflow.RegistryAction("action4",action4)
}
func action4(context *sflow.ProcessContext)bool{
	fmt.Println("4. hello word!")
	return true
}
