package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
)

func init() {
	sflow.RegistryAction("action2", action2)
}
func action2(context *sflow.ProcessContext) bool {
	fmt.Println("2. hello word!")
	return true
}
