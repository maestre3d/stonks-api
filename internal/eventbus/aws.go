package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/hashicorp/go-multierror"
	"github.com/maestre3d/stonks-api/internal/event"
)

// ASWSBus is the Amazon Web Services (AWS) serverless implementation using a topic-queue-chaining pattern with AWS SNS and AWS SQS
type AWSBus struct {
	snsClient *sns.Client
	sqsClient *sqs.Client
	region    string
	accountID string
	mu        sync.Mutex
}

var _ event.Bus = &AWSBus{}

func NewAWBus(snsClient *sns.Client, sqsClient *sqs.Client, region, accoundID string) *AWSBus {
	return &AWSBus{
		snsClient: snsClient,
		sqsClient: sqsClient,
		region:    region,
		accountID: accoundID,
		mu:        sync.Mutex{},
	}
}

func (b *AWSBus) Publish(ctx context.Context, events ...event.DomainEvent) error {
	errs := new(multierror.Error)
	for _, e := range events {
		integrationEv := new(event.IntegrationEvent)
		integrationEv.FromDomainEvent(e)
		if err := b.publish(ctx, *integrationEv); err != nil {
			errs = multierror.Append(err, errs)
		}
	}
	return errs.ErrorOrNil()
}

func (b *AWSBus) publish(ctx context.Context, e event.IntegrationEvent) error {
	dataJSON, err := json.Marshal(e.Data)
	if err != nil {
		return err
	}

	topic := newTopicArn(e, b.region, b.accountID)
	log.Print(topic)
	_, err = b.snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn:          aws.String(topic),
		Message:           aws.String(string(dataJSON)),
		MessageAttributes: newSNSAttributesFromEvent(e),
	})
	return err
}

func newTopicArn(e event.IntegrationEvent, region, awsAccountID string) string {
	return "arn:aws:sns:" + strings.ToLower(region) + ":" + awsAccountID +
		":" + newTopicNameFromEvent(e)
}

func newTopicNameFromEvent(e event.IntegrationEvent) string {
	log.Print(e.Type)
	topicSplit := strings.Split(e.Type, ".")
	if len(topicSplit) < 7 {
		// does not uses application event name nomeclature
		// (dns_ext.org_name.major_version.app_name.context_name.aggregate_name.action)
		return ""
	}

	return strings.Join(topicSplit[3:], "-")
}

func newSNSAttributesFromEvent(e event.IntegrationEvent) map[string]types.MessageAttributeValue {
	dataType := aws.String("String")
	return map[string]types.MessageAttributeValue{
		"ce_id": {
			StringValue: aws.String(e.Id),
			DataType:    dataType,
		},
		"ce_source": {
			StringValue: aws.String(e.Source),
			DataType:    dataType,
		},
		"ce_specversion": {
			StringValue: aws.String(e.SpecVersion),
			DataType:    dataType,
		},
		"ce_type": {
			StringValue: aws.String(e.Type),
			DataType:    dataType,
		},
		"datacontenttype": {
			StringValue: aws.String(e.DataContentType),
			DataType:    dataType,
		},
		"ce_dataschema": {
			StringValue: aws.String(e.DataSchema),
			DataType:    dataType,
		},
		"ce_subject": {
			StringValue: aws.String(e.Subject),
			DataType:    dataType,
		},
		"ce_time": {
			StringValue: aws.String(fmt.Sprint(e.Time.Unix())),
			DataType:    aws.String("Number"),
		},
	}
}
