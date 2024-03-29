package sflow

const (
	Int64Type = iota
	Float64Type
	StringType //bool can exchange to int[0,1] or string[true,false]
)

type ProcessContext map[string]Value ////string float64 int64 bool struct

type Value struct {
	Key  string `json:"key"`
	Type int    `json:"type"`
	Data string `json:"data"`
}
