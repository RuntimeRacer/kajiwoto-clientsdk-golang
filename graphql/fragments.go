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
	gql "github.com/runtimeracer/go-graphql-client"
)

type Plus struct {
	ExpireAt  uint64
	Cancelled gql.Boolean
	Icon      gql.Int
	Coins     gql.Int
	Type      gql.String
}

type Creator struct {
	AllowSubscriptions gql.Boolean
	DatasetTags        []gql.String
}

type Profile struct {
	ID          gql.String
	FirstName   gql.String
	LastName    gql.String
	Description gql.String
	Gender      gql.String
	Birthday    gql.String
	PhotoUri    gql.String
}

type Email struct {
	Address  gql.String
	Verified gql.Boolean
}

type Settings struct {
	PersonalRoomOrder []gql.String
	FavoriteRoomIds   []gql.String
	FavoriteEmojis    []gql.String
}

type User struct {
	ID          gql.String
	Activated   gql.Boolean
	Moderator   gql.Boolean
	Username    gql.String
	DisplayName gql.String
	Plus        Plus
	Creator     Creator
	Profile     Profile
	Email       Email
}

type Login struct {
	AuthToken string
	User      User
	Usage     Usage
	Settings  Settings
}

type Usage struct {
	Generator gql.Int
}

type Announcement struct {
	Date      uint64
	Title     gql.String
	Emojis    gql.String
	Content   []gql.String
	TextColor gql.String
}

type Welcome struct {
	WebVersion   gql.String
	Announcement Announcement
}

type LoginResult struct {
	Login   Login
	Welcome Welcome
}

type AITrainerGroup struct {
	ID              gql.String
	Name            gql.String
	Count           gql.Int
	Deleted         gql.Boolean
	Description     gql.String
	Documents       []AIDocument
	DominantColors  []gql.String
	Kudos           Kudos
	NSFW            gql.Boolean
	Personalities   [][]gql.String
	PetSpeciesIds   []gql.String
	Price           gql.Int
	ProfilePhotoUri gql.String
	Purchased       gql.Boolean
	Status          gql.String
	Tags            []gql.String
	UpdatedAt       uint64
	User            User
}

type Kudos struct {
	ID       gql.String
	Upvoted  gql.Boolean
	Upvotes  gql.Int
	Comments gql.Int
}

type DatasetLine struct {
	ID               gql.String
	UserMessage      gql.String
	Message          gql.String
	ASM              gql.String
	Endearment       gql.String
	Recent           gql.String
	Time             gql.String
	Deleted          gql.Boolean
	History          []gql.String
	AITrainerGroupID gql.String
}

type AIDocument struct {
	ID          gql.String
	Order       gql.Int
	Title       gql.String
	Content     gql.String
	QueueStatus gql.String
	QueuedAt    uint64
	BuiltAt     uint64
	CreatedAt   uint64
	UpdatedAt   uint64
}

type AiDialogueInput struct {
	Conditions  AITrainingCondition `json:"conditions"`
	Generated   gql.Boolean         `json:"generated"`
	History     []gql.String        `json:"history"`
	Message     gql.String          `json:"message"`
	UserMessage gql.String          `json:"userMessage"`
}

type AITrainingCondition struct {
	ASM        *string `json:"asm"`
	Endearment *string `json:"endearment"`
	Recent     *string `json:"recent"`
	Time       *string `json:"time"`
}

type AIEditorResult struct {
	Added            []DatasetLine
	AITrainerGroupID gql.String
	Count            gql.Int
	DeletedIDs       []gql.String
	Generated        []DatasetLine
	Message          gql.String
	MessageType      gql.String
}
