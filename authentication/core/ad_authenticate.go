package authentication

import "pulsecheck/authentication/model"

func (backend *JWTAuthenticationBackend) Authenticate(user *authmodel.User) bool {
	//	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)
	//
	testUser := authmodel.User{
		//	UUID:     uuid.New(),
		Username: "gokul",
		Password: "test",
	}

	if user.Username == testUser.Username {
		return true
	}
	return false
	//	//return user.Username == testUser.Username && bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte(user.Password)) == nil
	//}
}
