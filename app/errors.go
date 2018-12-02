package app

//ErrPasswordDoNotMatch ...
type ErrPasswordDoNotMatch string

func (e ErrPasswordDoNotMatch) Error() string {
	return string(e)
}

//ErrNotChangePoll ...
type ErrNotChangePoll string

func (e ErrNotChangePoll) Error() string {
	return string(e)
}
