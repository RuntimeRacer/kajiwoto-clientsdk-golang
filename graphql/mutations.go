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

/*
GraphQL for requests
*/

type kajiwotoLoginUserPWMutation struct {
	Login   Login `graphql:"login (usernameOrEmail: $usernameOrEmail, password: $password, deviceType: WEB)"`
	Welcome Welcome
}

type kajiwotoLoginAuthTokenMutation struct {
	Login   Login `graphql:"loginWithToken (authToken: $authToken, action: $action, deviceType: WEB)"`
	Welcome Welcome
}

type kajiwotoDatasetAITrainerGroupQuery struct {
	AITrainerGroup AITrainerGroup `graphql:"aiTrainerGroup (aiTrainerGroupId: $aiTrainerGroupId)"`
}

type kajiwotoDatasetLinesQuery struct {
	DatasetLines []DatasetLine `graphql:"datasetLines (aiTrainerGroupId: $aiTrainerGroupId, searchQuery: $searchQuery, limit: $limit, offset: $offset )"`
}

type kajiwotoAddToDatasetMutation struct {
	AIEditorResult AIEditorResult `graphql:"addToDataset (aiTrainerGroupId: $aiTrainerGroupId, editorType: $editorType, generateResults: $generateResults, dialogues: $dialogues )"`
}
