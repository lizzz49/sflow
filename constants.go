package sflow

const (
	ProcessInstanceStatusNone = iota
	ProcessInstanceStatusNew
	ProcessInstanceStatusStarted
	ProcessInstanceStatusSuspended
	ProcessInstanceStatusTerminated
	ProcessInstanceStatusError
	ProcessInstaneStatusFinish
)

const (
	StartActivity = iota
	EndActivity
	BusinessActivity
)
