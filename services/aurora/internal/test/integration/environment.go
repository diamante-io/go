//lint:file-ignore U1001 Ignore all unused code, this is only used in tests.
package integration

import (
	"os"
)

type EnvironmentManager struct {
	oldEnvironment, newEnvironment map[string]string
}

func NewEnvironmentManager() *EnvironmentManager {
	env := &EnvironmentManager{}
	env.oldEnvironment = make(map[string]string)
	env.newEnvironment = make(map[string]string)
	return env
}

// Add sets a new environment variable, saving the original value (if any).
func (envmgr *EnvironmentManager) Add(key, value string) error {
	// If someone pushes an environmental variable more than once, we don't want
	// to lose the *original* value, so we're being careful here.
	if _, ok := envmgr.newEnvironment[key]; !ok {
		if oldValue, ok := os.LookupEnv(key); ok {
			envmgr.oldEnvironment[key] = oldValue
		}
	}

	envmgr.newEnvironment[key] = value
	return os.Setenv(key, value)
}

// Restore restores the environment prior to any modifications.
//
// You should probably use this alongside `defer` to ensure the global
// environment isn't modified for longer than you intend.
func (envmgr *EnvironmentManager) Restore() {
	for key := range envmgr.newEnvironment {
		if oldValue, ok := envmgr.oldEnvironment[key]; ok {
			os.Setenv(key, oldValue)
		} else {
			os.Unsetenv(key)
		}
	}
}
