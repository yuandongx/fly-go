package executor

const (
	StatusIdle     = "idle"
	StatusRunning  = "running"
	StatusStopped  = "stopped"
	StatusError    = "error"
	StatusSuccess  = "success"
	StatusTimeout  = "timeout"
	StatusPending  = "pending"
	StatusCanceled = "canceled"
)

const (
	TriggerInterval = "interval"
	TriggerOnce     = "once"
)
