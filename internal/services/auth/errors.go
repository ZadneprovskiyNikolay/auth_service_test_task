package auth

type UnauthorizedError struct {
	message string
}

func (err UnauthorizedError) Error() string {
	return err.message
}
