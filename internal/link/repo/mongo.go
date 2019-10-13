package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	db *mongo.Database
}

const shortlinkCollection = "shortlinks"

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{db: db}
}

func (r Mongo) Create(link ShortLink) (*ShortLink, error) {
	_, err := r.db.Collection(shortlinkCollection).InsertOne(context.TODO(), link)
	return &link, err
}

func (r Mongo) ReadOne(slug string) (*ShortLink, error) {
	result := ShortLink{}
	err := r.db.Collection(shortlinkCollection).FindOne(context.TODO(), bson.M{"slug": slug}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &result, err
}

