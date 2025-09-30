package models

type ParamError struct{}
type InvalidArticleError struct{}
type EmailOrUserTakenError struct{}
type InvalidLoginError struct{}
type InvalidTokenError struct{}
type UnauthorizedError struct{}
type DBInternalError struct{}

func (e *ParamError) Error() string {
	return "some request parameters are invalid or missing"
}

func (e *InvalidArticleError) Error() string {
	return "the article does not exist"
}

func (e *EmailOrUserTakenError) Error() string {
	return "email or username already taken"
}

func (e *InvalidLoginError) Error() string {
	return "invalid username or password"
}

func (e *InvalidTokenError) Error() string {
	return "invalid auth token supplied"
}

func (e *UnauthorizedError) Error() string {
	return "unauthorized"
}

func (e *DBInternalError) Error() string {
	return "an internal error occurred in the db module"
}
