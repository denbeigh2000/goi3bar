package util

import (
	"fmt"
	"time"
)

type DeprecationError struct{}

func (d DeprecationError) Error() string {
	return fmt.Sprintf("We don't do things that way anymore. I mean come on, it's %v!", time.Now().Year())
}
