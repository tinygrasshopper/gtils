package bosh

import (
	"errors"
)

var (
	//ErrorTaskStatusCode -
	ErrorTaskStatusCode = errors.New("The resp code from return task should return 200")
	//ErrorManifestStatusCode -
	ErrorManifestStatusCode = errors.New("The retriveing bosh manifest API response code is not equal to 200")
	//ErrorTaskRedirectStatusCode -
	ErrorTaskRedirectStatusCode = errors.New("The resp code after task creation should return 302")
	//ErrorTaskResultUnknown -
	ErrorTaskResultUnknown = errors.New("TASK processed result is unknown")
)
