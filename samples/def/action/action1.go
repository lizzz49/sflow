package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
)

func init() {
	sflow.RegistryAction("action1", action1)
}
func action1(context *sflow.ProcessContext) bool {
	fmt.Println("1. hello word!")
	return true
}
