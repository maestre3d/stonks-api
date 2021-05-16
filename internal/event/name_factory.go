package event

import (
	"strings"

	"github.com/maestre3d/stonks-api/internal/domain"
)

func newBaseEventName(domainContext string) string {
	return domain.OrganizationDNSExtension + "." + domain.OrganizationName + "." + domain.ApplicationMajorVersion +
		"." + domain.ApplicationName + "." + domainContext + "."
}

func newWatcherContextEventName(aggregateName string) string {
	return newBaseEventName(domain.WatcherDomainContextName) + aggregateName + "."
}

func newAssetEventName(action string) string {
	return newWatcherContextEventName(domain.AssetAggregateName) + strings.ToLower(action) // avoid sprintf since it relays on reflection
}
