package profilex

import (
	"os"

	"github.com/pkg/profile"
)

type noop struct{}

// Stop is a noop.
func (p *noop) Stop() {}

// Profile parses the PROFILING environment variable and executes the proper profiling task.
func Profile() interface {
	Stop()
} {
	switch os.Getenv("PROFILING") {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.NoShutdownHook)
	case "mem":
		return profile.Start(profile.MemProfile, profile.NoShutdownHook)
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.NoShutdownHook)
	case "block":
		return profile.Start(profile.BlockProfile, profile.NoShutdownHook)
	}
	return new(noop)
}
