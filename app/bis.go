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

//Visit ...
func Visit(helper HTTPHelper, userHandler UserHandler, sessionHandler SessionHandler) {
	createAnonUser := func(v interface{}) (interface{}, error) {
		return userHandler.CreateAnonUser(), nil
	}

	createSession := func(v interface{}) (interface{}, error) {
		return sessionHandler.CreateSession(v.(User).ID), nil
	}

	helper.Process(nil, createAnonUser, createSession)
}
