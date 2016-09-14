package pushbullet

import "os"
import "fmt"

// EnvGetter is something that reads an environment variable. In code, this
// could read an environment variable from the OS. In tests, this can read
// variables from a mock.
type EnvGetter interface {
	Get(v string) string
}

// OSEnv is an EnvGetter that reads real OS environment variables.
type OSEnv struct{}

// Get reads an environment variable from the OS.
func (e OSEnv) Get(v string) string {
	return os.Getenv(v)
}

// MockEnv is an EnvGetter that mocks the OS environment.
type MockEnv map[string]string

// Get looks up a variable from the MockEnv map.
func (e MockEnv) Get(v string) string {
	return e[v]
}

// ConfigErrror is a configuration error for noti packages.
type ConfigErrror struct {
	Env    string
	Reason string
}

func (e ConfigErrror) Error() string {
	return fmt.Sprintf("invalid configuration for %s: %s", e.Env, e.Reason)
}

// APIError is an API error that's returned if a notification API request
// failed.
type APIError struct {
	Site string
	Msg  string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%s API: %s", e.Site, e.Msg)
}

type Params struct {
	Title   string
	Message string

	Failure bool
	API     string
	Token  string
	Config  EnvGetter
}
