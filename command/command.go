package command

const (
	// ExitCodeOK means that this app successfully finished to run.
	ExitCodeOK int = iota
	// ExitCodeFailed means that this app failed to run. Please check the error message.
	ExitCodeFailed
	// ExitCodeCommandNotFound means that your inputed sub-command is invalid.
	ExitCodeCommandNotFound
	// ExitCodeFunctionError means that this app's any function cannot work fine.
	ExitCodeFunctionError
	// ExitCodeInvalidArguments means that the arguments are invalid.
	ExitCodeInvalidArguments
	// ExitCodeIOError means that this app cannot conplete any IO function.
	ExitCodeIOError
)
