// Package websocket
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
package websocket

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

/*
 * handlers.go defines basic handlers for common events from the Kajiwoto backend
 * use these as inspiration for your own handler implementations when working with the SDK
 */

// NewKajiwotoWebSocketAuthResponseHandler is used to handle an auth message response
func NewKajiwotoWebSocketAuthResponseHandler(c *KajiwotoClient) MessageHandlerFunc {
	return func(message *KajiwotoWebSocketMessage) error {
		// Try to umarshall into required response
		// If this won't work, message is not of expected type
		response := &KaiwotoWebSocketAuthResponse{}
		if errUnmarshall := json.Unmarshal(message.MessageContent.([]byte), response); errUnmarshall != nil {
			return errUnmarshall
		}
		c.socketID = response.Sid
		log.Debugf("Assigned Socket ID: %v", c.socketID)
		return nil
	}
}

// NewKajiwotoWebSocketPingHandler is used to handle a ping event from the backend
func NewKajiwotoWebSocketPingHandler(c *KajiwotoClient) MessageHandlerFunc {
	return func(message *KajiwotoWebSocketMessage) error {
		if message.MessageCode == SocketCodePing {
			// Send Ping
			pongResponse := &KajiwotoWebSocketMessage{
				MessageCode:    SocketCodePong,
				MessageContent: nil,
			}
			if errPong := c.SendMessage(pongResponse); errPong != nil {
				return errPong
			}
		}
		return nil
	}
}
