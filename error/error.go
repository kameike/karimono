package apperror

const ErrorAccountNameAlreadyTaken = 0
const ErrorInvalidAccessToken = 1
const ErrorInvalidTeamPassword = 2
const ErrorInvalidAccountName = 3
const ErrorDataNotFount = 4
const ErrorTeamNameAlreadyTaken = 5
const ErrorAlreadyJoin = 6

type ApplicationError struct {
	Code int
}

func (e ApplicationError) Error() string {
	return ""
}

func (e ApplicationError) StatusCode() int {
	return 400
}
