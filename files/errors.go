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

type ErrorInvalidMonth struct{}

func (e *ErrorInvalidMonth) Error() string {
	return "Error: Invalid month"
}

type ErrorInvalidFileFormat struct {}

func (e *ErrorInvalidFileFormat) Error() string {
	return "Error: Can't reconstruct file. Incorrect format"
}

type ErrorInvalidAmountOfFields struct{}

func (e *ErrorInvalidAmountOfFields) Error() string {
	return "Error: Invalid amount of fields"
}

func compareErrors(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}
	if e2 == nil && e1 != nil {
		return false
	}
	if e1 == nil && e2 != nil {
		return false
	}
	return e1.Error() == e2.Error()
}