package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/instance"
)

func init() {
	instance.RegistryAction("action1", action1)
}
func action1(context *sflow.ProcessContext, action *instance.ActionInstance, form []sflow.Value) bool {
	fmt.Println("1. hello word!")
	return true
}
