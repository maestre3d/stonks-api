package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/maestre3d/stonks-api/internal/domain"
)

// IntegrationEvent is a kind of event based on the CNCF CloudEvent specification used to propagate
// domain events to different modules and systems
type IntegrationEvent struct {
	Id          string `json:"id"`
	Source      string `json:"source"`
	SpecVersion string `json:"specversion"`
	Type        string `json:"type"`

	// Optional region
	DataContentType string      `json:"datacontenttype,omitempty"`
	DataSchema      string      `json:"dataschema,omitempty"`
	Subject         string      `json:"subject,omitempty"`
	Time            time.Time   `json:"time,omitempty"`
	Data            interface{} `json:"data,omitempty"`
}

func (i *IntegrationEvent) FromDomainEvent(e DomainEvent) {
	i.Id = uuid.NewString()
	i.Source = generateSourceContextURI()
	i.SpecVersion = "1.0"
	i.Type = e.Name()
	i.DataContentType = "application/json"
	i.DataSchema = i.Source + "/cloudevents/schemas"
	i.Subject = e.AggregateID()
	i.Time = time.Now().UTC()
	i.Data = e
}

func generateSourceContextURI() string {
	return "https://" + domain.ApplicationName + "." + domain.OrganizationName + "." +
		domain.OrganizationDNSExtension
}
