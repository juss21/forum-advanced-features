package main

import (
	"fmt"
	"os"
)

func errorCheck(err error, exit bool) {
	if err != nil {
		fmt.Println(err)
		if exit {
			os.Exit(0)
		}
	}
}

func getLogin(uid string, password string) bool {
	for i := 0; i < len(userlist); i++ {
		// if username or email correct
		if userlist[i].Username == uid || userlist[i].Email == uid {
			// if password correct
			if password == userlist[i].Password {
				return true
			}
		}
	}
	return false
}

func getRegister(uid string, password string, email string) (bool, string) {
	for i := 0; i < len(userlist); i++ {
		// TODO add email check aswell?
		// check if user already exists
		if userlist[i].Username == uid {
			if userlist[i].Email == email {
				return false, "This e-mail is already used!"
			}
			return false, "This username is taken!"
		}
	}
	return true, "Account has been created!"
}
