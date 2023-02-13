// Package graphql
/*
Copyright © 2023 runtimeracer@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package graphql

import (
	"errors"
	"fmt"
	"github.com/runtimeracer/kajiwoto-clientsdk-golang/constants"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type GraphQLClientTestSuite struct {
	suite.Suite
}

func TestGraphQLClientTestSuite(t *testing.T) {
	suite.Run(t, new(GraphQLClientTestSuite))
}

func (s *GraphQLClientTestSuite) SetupTest() {
	// Set Log level for all tests
	log.SetLevel(log.DebugLevel)
}

func (s *GraphQLClientTestSuite) helperPerformLogin(client *KajiwotoGraphQLClient, username, password string) (sessionKey string, errLogin error) {
	// Check whether there is a Session key defined
	loginResult := LoginResult{}
	log.Debug("Performing login via Username / Password combo")
	loginResult, errLogin = client.DoLoginUserPW(username, password)

	// Check for error
	if errLogin != nil {
		return "", fmt.Errorf("unable to login, response: %q", errLogin)
	}

	// Validate response
	if loginResult.Login.AuthToken == "" {
		return "", errors.New("invalid response from server: Auth token empty")
	}

	// Seems like Login worked
	userInfo := &loginResult.Login.User
	log.Infof("Login successful! Hello %v!", userInfo.DisplayName)

	// Update Auth token in config file
	sessionKey = loginResult.Login.AuthToken
	return sessionKey, nil
}

func (s *GraphQLClientTestSuite) TestGraphQLLoginWrongCredentials() {
	// Init params
	username := os.Getenv("GRAPHQL_USER_LOGIN")
	password := os.Getenv("GRAPHQL_USER_PASSWORD")[1:4]

	// Init Client
	client := GetKajiwotoGraphQLClient(constants.KWGraphQLEndpoint)
	sessionKey, errLogin := s.helperPerformLogin(client, username, password)
	assert.NotNil(s.T(), errLogin)
	assert.Empty(s.T(), sessionKey)
}

func (s *GraphQLClientTestSuite) TestGraphQLLoginCorrectCredentials() {
	// Init params
	username := os.Getenv("GRAPHQL_USER_LOGIN")
	password := os.Getenv("GRAPHQL_USER_PASSWORD")

	// Init Client
	client := GetKajiwotoGraphQLClient(constants.KWGraphQLEndpoint)
	sessionKey, errLogin := s.helperPerformLogin(client, username, password)
	assert.Nil(s.T(), errLogin)
	assert.NotEmpty(s.T(), sessionKey)
}

func (s *GraphQLClientTestSuite) TestGraphQLGetRoom() {
	// Init params
	username := os.Getenv("GRAPHQL_USER_LOGIN")
	password := os.Getenv("GRAPHQL_USER_PASSWORD")
	roomID := os.Getenv("GRAPHQL_ROOM_ID")

	// Init Client
	client := GetKajiwotoGraphQLClient(constants.KWGraphQLEndpoint)
	sessionKey, errLogin := s.helperPerformLogin(client, username, password)
	assert.Nil(s.T(), errLogin)
	assert.NotEmpty(s.T(), sessionKey)

	// Get Room Data
	room, errRoom := client.GetRoom(roomID, "", sessionKey)
	assert.Nil(s.T(), errRoom)
	assert.NotEmpty(s.T(), room.ID)
	assert.Equal(s.T(), strings.ToLower(string(room.ChatRoomID)), strings.ToLower(roomID))
}

func (s *GraphQLClientTestSuite) TestGraphQLGetRoomHistory() {
	// Init params
	username := os.Getenv("GRAPHQL_USER_LOGIN")
	password := os.Getenv("GRAPHQL_USER_PASSWORD")
	roomID := os.Getenv("GRAPHQL_ROOM_ID")

	// Init Client
	client := GetKajiwotoGraphQLClient(constants.KWGraphQLEndpoint)
	sessionKey, errLogin := s.helperPerformLogin(client, username, password)
	assert.Nil(s.T(), errLogin)
	assert.NotEmpty(s.T(), sessionKey)

	// Get Room Data
	roomHistory, errRoom := client.GetRoomHistory(roomID, "", sessionKey)
	assert.Nil(s.T(), errRoom)
	assert.NotEmpty(s.T(), roomHistory.ID)
	assert.Equal(s.T(), strings.ToLower(string(roomHistory.ChatRoomID)), strings.ToLower(roomID))
	assert.NotNil(s.T(), roomHistory.Messages)
}
