package definition

import (
	"github.com/lizzz49/sflow"
)

type TransitionDefinition struct {
	Definition
	From       int     `json:"from"`
	To         int     `json:"to"`
	AlwaysTrue bool    `json:"always_true"`
	Express    Express `json:"express"`
}

type Express struct {
	Key   string      `json:"key"` //can load data from process context
	OP    string      `json:"op"`
	Value sflow.Value `json:"value"`
}

func newTransitionDefinition(id int, name string, from, to int, exp Express) *TransitionDefinition {
	t := &TransitionDefinition{}
	t.Id = id
	t.Name = name
	t.From = from
	t.To = to
	t.AlwaysTrue = exp.Key == ""
	t.Express = exp
	return t
}
