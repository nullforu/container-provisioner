package stack

import "errors"

var (
	ErrNotFound            = errors.New("stack not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrPodSpecInvalid      = errors.New("invalid pod spec")
	ErrNoAvailableNodePort = errors.New("no available nodeport")
	ErrClusterSaturated    = errors.New("cluster saturated")
)
