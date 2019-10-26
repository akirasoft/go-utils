package events

import (
	"encoding/json"

	"github.com/keptn/go-utils/pkg/models"
)

// ServiceCreateEventType is a CloudEvent type for creating a new service
const ServiceCreateEventType = "sh.keptn.event.service.create"

// InternalServiceCreateEventType is a CloudEvent type for creating a new service
const InternalServiceCreateEventType = "sh.keptn.internal.event.service.create"

// ProjectCreateEventType is a CloudEvent type for creating a new project
const ProjectCreateEventType = "sh.keptn.event.project.create"

// ProjectDeleteEventType is a CloudEvent type for deleting a project
const ProjectDeleteEventType = "sh.keptn.event.project.delete"

// InternalProjectCreateEventType is a CloudEvent type for creating a new project
const InternalProjectCreateEventType = "sh.keptn.internal.event.project.create"

// InternalProjectDeleteEventType is a CloudEvent type for deleting a project
const InternalProjectDeleteEventType = "sh.keptn.internal.event.project.delete"

// ConfigurationChangeEventType is a CloudEvent type for changing the configuration
const ConfigurationChangeEventType = "sh.keptn.event.configuration.change"

// ProblemOpenEventType is a CloudEvent type to inform about an open problem
const ProblemOpenEventType = "sh.keptn.event.problem.open"

// ConfigureMonitoringEventType is a CloudEvent for configuring monitoring
const ConfigureMonitoringEventType = "sh.keptn.event.monitoring.configure"

// TestsFinishedEventType is a CloudEvent for indicating that tests have finished
const TestsFinishedEventType = "sh.keptn.event.tests.finished"

// ProjectCreateEventData represents the data for creating a new project
type ProjectCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Shipyard is a base64 encoded string of the shipyard file
	Shipyard string `json:"shipyard"`
	// GitUser is the name of a git user of an upstream repository
	GitUser string `json:"gitUser,omitempty"`
	// GitToken is the authentication token for the git user
	GitToken string `json:"gitToken,omitempty"`
	// GitRemoteURL is the remote url of a repository
	GitRemoteURL string `json:"gitRemoteURL,omitempty"`
}

// ProjectDeleteEventData represents the data for deleting a new project
type ProjectDeleteEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
}

// ServiceCreateEventData represents the data for creating a new service
type ServiceCreateEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// HelmChart are the data of a Helm chart packed as tgz and base64 encoded
	HelmChart string `json:"helmChart"`
	// DeploymentStrategies contains the deployment strategy for the stages
	DeploymentStrategies map[string]DeploymentStrategy `json:"deploymentStrategies"`
}

// ConfigurationChangeEventData represents the data for changing the service configuration
type ConfigurationChangeEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// ValuesCanary contains new Helm values for canary
	ValuesCanary map[string]interface{} `json:"valuesCanary,omitempty"`
	// Canary contains a new configuration for canary releases
	Canary *Canary `json:"canary,omitempty"`
	// DeploymentChanges contains changes of the primary deployment
	DeploymentChanges []PropertyChange `json:"deploymentChanges,omitempty"`
}

// TestsFinishedEventData represents the data for a test finished event
type TestsFinishedEventData struct {
	// Project is the name of the project
	Project string `json:"project"`
	// Service is the name of the new service
	Service string `json:"service"`
	// Stage is the name of the stage
	Stage string `json:"stage"`
	// TestStrategy is the testing strategy
	TestStrategy string `json:"teststrategy"`
}

// PropertyChange describes the property to be changed
type PropertyChange struct {
	PropertyPath string      `json:"propertyPath"`
	Value        interface{} `json:"value"`
}

// Canary describes the new configuration in a canary release
type Canary struct {
	// Value represents the traffic percentage on the canary
	Value int32 `json:"value,omitempty"`
	// Action represents the action of the canary
	Action CanaryAction `json:"action"`
}

// ProblemEventData represents the data for describing a problem
type ProblemEventData struct {
	State          string `json:"state"`
	ProblemID      string `json:"problemID"`
	ProblemTitle   string `json:"problemtitle"`
	ProblemDetails string `json:"problemdetails"`
	ImpactedEntity string `json:"impactedEntity"`
}

// ConfigureMonitoringEventData represents the data necessary to configure monitoring for a service
type ConfigureMonitoringEventData struct {
	Type              string                    `json:"type"`
	Project           string                    `json:"project"`
	Service           string                    `json:"service"`
	ServiceIndicators *models.ServiceIndicators `json:"serviceIndicators"`
	ServiceObjectives *models.ServiceObjectives `json:"serviceObjectives"`
	Remediation       *models.Remediations      `json:"remediation"`
}

// EvaluationDoneEventData Keptn event payload for completed Pitometer evaluation Note: many elements are not strongly typed
type EvaluationDoneEventData struct {
	Deploymentstrategy string `json:"deploymentstrategy"`
	Evaluationdetails  []struct {
		Key string `json:"Key"`
		// we need to parse this later as it could be a string or array
		Value json.RawMessage `json:"Value"`
	} `json:"evaluationdetails"`
	Evaluationpassed bool   `json:"evaluationpassed"`
	Project          string `json:"project"`
	Service          string `json:"service"`
	Stage            string `json:"stage"`
	Teststrategy     string `json:"teststrategy"`
}

// DeploymentFinishedEventData Keptn event payload for completed deployment
type DeploymentFinishedEventData struct {
	Deploymentstrategy string `json:"deploymentstrategy"`
	Project            string `json:"project"`
	Service            string `json:"service"`
	Stage              string `json:"stage"`
	Teststrategy       string `json:"teststrategy"`
}
