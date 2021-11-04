package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
	"github.com/lizzz49/sflow/instance"
)

func init() {
	instance.RegistryAction("action4", action4)
}
func action4(context *sflow.ProcessContext, action *instance.ActionInstance, form []sflow.Value) bool {
	fmt.Println("4. hello word!")
	return true
}
