package app

//CreateUser ...
func CreateUser(helper HTTPHelper, handler UserHandler) {
	createUser := func(v interface{}) (interface{}, error) {
		return handler.CreateUserFromData(v.(*UserCreationData))
	}

	saveUser := func(v interface{}) (interface{}, error) {
		return handler.SaveUser(v.(User)), nil
	}

	helper.Process(&UserCreationData{}, createUser, saveUser)
}
