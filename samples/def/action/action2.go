package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/instance"
)

func init() {
	instance.RegistryAction("action2", action2)
}
func action2(context *sflow.ProcessContext, action *instance.ActionInstance, form []sflow.Value) bool {
	fmt.Println("2. hello word!")
	return true
}
