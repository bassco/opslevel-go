package opslevel

import (
	"fmt"
	"github.com/shurcooL/graphql"
)

type CustomActionsId struct {
	Aliases []string   `graphql:"aliases"`
	Id      graphql.ID `graphql:"id"`
}

type CustomActionsExternalAction struct {
	CustomActionsId

	Description    string `graphql:"description"`
	LiquidTemplate string `graphql:"liquidTemplate"`
	Name           string `graphql:"name"`

	CustomActionsWebhookAction `graphql:"... on CustomActionsWebhookAction"`
}

type CustomActionsWebhookAction struct {
	BasicAuthUserName string                      `graphql:"basicAuthUserName"`
	BasicAuthPassword string                      `graphql:"basicAuthPassword"`
	Headers           JSON                        `graphql:"headers"`
	HTTPMethod        CustomActionsHttpMethodEnum `graphql:"httpMethod"`
	WebhookURL        string                      `graphql:"webhookUrl"`
}

type CustomActionsTriggerDefinition struct {
	Action      CustomActionsId `graphql: "action"`
	Aliases     []string        `graphql:"aliases"`
	Description string          `graphql:"description"`
	Filter      FilterId        `graphql:"filter"`
	Id          graphql.ID      `graphql:"id"`
	Name        string          `graphql:"name"`
	Owner       TeamId          `graphql:"owner"`
	Timestamps  Timestamps      `graphql:"timestamps"`
}

type CustomActionsExternalActionsConnection struct {
	Nodes      []CustomActionsExternalAction
	PageInfo   PageInfo
	TotalCount int
}

type CustomActionsTriggerDefinitionsConnection struct {
	Nodes      []CustomActionsTriggerDefinition
	PageInfo   PageInfo
	TotalCount int
}

type CustomActionsWebhookActionCreateInput struct {
	Name              string                      `json:"name"`
	Description       *graphql.String             `json:"description,omitempty"`
	LiquidTemplate    string                      `json:"liquidTemplate"`
	WebhookURL        string                      `json:"webhookUrl"`
	HTTPMethod        CustomActionsHttpMethodEnum `json:"httpMethod"`
	BasicAuthUserName *graphql.String             `json:"basicAuthUserName,omitempty"`
	BasicAuthPassword *graphql.String             `json:"basicAuthPassword,omitempty"`
	Headers           JSON                        `json:"headers,omitempty"`
}

type CustomActionsWebhookActionUpdateInput struct {
	Id                graphql.ID                  `json:"id"`
	Name              *graphql.String             `json:"name,omitempty"`
	Description       *graphql.String             `json:"description,omitempty"`
	LiquidTemplate    *graphql.String             `json:"liquidTemplate,omitempty"`
	WebhookURL        *graphql.String             `json:"webhookUrl,omitempty"`
	HTTPMethod        CustomActionsHttpMethodEnum `json:"httpMethod,omitempty"`
	BasicAuthUserName *graphql.String             `json:"basicAuthUserName,omitempty"`
	BasicAuthPassword *graphql.String             `json:"basicAuthPassword,omitempty"`
	Headers           JSON                        `json:"headers,omitempty"`
}

type CustomActionsTriggerDefinitionCreateInput struct {
	Name        string          `json:"name"`
	Description *graphql.String `json:"description,omitempty"`
	Owner       graphql.ID      `json:"ownerId"`
	Action      *graphql.ID     `json:"actionId,omitempty"`
	Filter      *graphql.ID     `json:"filterId,omitempty"`
	// This is being explictly left out to reduce the complexity of the implementation
	// action *CustomActionsWebhookActionCreateInput
}

type CustomActionsTriggerDefinitionUpdateInput struct {
	Id          graphql.ID      `json:"id"`
	Name        *graphql.String `json:"name,omitempty"`
	Description *graphql.String `json:"description,omitempty"`
	Owner       *graphql.ID     `json:"ownerId,omitempty"`
	Action      *graphql.ID     `json:"actionId,omitempty"`
	Filter      *graphql.ID     `json:"filterId,omitempty"`
}

