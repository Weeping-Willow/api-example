package testt

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
)

var GetConnection func(dbHost string) (*mongo.Client, error)

func Database(t *testing.T) *mongo.Database {
	t.Helper()
	client, err := GetConnection(Conf.DatabaseUrl)
	assert.NoError(t, err)
	dbName := randStringRunes(12)
	db := client.Database(dbName)

	t.Cleanup(func() {
		err := db.Drop(context.TODO())
		assert.NoError(t, err)
		client.Disconnect(context.TODO())
	})
	return db
}

func DBWithModelMocks(t *testing.T, dataFunc func() interface{}, collectionName string) *mongo.Database {
	t.Helper()
	db := Database(t)
	insertModelsIntoDB(t, db, dataFunc, collectionName)
	return db
}

func insertModelsIntoDB(t *testing.T, db *mongo.Database, dataFunc func() interface{}, collectionName string) {
	t.Helper()
	data := dataFunc()
	dataToInsert := make([]interface{}, 0)
	switch kind := reflect.TypeOf(data).Kind(); kind {
	case reflect.Slice:
		s := reflect.ValueOf(data)
		for i := 0; i < s.Len(); i++ {
			dataToInsert = append(dataToInsert, s.Index(i).Interface())
		}
	case reflect.Struct:
		dataToInsert = append(dataToInsert, reflect.ValueOf(data).Interface())
	default:
		require.FailNowf(t, "invalid type given to model mocks : %s", kind.String())
	}

	collection := db.Collection(collectionName)
	_, err := collection.InsertMany(context.TODO(), dataToInsert)
	require.Nil(t, err, "failed to insert mock data into database")
}
