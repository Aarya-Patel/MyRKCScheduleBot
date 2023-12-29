package storage

import (
	"context"

	"cloud.google.com/go/datastore"
)

var dbClient *datastore.Client

func GetDataStoreClient(projectId string) (*datastore.Client, error) {
	if dbClient == nil {
		ctx := context.Background()
		client, err := datastore.NewClient(ctx, projectId)

		if err != nil {
			return nil, err
		}

		dbClient = client
	}

	return dbClient, nil
}