func (client *Client) CreateWebhookAction(input CustomActionsWebhookActionCreateInput) (*CustomActionsExternalAction, error) {
	var m struct {
		Payload struct {
			WebhookAction CustomActionsExternalAction
			Errors        []OpsLevelErrors
		} `graphql:"customActionsWebhookActionCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v)
	return &m.Payload.WebhookAction, HandleErrors(err, m.Payload.Errors)
}

// TODO: Not implemented in the API yet
//func (client *Client) GetCustomAction(input IdentifierInput) (*CustomActionsTriggerDefinition, error) {

//}

func (client *Client) ListCustomActions(variables *PayloadVariables) (CustomActionsExternalActionsConnection, error) {
	var q struct {
		Account struct {
			Actions CustomActionsExternalActionsConnection `graphql:"customActionsExternalActions(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = &PayloadVariables{
			"after": "",
			"first": client.pageSize,
		}
	}
	if err := client.Query(&q, nil); err != nil {
		return CustomActionsExternalActionsConnection{}, err
	}
	for q.Account.Actions.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Actions.PageInfo.End
		resp, err := client.ListCustomActions(variables)
		if err != nil {
			return CustomActionsExternalActionsConnection{}, err
		}
		q.Account.Actions.Nodes = append(q.Account.Actions.Nodes, resp.Nodes...)
		q.Account.Actions.PageInfo = resp.PageInfo
	}
	return q.Account.Actions, nil
}

func (client *Client) UpdateWebhookAction(input CustomActionsWebhookActionUpdateInput) (*CustomActionsExternalAction, error) {
	var m struct {
		Payload struct {
			WebhookAction CustomActionsExternalAction
			Errors        []OpsLevelErrors
		} `graphql:"customActionsWebhookActionUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v)
	return &m.Payload.WebhookAction, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) DeleteWebhookAction(input IdentifierInput) error {
	var m struct {
		Payload ResourceDeletePayload `graphql:"customActionsWebhookActionDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v)
	return HandleErrors(err, m.Payload.Errors)
}

func (client *Client) CreateTriggerDefinition(input CustomActionsTriggerDefinitionCreateInput) (*CustomActionsTriggerDefinition, error) {
	var m struct {
		Payload struct {
			TriggerDefinition CustomActionsTriggerDefinition
			Errors            []OpsLevelErrors
		} `graphql:"customActionsTriggerDefinitionCreate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v)
	return &m.Payload.TriggerDefinition, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) GetTriggerDefinition(input IdentifierInput) (*CustomActionsTriggerDefinition, error) {
	var q struct {
		Account struct {
			Definition CustomActionsTriggerDefinition `graphql:"customActionsTriggerDefinition(input: $input)"`
		}
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Query(&q, v)
	if q.Account.Definition.Id == nil {
		err = fmt.Errorf("CustomActionsTriggerDefinition with ID '%s' or Alias '%s' not found", input.Id, input.Alias)
	}
	return &q.Account.Definition, HandleErrors(err, nil)
}

func (client *Client) ListTriggerDefinitions(variables *PayloadVariables) (CustomActionsTriggerDefinitionsConnection, error) {
	var q struct {
		Account struct {
			Definitions CustomActionsTriggerDefinitionsConnection `graphql:"customActionsTriggerDefinitions(after: $after, first: $first)"`
		}
	}
	if variables == nil {
		variables = &PayloadVariables{
			"after": "",
			"first": client.pageSize,
		}
	}
	if err := client.Query(&q, nil); err != nil {
		return CustomActionsTriggerDefinitionsConnection{}, err
	}
	for q.Account.Definitions.PageInfo.HasNextPage {
		(*variables)["after"] = q.Account.Definitions.PageInfo.End
		resp, err := client.ListTriggerDefinitions(variables)
		if err != nil {
			return CustomActionsTriggerDefinitionsConnection{}, err
		}
		q.Account.Definitions.Nodes = append(q.Account.Definitions.Nodes, resp.Nodes...)
		q.Account.Definitions.PageInfo = resp.PageInfo
	}
	return q.Account.Definitions, nil
}

func (client *Client) UpdateTriggerDefinition(input CustomActionsTriggerDefinitionUpdateInput) (*CustomActionsTriggerDefinition, error) {
	var m struct {
		Payload struct {
			TriggerDefinition CustomActionsTriggerDefinition
			Errors            []OpsLevelErrors
		} `graphql:"customActionsTriggerDefinitionUpdate(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v)
	return &m.Payload.TriggerDefinition, HandleErrors(err, m.Payload.Errors)
}

func (client *Client) DeleteTriggerDefinition(input IdentifierInput) error {
	var m struct {
		Payload ResourceDeletePayload `graphql:"customActionsTriggerDefinitionDelete(input: $input)"`
	}
	v := PayloadVariables{
		"input": input,
	}
	err := client.Mutate(&m, v)
	return HandleErrors(err, m.Payload.Errors)
}
