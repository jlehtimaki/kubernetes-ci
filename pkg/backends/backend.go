package backends

import "os/exec"

type Backend interface {
	Login() []*exec.Cmd
}

type BaseBackend struct {
}