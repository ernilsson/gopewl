package gopewl

import "errors"

var (
	ErrNonPositivePoolSize = errors.New("pool size must be a positive integer")
	ErrNegativeQueueSize = errors.New("queue size must be a positive integer or 0")
	ErrIllegalPoolCapacity = errors.New("pool capacity must be larger than pool size or 0")
	ErrNegativePoolCapacity = errors.New("pool capacity must be a positive integer or 0")
)