package peers

//CouldNotBeVerified should be returned
//if something that has to be sent by a
//participating peer could not be verified.
type CouldNotBeVerified struct {
}

func (e *CouldNotBeVerified) Error() string {
	return "This could not be verified to a participating peer"
}

//NotVerifiedError returns a not verified error.
func NotVerifiedError() error {
	return &CouldNotBeVerified{}
}


type ErrorCouldNotAppend struct {
}

func (e *ErrorCouldNotAppend) Error() string {
	return "Could not add the backup"
}


type ErrorIncorrectFormat struct {

}

func (e *ErrorIncorrectFormat) Error() string {
	return "Incorrect format of data"
}

type ErrorCouldNotDecode struct {
}

func (e *ErrorCouldNotDecode) Error() string {
	return "Could not decode the given data"
}

type ErrorFailedProtocol struct {

}

func (e *ErrorFailedProtocol) Error() string {
	return "Something went wrong in the protocol execution"
}

type ErrorUnableToProveStorage struct {
}

func (e *ErrorUnableToProveStorage) Error() string {
	return "Error: Unable to prove they have the backup"
}
