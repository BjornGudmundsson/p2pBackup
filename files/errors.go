package files

type ErrorIncorrectLogFormat struct {
}

func (e *ErrorIncorrectLogFormat) Error() string {
	return "Error: Incorrect format for given log format. Could not decode"
}

type ErrorNoLogs struct{}

func (e *ErrorNoLogs) Error() string {
	return "Error: There were no logs"
}