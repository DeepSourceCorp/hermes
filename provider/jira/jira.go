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
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
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
		log.Errorf("jira: extracting payload: %v", err)
		return nil, err
	}

	if err := payload.Validate(); err != nil {
		log.Errorf("jira: validating payload: %v", err)
		return nil, err
	}

	// Extract and validate the configuration.
	opts := new(Opts)
	if err := opts.Extract(notifier.Config); err != nil {
		log.Errorf("jira: extracting options: %v", err)
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		log.Errorf("jira: validating options: %v", err)
		return nil, err
	}

	request := &CreateIssueRequest{
		Fields: Fields{
			Project: struct {
				Key string "json:\"key\""
			}{Key: opts.ProjectKey},
			IssueType: struct {
				ID string "json:\"id\""
			}{ID: opts.IssueType},
			Summary:     payload.Summary,
			Description: payload.Description,
		},
		CloudID:     opts.CloudID,
		BearerToken: opts.Secret.Token,
	}

	response, err := p.Client.CreateIssue(request)
	if err != nil {
		log.Errorf("jira: creating issue: %v", err)
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

type Rel struct {
	sync.Mutex
	Data map[string]*RelValue
}

type OptValueResponse struct {
	CloudID string `json:"cloud_id"`
	Rel     Rel    `json:"_rel"`
}

func (p *jiraSimple) GetOptValues(_ context.Context, opts *domain.NotifierSecret) (map[string]interface{}, error) {
	sites, err := p.getSites(opts.Token)
	if err != nil {
		log.Errorf("jira: getting options: sites: %v", err)
		return nil, err
	}

	g1, g2 := new(errgroup.Group), new(errgroup.Group)

	relValues := &Rel{Data: make(map[string]*RelValue)}
	for _, site := range sites {
		site := site
		g1.Go(func() error {
			result, err := p.getIssueTypes(opts.Token, site.ID)
			if err != nil {
				log.Errorf("jira: getting options: issueTypes: %v", err)
				return err
			}
			relValues.Mutex.Lock()
			_, ok := relValues.Data[site.ID]
			if !ok {
				relValues.Data[site.ID] = &RelValue{IssueTypes: result}
			} else {
				relValues.Data[site.ID].IssueTypes = result
			}
			relValues.Mutex.Unlock()
			return nil
		})

		g2.Go(func() error {
			result, err := p.getProjectKeys(opts.Token, site.ID)
			if err != nil {
				log.Errorf("jira: getting options: projectKeys: %v", err)
				return err
			}
			relValues.Mutex.Lock()
			_, ok := relValues.Data[site.ID]
			if !ok {
				relValues.Data[site.ID] = &RelValue{ProjectKeys: result}
			} else {
				relValues.Data[site.ID].ProjectKeys = result
			}
			relValues.Mutex.Unlock()
			return nil
		})
	}

	if err := g1.Wait(); err != nil {
		log.Errorf("jira: waiting for g1: %v", err)
		return nil, err
	}

	if err := g2.Wait(); err != nil {
		log.Errorf("jira: waiting for g2: %v", err)
		return nil, err
	}

	return map[string]interface{}{
			"cloud_id": sites,
			"_rel":     map[string]interface{}{"cloud_id": relValues.Data},
		},
		nil
}

func (p *jiraSimple) getSites(token string) (Values, error) {
	request := &AccessibleResourcesRequest{BearerToken: token}
	response, err := p.Client.GetAccessibleResources(request)
	if err != nil {
		log.Errorf("jira: getting accessible resources: %v", err)
		return nil, err
	}

	sites := make([]Value, 0, len(*response))
	for idx := range *response {
		sites = append(sites,
			Value{
				ID:   (*response)[idx].ID,
				Name: (*response)[idx].Name,
			},
		)
	}
	return sites, nil
}

func (p *jiraSimple) getIssueTypes(token, cloudID string) (Values, error) {
	request := &GetIssueTypesRequest{
		BearerToken: token,
		CloudID:     cloudID,
	}

	response, err := p.Client.GetIssueTypes(request)
	if err != nil {
		log.Errorf("jira: getting issue types: %v", err)
		return nil, err
	}

	issueTypes := make([]Value, 0, len(*response))
	for i := range *response {
		issueTypes = append(issueTypes,
			Value{
				ID:   (*response)[i].ID,
				Name: (*response)[i].Name,
			},
		)
	}
	return issueTypes, nil
}

func (p *jiraSimple) getProjectKeys(token, cloudID string) (Values, error) {
	request := &GetProjectsRequest{
		BearerToken: token,
		CloudID:     cloudID,
	}

	response, err := p.Client.GetProjects(request)
	if err != nil {
		log.Errorf("jira: getting projects: %v", err)
		return nil, err
	}
	values := (*response).Values

	projectKeys := make([]Value, 0, len(values))
	for i := range values {
		projectKeys = append(projectKeys,
			Value{
				ID:   values[i].Key,
				Name: values[i].Name,
			},
		)
	}
	return projectKeys, nil
}

// Payload defines the primary content payload for the JIRA provider.
type Payload struct {
	Summary     string                 `json:"summary"`
	Description map[string]interface{} `json:"description"`
}

// Extract unmarshals body to JIRA payload.
func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		log.Errorf("jira: unmarshalling body: %v", err)
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
		log.Errorf("jira: decoding options: %v", err)
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
