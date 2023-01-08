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
	"context"
	"fmt"
	gql "github.com/runtimeracer/go-graphql-client"
	"net/http"
)

// headerTransport is used to add custom headers to the request
// shootout to tgwizard; https://github.com/shurcooL/graphql/issues/28
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req)
	for key, val := range h.headers {
		req2.Header.Set(key, val)
	}
	return h.base.RoundTrip(req2)
}

func (h *headerTransport) GetHeaders() map[string]string {
	return h.headers
}

func (h *headerTransport) AddHeaders(newHeaders map[string]string) {
	for k, v := range newHeaders {
		h.headers[k] = v
	}
}

// KajiwotoClient is a custom graphql client for kajiwoto reqeusts
type KajiwotoClient struct {
	client          *gql.Client
	transportClient *http.Client
}

func GetKajiwotoClient(endpoint string) *KajiwotoClient {
	// Init HTTP Client
	transportClient := &http.Client{
		Transport: &headerTransport{
			base: http.DefaultTransport,
			headers: map[string]string{
				// Default Headers
				"Content-Type": "application/json",
			},
		},
	}

	return &KajiwotoClient{
		client:          gql.NewClient(endpoint, transportClient),
		transportClient: transportClient,
	}
}

func (c *KajiwotoClient) GetHeaders() map[string]string {
	return c.transportClient.Transport.(*headerTransport).GetHeaders()
}

func (c *KajiwotoClient) AddHeaders(newHeaders map[string]string) {
	c.transportClient.Transport.(*headerTransport).AddHeaders(newHeaders)
}

// DoLoginUserPW performs login via user / pw combination
func (c *KajiwotoClient) DoLoginUserPW(username, password string) (result LoginResult, err error) {
	// Sanity check
	if username == "" || password == "" {
		return result, fmt.Errorf("invalid login credentials")
	}

	vars := map[string]interface{}{
		"usernameOrEmail": gql.String(username),
		"password":        gql.String(password),
	}

	loginResult := kajiwotoLoginUserPWMutation{}
	if errLogin := c.performGraphMutation(vars, &loginResult); errLogin != nil {
		return result, errLogin
	}

	// Build generic Result object
	result = LoginResult{
		Login:   loginResult.Login,
		Welcome: loginResult.Welcome,
	}

	return result, nil
}

// DoLoginAuthToken performs login via session key if available
func (c *KajiwotoClient) DoLoginAuthToken(authToken string) (result LoginResult, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid login credentials")
	}

	vars := map[string]interface{}{
		"authToken": gql.Token(authToken),
		"action":    gql.String(""),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	loginResult := kajiwotoLoginAuthTokenMutation{}
	if errLogin := c.performGraphMutation(vars, &loginResult); errLogin != nil {
		return result, fmt.Errorf("unable to login, response: %q", errLogin)
	}

	// Build generic Result object
	result = LoginResult{
		Login:   loginResult.Login,
		Welcome: loginResult.Welcome,
	}

	return result, nil
}

func (c *KajiwotoClient) GetAITrainerGroup(aiTrainerGroupID, authToken string) (result AITrainerGroup, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid auth token")
	}
	if aiTrainerGroupID == "" {
		return result, fmt.Errorf("invalid trainer group ID")
	}

	vars := map[string]interface{}{
		"aiTrainerGroupId": gql.String(aiTrainerGroupID),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	// Execute Query
	aiTrainerGroupResult := kajiwotoDatasetAITrainerGroupQuery{}
	if errLogin := c.performGraphQuery(vars, &aiTrainerGroupResult); errLogin != nil {
		return result, fmt.Errorf("unable to fetch AI trainer group, response: %q", errLogin)
	}

	// Build generic Result object
	result = aiTrainerGroupResult.AITrainerGroup
	return result, nil
}

func (c *KajiwotoClient) GetDatasetLines(aiTrainerGroupID, searchQuery, authToken string, limit, offset int) (result []DatasetLine, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid auth token")
	}
	if aiTrainerGroupID == "" {
		return result, fmt.Errorf("invalid trainer group ID")
	}
	if limit < 1 || limit > 100 {
		return result, fmt.Errorf("limit exceeds allowed range")
	}
	if offset < 0 {
		return result, fmt.Errorf("offset cannot be negative")
	}

	vars := map[string]interface{}{
		"aiTrainerGroupId": gql.String(aiTrainerGroupID),
		"searchQuery":      gql.String(searchQuery),
		"limit":            gql.Int(limit),
		"offset":           gql.Int(offset),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	// Execute Query
	datasetLinesResult := kajiwotoDatasetLinesQuery{}
	if errQuery := c.performGraphQuery(vars, &datasetLinesResult); errQuery != nil {
		return result, fmt.Errorf("unable to fetch dataset lines, response: %q", errQuery)
	}

	// Build generic Result object
	result = datasetLinesResult.DatasetLines
	return result, nil
}

func (c *KajiwotoClient) AddToDataset(aiTrainerGroupID, authToken string, dialogues []*AiDialogueInput) (result AIEditorResult, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid login credentials")
	}

	vars := map[string]interface{}{
		"aiTrainerGroupId": gql.String(aiTrainerGroupID),
		"dialogues":        dialogues,
		"editorType":       gql.String("kajitool"),
		"generateResults":  gql.Boolean(false),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	trainingResult := kajiwotoAddToDatasetMutation{}
	if errTrain := c.performGraphMutation(vars, &trainingResult); errTrain != nil {
		return result, fmt.Errorf("unable to train dataset, response: %q", errTrain)
	}

	// Build generic Result object
	result = trainingResult.AIEditorResult
	return result, nil
}

func (c *KajiwotoClient) performGraphMutation(vars map[string]interface{}, mutation interface{}) error {
	return c.client.Mutate(context.Background(), mutation, vars)
}

func (c *KajiwotoClient) performGraphQuery(vars map[string]interface{}, query interface{}) error {
	return c.client.Query(context.Background(), query, vars)
}

// cloneRequest creates a shallow copy of the request along with a deep copy of the Headers.
func cloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = cloneHeader(req.Header)

	return r
}

// cloneHeader creates a deep copy of an http.Header.
func cloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}
