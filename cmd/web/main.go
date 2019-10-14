package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jasondeutsch/shorty/cmd/web/handler"
	"github.com/jasondeutsch/shorty/internal/link"
	"github.com/jasondeutsch/shorty/internal/stats"
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

	postgresDB, err := initPostgres()
	if err != nil {
		log.Fatal(err)
	}

	// services
	statsRepo := stats.NewPostgres(postgresDB)
	statsService := stats.NewStatsService(statsRepo)

	shortlinkRepo := link.NewMongo(mongoDB)
	shortlinkService := link.NewShortLinkService(shortlinkRepo)
	shortLinkHandlers := handler.NewHandler(shortlinkService, statsService)

	// router
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Mount("/", shortLinkHandlers)
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

func initPostgres() (*sql.DB, error) {
	connectionInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("postgres",connectionInfo )
	return db, err
}

