package begonia

type MethodNotAllowError struct{}

func (err MethodNotAllowError) Error() string {
	return "405 method not allow"
}

type NotFoundError struct{}

func (err NotFoundError) Error() string {
	return "404 not found"
}
