package data

import (
	"context"
	"fmt"
	"os"
	"time"
	"user-service/util"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var L *log.Logger
var collection *mongo.Collection

func init() {
	_ = godotenv.Load("connection.env")
	_, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	clientOptions := options.Client().
		ApplyURI(os.Getenv("CONNECTION_DB"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
		L.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		L.Fatal(err)
	}
	c := client.Database("ott_platform").Collection("userDetails")
	collection = c
}

func init() {
	L = &log.Logger{}
	L.SetFormatter(&log.TextFormatter{
		TimestampFormat: "02-01-2006 15:04:05",
		FullTimestamp:   true,
	})
	L.SetOutput(os.Stdout)
}
func GetMongoDB() *mongo.Collection {
	return collection
}
