package mock

import "io"

type Executor struct {
	Command string
	Args    []string

	Output string
	Error  error
}

func (e *Executor) Execute(name string, args ...string) (string, error) {
	e.Command = name
	e.Args = args

	return e.Output, e.Error
}

func (e *Executor) ExecuteStdout(name string, args ...string) error {
	e.Command = name
	e.Args = args

	return e.Error
}

func (e *Executor) ExecuteWithWriters(_, _ io.Writer, name string, args ...string) error {
	e.Command = name
	e.Args = args

	return e.Error
}
