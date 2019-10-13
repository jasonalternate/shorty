package main

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jasondeutsch/shorty/cmd/web/handler"
	"github.com/jasondeutsch/shorty/internal/link"
)

func main() {
	log.Printf("starting...")

	// config

	const port = ":8080"

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// persistence

	mongoDB, err := initMongo()
	if err != nil {
		log.Fatal(err)
	}

	// services
	shortlinkRepo := link.NewMongo(mongoDB)
	shortlinkService := link.NewShortLinkService(shortlinkRepo)
	shortLinkHandlers := handler.NewHandler(shortlinkService)

	// router

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Mount("/", shortLinkHandlers)
		r.NotFound(handler.Cow) // TODO
	})

	// serve

	log.Printf("running on port: %s", port)

	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func initMongo() (*mongo.Database, error) {
	mongoURL := os.Getenv("MONGO_TEST_URL")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, err
	}

	mongoClient.Connect(ctx)
	mongoDB := mongoClient.Database("shorty-test")
	shortlinkCollection := mongoDB.Collection("shortlinks")

	t := true // wtf mongo?
	mod := mongo.IndexModel{
		Keys: bson.M{
			"slug": 1,
		}, Options: &options.IndexOptions{Unique: &t},
	}

	_, err = shortlinkCollection.Indexes().CreateOne(ctx, mod)

	return mongoDB, err
}
