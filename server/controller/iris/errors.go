package iris

type internalError struct {
	err error
}

func (err *internalError) Error() string {
	return err.err.Error()
}

func InternalError(err error) error {
	return &internalError{
		err: err,
	}
}
