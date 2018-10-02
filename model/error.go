package model

type notFoundError string

func (err notFoundError) Error() string {
	return string(err) + " not found"
}

func NotFoundError(typ string) error {
	return notFoundError(typ)
}

func IsNotFound(err error) bool {
	_, ok := err.(notFoundError)
	return ok
}
