package event

type DomainEvent interface {
	AggregateID() string
	Name() string
}
