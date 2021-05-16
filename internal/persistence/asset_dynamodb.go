package persistence

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/maestre3d/stonks-api/internal/aggregate"
	"github.com/maestre3d/stonks-api/internal/repository"
	"github.com/maestre3d/stonks-api/internal/valueobject"
)

const (
	assetDynamoTable        = "stonks-tickers"
	assetDynamoPartitionKey = "ticker_symbol"
)

type assetDynamo struct {
	TickerSymbol string  `dynamodbav:"ticker_symbol"`
	SpreadPrice  float64 `dynamodbav:"spread_price"`
}

func (m *assetDynamo) FromAggregate(asset aggregate.Asset) {
	m.TickerSymbol = asset.TickerSymbol.Value()
	m.SpreadPrice = asset.SpreadPrice.Value()
}

func (m assetDynamo) ToAggregate() *aggregate.Asset {
	return &aggregate.Asset{
		TickerSymbol: valueobject.TickerSymbol(m.TickerSymbol),
		SpreadPrice:  valueobject.Spread(m.SpreadPrice),
	}
}

type AssetDynamoRepository struct {
	client *dynamodb.Client
	mu     sync.RWMutex
}

var _ repository.Asset = &AssetDynamoRepository{}

func NewAssetDynamoRepository(c *dynamodb.Client) *AssetDynamoRepository {
	return &AssetDynamoRepository{
		client: c,
		mu:     sync.RWMutex{},
	}
}

func (a *AssetDynamoRepository) Save(ctx context.Context, asset aggregate.Asset) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	model := new(assetDynamo)
	model.FromAggregate(asset)

	assetDynamo, err := attributevalue.MarshalMap(model)
	if err != nil {
		return err
	}

	_, err = a.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      assetDynamo,
		TableName: aws.String(assetDynamoTable),
	})
	return err
}

func (a *AssetDynamoRepository) Update(ctx context.Context, asset aggregate.Asset) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return nil
}

func (a *AssetDynamoRepository) Remove(ctx context.Context, asset aggregate.Asset) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	model := new(assetDynamo)
	model.FromAggregate(asset)

	assetDynamo, err := attributevalue.MarshalMap(model)
	if err != nil {
		return err
	}

	_, err = a.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key: map[string]types.AttributeValue{
			assetDynamoPartitionKey: assetDynamo[assetDynamoPartitionKey],
		},
		TableName: aws.String(assetDynamoTable),
	})
	return err
}

func (a *AssetDynamoRepository) Find(ctx context.Context, ticker valueobject.TickerSymbol) (*aggregate.Asset, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	o, err := a.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			assetDynamoPartitionKey: &types.AttributeValueMemberS{
				Value: ticker.Value(),
			},
		},
		TableName: aws.String(assetDynamoTable),
	})
	if err != nil {
		return nil, err
	} else if o.Item == nil {
		return nil, nil
	}

	model := new(assetDynamo)
	if err = attributevalue.UnmarshalMap(o.Item, model); err != nil {
		return nil, err
	}
	return model.ToAggregate(), nil
}

func (a *AssetDynamoRepository) Search(ctx context.Context) ([]*aggregate.Asset, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return nil, nil
}
