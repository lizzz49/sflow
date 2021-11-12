package sflow

const (
	StatusNone = iota
	StatusNew
	StatusStarted
	StatusSuspended
	StatusTerminated
	StatusError
	StatusFinish
)

const (
	StartActivity = iota
	EndActivity
	BusinessActivity
)
