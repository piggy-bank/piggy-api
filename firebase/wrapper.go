package firebase

import (
	"firebase.google.com/go/auth"
)

func GetFireBaseAuthClient() *auth.Client {
	// configure firebase
	return SetupFirebase()
}
