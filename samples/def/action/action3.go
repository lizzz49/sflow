package action

import (
	"fmt"
	"github.com/lizzz49/sflow"
)

func init() {
	sflow.RegistryAction("action3", action3)
}
func action3(context *sflow.ProcessContext) bool {
	fmt.Println("3. hello word!")
	return true
}
