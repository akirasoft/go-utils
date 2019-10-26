package utils

import (
	"context"
	"errors"
	"log"

	keptnevents "github.com/akirasoft/go-utils/pkg/events"
	cloudevents "github.com/cloudevents/sdk-go"
)

// RcvConfig stores configuration elements for cloudevents listener
type RcvConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

// KeptnEventHandler parses Keptn events and returns the Keptn event payload
func KeptnEventHandler(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
	eventID := event.Context.GetID()
	switch event.Type() {
	case "sh.keptn.event.configuration.change":
		//log.Println("sh.keptn.events.configuration-changed")
		// should be sh.keptn.events.configuration.change
		data := &keptnevents.ConfigurationChangeEventData{}
		if err := event.DataAs(data); err != nil {
			return err
		}
		receiver(data, shkeptncontext, eventID)
	case "sh.keptn.events.deployment-finished":
		//log.Println("sh.keptn.events.deployment-finished")
		data := &keptnevents.DeploymentFinishedEventData{}
		if err := event.DataAs(data); err != nil {
			return err
		}
		receiver(data, shkeptncontext, eventID)
	case "sh.keptn.events.evaluation-done":
		//log.Println("sh.keptn.events.evaluation-done")
		data := &keptnevents.EvaluationDoneEventData{}
		if err := event.DataAs(data); err != nil {
			return err
		}
		receiver(data, shkeptncontext, eventID)
	case "sh.keptn.events.tests-finished":
		//log.Println("sh.keptn.events.tests-finished")
		data := &keptnevents.TestsFinishedEventData{}
		if err := event.DataAs(data); err != nil {
			return err
		}
		receiver(data, shkeptncontext, eventID)
	case "sh.keptn.events.problem":
		//log.Println("sh.keptn.events.problem")
		data := &keptnevents.ProblemEventData{}
		if err := event.DataAs(data); err != nil {
			return err
		}
		receiver(data, shkeptncontext, eventID)
	default:
		const errorMsg = "Received unexpected keptn event"
		return errors.New(errorMsg)
	}

	return nil
}

var receiver func(interface{}, string, string) error

// KeptnReceiver listens for Keptn events on the path and port defined via Rcv
// returns Keptn event data, shkeptncontext and cloudevent ID
func KeptnReceiver(Rcv RcvConfig, keptnEventCallback func(interface{}, string, string) error) error {

	ctx := context.Background()

	receiver = keptnEventCallback

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithPort(Rcv.Port),
		cloudevents.WithPath(Rcv.Path),
	)
	if err != nil {
		log.Printf("failed to create transport, %v", err)
		return err
	}
	c, err := cloudevents.NewClient(t)
	if err != nil {
		log.Printf("failed to create client, %v", err)
		return err
	}

	return c.StartReceiver(ctx, KeptnEventHandler)
}
