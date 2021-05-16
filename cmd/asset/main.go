package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/maestre3d/stonks-api/internal/command"
	"github.com/maestre3d/stonks-api/internal/eventbus"
	"github.com/maestre3d/stonks-api/internal/persistence"
	"github.com/maestre3d/stonks-api/internal/usecase"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	dbClient := dynamodb.NewFromConfig(cfg)
	snsClient := sns.NewFromConfig(cfg)
	eventBus := eventbus.NewAWBus(snsClient, nil, "us-east-1", "AWS_ACCOUNT_ID")

	assetRepo := persistence.NewAssetDynamoRepository(dbClient)
	assetCase := usecase.NewRegisterAsset(assetRepo, eventBus)

	err = command.RegisterAssetHandler(ctx, assetCase, command.RegisterAsset{
		TickerSymbol: "VOO",
		SpreadPrice:  330,
	})
	if err != nil {
		panic(err)
	}
}
