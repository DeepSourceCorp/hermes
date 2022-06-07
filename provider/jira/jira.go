package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

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
	var payload = new(Payload)
	if err := payload.Extract(body); err != nil {
		return nil, err
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Extract and validate the configuration.
	var opts = new(Opts)
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

type Value struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Values []Value

type RelValue struct {
	ProjectKeys Values `json:"project_key"`
	IssueTypes  Values `json:"issue_type"`
}

type Rel map[string]RelValue

type OptValueResponse struct {
	CloudID string `json:"cloud_id"`
	Rel     Rel    `json:"_rel"`
}

func (p *jiraSimple) GetOptValues(_ context.Context, opts *domain.NotifierSecret) (map[string]interface{}, error) {
	sites, err := p.getSites(opts.Token)
	if err != nil {
		return nil, err
	}

	relValues := Rel{}

	itChan := make(chan Values, 10)
	pkChan := make(chan Values, 10)

	for idx := range sites {
		go p.getIssueTypes(opts.Token, sites[idx].ID, itChan)
		if err != nil {
			return nil, err
		}
		go p.getProjectKeys(opts.Token, sites[idx].ID, pkChan)
		if err != nil {
			return nil, err
		}

	}

	wg := new(sync.WaitGroup)
	for idx := range sites {
		projectKeys := Values{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			projectKeys = <-pkChan
		}()
		issueTypes := <-itChan
		wg.Wait()
		relValues[sites[idx].ID] = RelValue{
			ProjectKeys: projectKeys,
			IssueTypes:  issueTypes,
		}
	}

	return map[string]interface{}{
		"cloud_id": sites,
		"_rel": map[string]interface{}{
			"cloud_id": relValues,
		},
	}, nil
}

func (p *jiraSimple) getSites(token string) (Values, error) {
	sites := Values{}
	request := &AccessibleResourcesRequest{
		BearerToken: token,
	}
	response, err := p.Client.GetAccessibleResources(request)
	if err != nil {
		return nil, err
	}
	for idx := range *response {
		sites = append(sites, Value{
			ID:   (*response)[idx].ID,
			Name: (*response)[idx].Name,
		})
	}
	return sites, nil
}

func (p *jiraSimple) getIssueTypes(token string, cloudID string, itChan chan Values) error {
	issueTypes := Values{}

	request := &GetIssueTypesRequest{
		BearerToken: token,
		CloudID:     cloudID,
	}
	response, err := p.Client.GetIssueTypes(request)
	if err != nil {
		return err
	}
	for i, _ := range *response {
		issueTypes = append(issueTypes, Value{
			ID:   (*response)[i].ID,
			Name: (*response)[i].Name,
		})
	}
	itChan <- issueTypes
	return nil
}

func (p *jiraSimple) getProjectKeys(token string, cloudID string, pkChan chan Values) error {
	projectKeys := Values{}
	request := &GetProjectsRequest{
		BearerToken: token,
		CloudID:     cloudID,
	}
	response, err := p.Client.GetProjects(request)
	if err != nil {
		return err
	}

	values := (*response).Values
	for i, _ := range values {
		projectKeys = append(projectKeys, Value{
			ID:   values[i].Key,
			Name: values[i].Name,
		})
	}
	pkChan <- projectKeys
	return nil
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
		return errFailedOptsValidation("notifier config emtpy")
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
		return errFailedOptsValidation("issue_type or project_key is emtpy")
	}

	if o.Secret == nil || o.Secret.Token == "" {
		return errFailedOptsValidation("secret not defined in configuration")
	}
	return nil
}
