package apperror

const AccountNameAlreadyTaken = 0

type ApplicationError struct {
	Code int
}

func (e ApplicationError) Error() string {
	return ""
}
