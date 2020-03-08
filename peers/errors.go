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
