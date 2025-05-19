package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepository(client *mongo.Client, dbName, collectionName string) *ProductRepo {
	coll := client.Database(dbName).Collection(collectionName)
	return &ProductRepo{collection: coll}
}

func (r *ProductRepo) GetProduct(ctx context.Context, productId primitive.ObjectID) (Product, error) {
	var product Product
	filter := bson.M{"_id": productId}
	err := r.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return product, err
	}

	return product, nil
}

type ProductGetter interface {
	GetProduct(ctx context.Context, productId primitive.ObjectID) (Product, error)
}
