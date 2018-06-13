package sflow

type TransitionDefinition struct {
	Definition
	From       string  `json:"from"`
	To         string  `json:"to"`
	AlwaysTrue bool    `json:"always_true"`
	Express    Express `json:"express"`
}

type Express struct {
	Key   string `json:"key"` //can load data from process context
	OP    string `json:"op"`
	Value Value  `json:"value"`
}

func newTransitionDefinition(id, name, from, to string, exp Express) *TransitionDefinition {
	t := &TransitionDefinition{}
	t.Id = id
	t.Name = name
	t.From = from
	t.To = to
	t.AlwaysTrue = exp.Key == ""
	t.Express = exp
	return t
}
