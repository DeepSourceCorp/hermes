package jira

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type jiraSimple struct {
	Client *Client
}

const ProviderType = domain.ProviderType("jira")

func NewJIRAProvider(httpClient *http.Client) provider.Provider {
	return &jiraSimple{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *jiraSimple) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
	// Extract and validate the payload.
	payload := new(Payload)
	if err := payload.Extract(body); err != nil {
		return nil, err
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Extract and validate the configuration.
	opts := new(Opts)
	if err := opts.Extract(notifier.Config); err != nil {
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	request := &CreateIssueRequest{
		Fields: Fields{
			Project:     Project{Key: opts.ProjectKey},
			IssueType:   IssueType{Name: opts.IssueType},
			Summary:     payload.Summary,
			Description: payload.Description,
		},
		CloudID:     opts.CloudID,
		BearerToken: opts.Secret.Token,
	}

	response, err := p.Client.CreateIssue(request)
	if err != nil {
		return nil, err
	}

	return &domain.Message{
		ID:               ksuid.New().String(),
		Ok:               true,
		Payload:          payload,
		ProviderResponse: response,
	}, nil
}

func (p *jiraSimple) GetOptValues(context.Context, *domain.NotifierSecret) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// Payload defines the primary content payload for the JIRA provider.
type Payload struct {
	Summary     string                 `json:"summary"`
	Description map[string]interface{} `json:"description"`
}

// Extract unmarshals body to JIRA payload.
func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		return errFailedBodyValidation(err.Error())
	}
	return nil
}

// Validate() validates the payload ensuring all mandatory properties are set.
// Description should ideally be validated agaings JDF (Jira Document Format)
func (p *Payload) Validate() domain.IError {
	if p.Summary == "" {
		return errFailedBodyValidation(
			"generated payload does not contaion mandatory param summary",
		)
	}
	if p.Description == nil {
		return errFailedBodyValidation(
			"generated payload does not contain mandatory param description",
		)
	}
	return nil
}

// Opts defines JIRA specific options.
type Opts struct {
	Secret     *domain.NotifierSecret
	ProjectKey string `mapstructure:"project_key"`
	IssueType  string `mapstructure:"issue_type"`
	CloudID    string `mapstructure:"cloud_id"`
}

func (o *Opts) Extract(c *domain.NotifierConfiguration) domain.IError {
	if c == nil {
		return errFailedOptsValidation("notifier config empty")
	}
	if err := mapstructure.Decode(c.Opts, o); err != nil {
		return errFailedOptsValidation("failed to decode configuration")
	}
	o.Secret = c.Secret
	return nil
}

// ValidateAndExtractOpts validates the notifier configuration and returns
// JIRA specific opts.
func (o *Opts) Validate() domain.IError {
	if o == nil {
		return errFailedOptsValidation("options empty")
	}
	if o.IssueType == "" || o.ProjectKey == "" {
		return errFailedOptsValidation("issue_type or project_key is empty")
	}

	if o.Secret == nil || o.Secret.Token == "" {
		return errFailedOptsValidation("secret not defined in configuration")
	}
	return nil
}
