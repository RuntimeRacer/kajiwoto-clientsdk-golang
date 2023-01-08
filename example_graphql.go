package main

import (
	"errors"
	"fmt"
	"github.com/runtimeracer/kajiwoto-clientsdk-golang/constants"
	"github.com/runtimeracer/kajiwoto-clientsdk-golang/graphql"
	"os"
)

func main() {
	username := ""
	password := ""
	sessionKey := ""

	// Init Client
	client := graphql.GetKajiwotoClient(constants.KWGraphQLEndpoint)

	// Check whether there is a Session key defined
	loginResult := graphql.LoginResult{}
	var errLogin error
	if sessionKey != "" {
		fmt.Println(fmt.Sprintf("Performing login via Session key: %v", sessionKey))
		loginResult, errLogin = client.DoLoginAuthToken(sessionKey)
		if errLogin != nil {
			fmt.Println(fmt.Sprintf("Unable to login via auth token, trying with username / password. error: %v", errLogin))
			loginResult, errLogin = client.DoLoginUserPW(username, password)
		} else if loginResult.Login.AuthToken == "" {
			fmt.Println(fmt.Sprintf("No User information returned from server. Session may be outdated. Trying with username / password."))
			loginResult, errLogin = client.DoLoginUserPW(username, password)
		}
	} else {
		fmt.Println("Performing login via Username / Password combo")
		loginResult, errLogin = client.DoLoginUserPW(username, password)
	}

	// Check for error
	if errLogin != nil {
		fmt.Println(fmt.Errorf("unable to login, response: %q", errLogin))
		os.Exit(1)
	}

	// Validate response
	if loginResult.Login.AuthToken == "" {
		fmt.Println(errors.New("invalid response from server: Auth token empty"))
		os.Exit(1)
	}

	// Seems like Login worked
	userInfo := &loginResult.Login.User
	fmt.Println(fmt.Sprintf("Login successful! Hello %v!", userInfo.DisplayName))

	// Update Auth token in config file
	sessionKey = loginResult.Login.AuthToken

	// ... write your cfg ...

}
