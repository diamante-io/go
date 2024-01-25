package health

// Status indicates whether the service is health or not.
type Status string

const (
	// StatusPass indicates that the service is healthy.
	StatusPass Status = "pass"
	// StatusFail indicates that the service is unhealthy.
	StatusFail Status = "fail"
)
