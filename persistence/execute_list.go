package persistence

import (
	"os"

	"github.com/pivotalservices/gtils/command"
	"github.com/xchapter7x/lo"
)

func executeList(callList []string, caller command.Executer) (err error) {
	for _, callstring := range callList {
		lo.G.Debug(callstring)

		if err = caller.Execute(os.Stdout, callstring); err != nil {
			lo.G.Error(err.Error())
			break
		}
	}
	return
}
