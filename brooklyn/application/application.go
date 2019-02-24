package application

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"text/template"
	"time"

	api_application "github.com/apache/brooklyn-client/cli/api/application"
	"github.com/apache/brooklyn-client/cli/api/entities"
	"github.com/jittakal/terraform-provider-brooklyn/brooklyn/config"
	"github.com/jittakal/terraform-provider-brooklyn/brooklyn/utils"
)

var (
	errorNotStarting = errors.New("Brooklyn application state should be Starting: Maximum number of retries (10) exceeded")
)

// Configuration configuration
type Configuration struct {
	Key   string
	Value string
}

// Application structure
type Application struct {
	Client config.Client

	ID             string
	Name           string
	Location       string
	Type           string
	Configurations []Configuration
}

// Template for Application Yaml
const applicationYamlTmpl = `name: {{.Name}}
location: {{.Location}}
services:
  - type: {{.Type}} {{ if gt (len .Configurations) 0 }}
    brooklyn.config:{{ range $configuration := .Configurations }}
      {{ $configuration.Key }}: {{ $configuration.Value }}{{ end }}{{ end }}`

func (a *Application) yamlAsBytes() ([]byte, error) {
	log.Printf("[INFO] Calling .applicationYamlBytes()")
	// Create a new template and parse the application into it.
	var t *template.Template
	t = template.Must(template.New("provider").Parse(applicationYamlTmpl))
	var appYml bytes.Buffer
	err := t.Execute(&appYml, a)
	log.Printf("[DEBUG] %s", appYml.String())
	var app []byte
	if err != nil {
		return app, err
	}
	return appYml.Bytes(), nil
}

// Create validate and deploy the application
func (a *Application) Create() error {
	appYaml, err := a.yamlAsBytes()
	if err != nil {
		return err
	}

	taskSummary, err := api_application.CreateFromBytes(a.Client, appYaml)
	if err != nil {
		return err
	}
	a.ID = taskSummary.EntityId
	return nil
}

// Expunge delete application
func (a *Application) Expunge() error {
	appID := a.ID
	if appID != "" {
		resp, err := entities.Expunge(a.Client, appID, appID)

		if err != nil {
			return err
		}
		log.Printf("[INFO] %s", resp)
	}
	return nil
}

// Delete invokes expunge request for application.
func (a *Application) Delete() error {
	if a.ID == "" {
		return nil
	}

	url := fmt.Sprintf("/v1/applications/%s/entities/%s/expunge?release=true", a.ID, a.ID)
	//var response models.TaskSummary
	_, err := a.Client.SendEmptyPostRequest(url)
	// if err != nil {
	// 	return err
	// }
	//err = json.Unmarshal(body, &response)
	return err
}

// Rename rename application name
func (a *Application) Rename() error {
	appID := a.ID
	if appID != "" && a.Name != "" {
		resp, err := entities.Rename(a.Client, appID, appID, a.Name)

		if err != nil {
			return err
		}
		log.Printf("[INFO] %s", resp)
	}
	return nil
}

// WaitForApplicationRunning wait for application running
func (a *Application) WaitForApplicationRunning() error {
	log.Println("[INFO] Calling .waitForApplicationRunning()")
	return utils.WaitForSpecific(a.isApplicationRunning, 120, 5*time.Second)
}

// WaitForApplicationStarting wait for application starting
func (a *Application) WaitForApplicationStarting() error {
	log.Println("[INFO] Calling .waitForApplicationStarting()")
	if err := utils.WaitForSpecific(a.isApplicationStarting, 20, 3*time.Second); err != nil {
		return errorNotStarting
	}
	return nil
}

func (a *Application) isApplicationRunning() bool {
	log.Println("[INFO] Calling .isApplicationRunning()")
	st, err := a.GetApplicationState()
	if err != nil {
		log.Println("[ERROR] %s", err)
	}
	if st == Running {
		return true
	}
	return false
}

func (a *Application) isApplicationStarting() bool {
	log.Println("[INFO] Calling .isApplicationStarting()")
	st, err := a.GetApplicationState()
	if err != nil {
		log.Println("[ERROR] %s", err)
	}
	if st == Starting {
		return true
	}
	return false
}

// GetApplicationState returns the state that the host is in (running, stopped, etc)
func (a *Application) GetApplicationState() (State, error) {
	log.Println("[INFO] Calling .GetApplicationState()")

	if a.ID == "" {
		log.Println("[INFO] Application id is nil.")
		return Stopped, nil
	}

	applicationSummary, err := api_application.Application(a.Client, a.ID)
	if err != nil {
		log.Println("[INFO] Application does not exists.")
		return Stopped, nil
	}

	log.Println("[INFO] Application status is %s", applicationSummary.Status)
	switch applicationSummary.Status {
	case "RUNNING":
		return Running, nil
	case "STARTING":
		return Starting, nil
	case "STOPPING":
		return Stopping, nil
	case "ERROR":
		return Error, nil
	case "STOPPED":
		return Stopped, nil
	default:
		return None, nil
	}
}
