package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/instance"
)

func init() {
	instance.RegistryAction("action3", action3)
}
func action3(context *sflow.ProcessContext, action *instance.ActionInstance, form []sflow.Value) bool {
	fmt.Println("3. hello word!")
	return true
}
