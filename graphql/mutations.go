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
