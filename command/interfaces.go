package command

import "io"

//Executer -
type Executer interface {
	Execute(destination io.Writer, command string) error
}
