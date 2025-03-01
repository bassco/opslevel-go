package opslevel_test

import (
	"testing"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/rocktavious/autopilot/v2023"
)

func TestCreateAWSIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation AWSIntegrationCreate($input:AwsIntegrationInput!){awsIntegrationCreate(input: $input){integration{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},errors{message,path}}}`,
		`{"input": { "iamRole": "arn:aws:iam::XXXX:role/aws-integration-role", "externalId": "123456789", "ownershipTagKeys": ["owner"] }}`,
		`{"data": {
      "awsIntegrationCreate": {
        "integration": {
          {{ template "id1" }},
          "name": "AWS - XXXX",
          "type": "aws",
          "createdAt": "2023-04-26T16:25:29.574450Z",
          "installedAt": "2023-04-26T16:25:28.541124Z",
          "iamRole": "arn:aws:iam::XXXX:role/aws-integration-role",
          "externalId": "123456789",
          "awsTagsOverrideOwnership": true,
          "ownershipTagKeys": [
            "owner"
          ]
        },
        "errors": []
      }}}`,
	)
	client := BestTestClient(t, "integration/create_aws", testRequest)
	// Act
	result, err := client.CreateIntegrationAWS(opslevel.AWSIntegrationInput{
		IAMRole:    opslevel.RefOf("arn:aws:iam::XXXX:role/aws-integration-role"),
		ExternalID: opslevel.RefOf("123456789"),
	})
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "AWS - XXXX", result.Name)
}

func TestCreateAzureResourcesIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation AzureResourcesIntegrationCreate($input:AzureResourcesIntegrationInput!){azureResourcesIntegrationCreate(input: $input){integration{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},errors{message,path}}}`,
		`{"input": { "name": "new azure resources", "tenantId": "12345678-1234-1234-1234-123456789abc", "subscriptionId": "12345678-1234-1234-1234-123456789def", "clientId": "XXX_client_id_XXX", "clientSecret": "XXX_client_secret_XXX"}}`,
		`{"data": {
      "azureResourcesIntegrationCreate": {
        "integration": {
          {{ template "id1" }},
          "name": "new azure resources",
          "type": "azureResources",
          "createdAt": "2024-07-04T16:25:29.574450Z",
          "installedAt": "2024-07-04T16:25:28.541124Z",
          "tenantId": "12345678-1234-1234-1234-123456789abc",
          "subscriptionId": "12345678-1234-1234-1234-123456789def",
          "aliases": []
        },
        "errors": []
      }}}`,
	)
	client := BestTestClient(t, "integration/create_azure_resources", testRequest)
	// Act
	result, err := client.CreateIntegrationAzureResources(opslevel.AzureResourcesIntegrationInput{
		Name:           opslevel.RefOf("new azure resources"),
		TenantId:       opslevel.RefOf("12345678-1234-1234-1234-123456789abc"),
		SubscriptionId: opslevel.RefOf("12345678-1234-1234-1234-123456789def"),
		ClientId:       opslevel.RefOf("XXX_client_id_XXX"),
		ClientSecret:   opslevel.RefOf("XXX_client_secret_XXX"),
	})
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "new azure resources", result.Name)
	autopilot.Equals(t, "12345678-1234-1234-1234-123456789abc", result.TenantId)
	autopilot.Equals(t, "12345678-1234-1234-1234-123456789def", result.SubscriptionId)
	autopilot.Equals(t, []string{}, result.Aliases)
}

func TestCreateNewRelicIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation NewRelicIntegrationCreate($input:NewRelicIntegrationInput!){newRelicIntegrationCreate(input: $input){integration{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},errors{message,path}}}`,
		`{ "input": { "apiKey": "123456789", "baseUrl": "https://api.newrelic.com/graphql" }}`,
		`{"data": {
      "newRelicIntegrationCreate": {
        "integration": {
          {{ template "id1" }},
          "name": "New Relic - XXXX",
          "type": "new_relic",
          "createdAt": "2023-04-26T16:25:29.574450Z",
          "installedAt": "2023-04-26T16:25:28.541124Z",
          "accountKey": "XXXX",
          "baseUrl": "https://api.newrelic.com/graphql"
        },
        "errors": []
      }}}`,
	)
	client := BestTestClient(t, "integration/create_new_relic", testRequest)
	// Act
	result, err := client.CreateIntegrationNewRelic(opslevel.NewRelicIntegrationInput{
		ApiKey:  opslevel.RefOf("123456789"),
		BaseUrl: opslevel.RefOf("https://api.newrelic.com/graphql"),
	})
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "New Relic - XXXX", result.Name)
}

func TestGetIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`query IntegrationGet($id:ID!){account{integration(id: $id){id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}}}}`,
		`{ {{ template "id1" }} }`,
		`{"data": {
      "account": {
        "integration": {
          {{ template "id1" }},
          "name": "Deploy",
          "type": "deploy"
        }
      }}}`,
	)
	client := BestTestClient(t, "integration/get", testRequest)
	// Act
	result, err := client.GetIntegration(id1)
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "Deploy", result.Name)
}

func TestGetMissingIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`query IntegrationGet($id:ID!){account{integration(id: $id){id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}}}}`,
		`{ {{ template "id2" }} }`,
		`{"data": { "account": { "integration": null }}}`,
	)
	client := BestTestClient(t, "integration/get_missing", testRequest)
	// Act
	_, err := client.GetIntegration(id2)
	// Assert
	autopilot.Assert(t, err != nil, "This test should throw an error.")
}

func TestListIntegrations(t *testing.T) {
	// Arrange
	testRequestOne := autopilot.NewTestRequest(
		`query IntegrationList($after:String!$first:Int!){account{integrations(after: $after, first: $first){nodes{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},pageInfo{hasNextPage,hasPreviousPage,startCursor,endCursor},totalCount}}}`,
		`{{ template "pagination_initial_query_variables" }}`,
		`{ "data": { "account": { "integrations": { "nodes": [ { {{ template "deploy_integration_response" }} }, { {{ template "payload_integration_response" }} } ], {{ template "pagination_initial_pageInfo_response" }}, "totalCount": 2 }}}}`,
	)
	testRequestTwo := autopilot.NewTestRequest(
		`query IntegrationList($after:String!$first:Int!){account{integrations(after: $after, first: $first){nodes{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},pageInfo{hasNextPage,hasPreviousPage,startCursor,endCursor},totalCount}}}`,
		`{{ template "pagination_second_query_variables" }}`,
		`{ "data": { "account": { "integrations": { "nodes": [ { {{ template "kubernetes_integration_response" }} } ], {{ template "pagination_second_pageInfo_response" }}, "totalCount": 1 }}}}`,
	)
	requests := []autopilot.TestRequest{testRequestOne, testRequestTwo}

	client := BestTestClient(t, "integration/list", requests...)
	// Act
	response, err := client.ListIntegrations(nil)
	result := response.Nodes
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, 3, response.TotalCount)
	autopilot.Equals(t, "Payload", result[1].Name)
	autopilot.Equals(t, "Kubernetes", result[2].Name)
}

func TestUpdateAWSIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation AWSIntegrationUpdate($input:AwsIntegrationInput!$integration:IdentifierInput!){awsIntegrationUpdate(integration: $integration input: $input){integration{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},errors{message,path}}}`,
		`{"integration": { {{ template "id1" }} }, "input": { "name": "Dev2", "externalId": "123456789", "ownershipTagKeys": null }}`,
		`{"data": {
      "awsIntegrationUpdate": {
        "integration": {
        {{ template "id1" }},
        "name": "Dev2",
        "type": "aws",
        "createdAt": "2023-04-26T16:25:29.574450Z",
        "installedAt": "2023-04-26T16:25:28.541124Z",
        "iamRole": "arn:aws:iam::XXXX:role/aws-integration-role",
        "externalId": "123456789",
        "awsTagsOverrideOwnership": true,
        "ownershipTagKeys": [
          "owner"
        ]
      },
      "errors": []
    }}}`,
	)
	client := BestTestClient(t, "integration/update_aws", testRequest)
	// Act
	result, err := client.UpdateIntegrationAWS(string(id1), opslevel.AWSIntegrationInput{
		Name:       opslevel.RefOf("Dev2"),
		ExternalID: opslevel.RefOf("123456789"),
	})
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "Dev2", result.Name)
}

func TestUpdateAzureResourcesIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation AzureResourcesIntegrationUpdate($input:AzureResourcesIntegrationInput!$integration:IdentifierInput!){azureResourcesIntegrationUpdate(integration: $integration input: $input){integration{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},errors{message,path}}}`,
		`{"integration": { {{ template "id1" }} }, "input": { "name": "updated azure resources", "clientId": "updated client id", "clientSecret": "updated client secret" }}`,
		`{"data": {
      "azureResourcesIntegrationUpdate": {
        "integration": {
          {{ template "id1" }},
          "name": "updated azure resources",
          "type": "azureResources",
          "createdAt": "2024-07-04T16:25:29.574450Z",
          "installedAt": "2024-07-04T16:25:28.541124Z",
          "tenantId": "12345678-1234-1234-1234-123456789abc",
          "subscriptionId": "12345678-1234-1234-1234-123456789def",
          "aliases": ["alias1", "alias2"]
      },
      "errors": []
    }}}`,
	)
	client := BestTestClient(t, "integration/update_azure_resources", testRequest)
	// Act
	result, err := client.UpdateIntegrationAzureResources(string(id1), opslevel.AzureResourcesIntegrationInput{
		Name:         opslevel.RefOf("updated azure resources"),
		ClientId:     opslevel.RefOf("updated client id"),
		ClientSecret: opslevel.RefOf("updated client secret"),
	})
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "updated azure resources", result.Name)
	autopilot.Equals(t, "12345678-1234-1234-1234-123456789abc", result.TenantId)
	autopilot.Equals(t, "12345678-1234-1234-1234-123456789def", result.SubscriptionId)
	autopilot.Equals(t, []string{"alias1", "alias2"}, result.Aliases)
}

func TestUpdateNewRelicIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation NewRelicIntegrationUpdate($input:NewRelicIntegrationInput!$resource:IdentifierInput!){newRelicIntegrationUpdate(input: $input resource: $resource){integration{id,name,type,createdAt,installedAt,... on AwsIntegration{iamRole,externalId,awsTagsOverrideOwnership,ownershipTagKeys},... on AzureResourcesIntegration{aliases,ownershipTagKeys,subscriptionId,tagsOverrideOwnership,tenantId},... on NewRelicIntegration{baseUrl,accountKey}},errors{message,path}}}`,
		`{"resource": { {{ template "id1" }} }, "input": { "baseUrl": "https://api-test.newrelic.com/graphql" }}`,
		`{"data": {
      "newRelicIntegrationUpdate": {
        "integration": {
        {{ template "id1" }},
        "name": "New Relic - XXXX",
        "type": "new_relic",
        "createdAt": "2023-04-26T16:25:29.574450Z",
        "installedAt": "2023-04-26T16:25:28.541124Z",
        "accountKey": "XXXX",
        "baseUrl": "https://api-test.newrelic.com/graphql"
      },
      "errors": []
    }}}`,
	)

	client := BestTestClient(t, "integration/update_new_relic", testRequest)
	// Act
	result, err := client.UpdateIntegrationNewRelic(
		string(id1),
		opslevel.NewRelicIntegrationInput{
			BaseUrl: opslevel.RefOf("https://api-test.newrelic.com/graphql"),
		},
	)
	// Assert
	autopilot.Equals(t, nil, err)
	autopilot.Equals(t, id1, result.Id)
	autopilot.Equals(t, "https://api-test.newrelic.com/graphql", result.BaseUrl)
}

func TestDeleteIntegration(t *testing.T) {
	// Arrange
	testRequest := autopilot.NewTestRequest(
		`mutation IntegrationDelete($input:IdentifierInput!){integrationDelete(resource: $input){errors{message,path}}}`,
		`{"input": { {{ template "id1" }} }}`,
		`{"data": { "integrationDelete": { "errors": [] }}}`,
	)
	client := BestTestClient(t, "integration/delete", testRequest)
	// Act
	err := client.DeleteIntegration(string(id1))
	// Assert
	autopilot.Equals(t, nil, err)
}
