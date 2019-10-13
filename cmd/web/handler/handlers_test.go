package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "github.com/lib/pq"

	"github.com/jasondeutsch/shorty/internal/link"
	"github.com/jasondeutsch/shorty/internal/stats"
)

var shortLinkHandler *Handler

const cloudflare = "{\"destination\":\"https://www.cloudflare.com\"}"

func TestMain(m *testing.M) {
	// TODO this whole setup needs improvement
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURL := os.Getenv("MONGO_TEST_URL")
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient.Connect(ctx)
	mongoDB := mongoClient.Database("shorty-test")
	shortlinkCollection := mongoDB.Collection("shortlinks")

	connectionInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DB"))

	postgresDB, err := sql.Open("postgres",connectionInfo )
	if err != nil {
		log.Fatal(err)
	}

	if err := postgresDB.Ping(); err != nil {
		log.Fatal(err)
	}

	t := true // wtf mongo?
	mod := mongo.IndexModel{
		Keys: bson.M{
			"slug": 1,
		}, Options: &options.IndexOptions{Unique: &t},
	}
	_, err = shortlinkCollection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		fmt.Println(err)
	}

	statsRepo := stats.NewPostgres(postgresDB)
	statsService := stats.NewStatsService(statsRepo)

	shortLinkRepo := link.NewMongo(mongoDB)
	shortLinkService := link.NewShortLinkService(shortLinkRepo)
	shortLinkHandler = NewHandler(shortLinkService, statsService)

	runTests := m.Run()
	os.Exit(runTests)
}


func TestRedirectFromShortLink(t *testing.T) {
	// first create the shirt link
	rr := testCreateShortlink(t, cloudflare)

	type testResponse struct {
		Slug string `json:"slug"`
	}
	var  resp testResponse

	json.NewDecoder(rr.Body).Decode(&resp)


	// now test the redirect
	req, err := http.NewRequest("get", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// this is a pain
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("slug", resp.Slug)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr = httptest.NewRecorder()

	shortLinkHandler.RedirectFromShortLink(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}
}

func TestCreateLink(t *testing.T) {
	rr := testCreateShortlink(t, cloudflare)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func testCreateShortlink(t *testing.T, body string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("post", "/links", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	shortLinkHandler.CreateShortLink(rr, req)
	return rr
}

func TestCreateLinkInvalidDomain(t *testing.T) {
	rr := testCreateShortlink(t, "www.not a valid- domain")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
