package system

import "os/user"

func GetCurrentUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "Unknown user", err
	}

	return currentUser.Username, nil
}
