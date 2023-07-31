package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error", err)
	}
	log.Println("loaded")
}

var DB *mongo.Client
var Collections struct {
	Accounts   *mongo.Collection
	Loans      *mongo.Collection
	Repayments *mongo.Collection
	Users      *mongo.Collection
}

func connectDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MongoUri")
	log.Print(mongoURI)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to create a new client: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	return client, nil
}

func Init() error {
	log.Print("Initializing database connection")

	var err error
	DB, err = connectDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	databaseName := os.Getenv("DatabaseName")
	log.Println("databaseName", databaseName)
	Collections.Accounts = DB.Database(databaseName).Collection(os.Getenv("AccountCollection"))
	Collections.Loans = DB.Database(databaseName).Collection(os.Getenv("LoanCollection"))
	Collections.Repayments = DB.Database(databaseName).Collection(os.Getenv("RepayCollection"))
	Collections.Users = DB.Database(databaseName).Collection(os.Getenv("UserCollection"))
	return nil
}

func InitializeWithDBSession() (mongo.Session, error) {
	var err error
	DB, err = connectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	session, err := DB.StartSession()
	if err != nil {
		fmt.Println("Error starting session:", err)
		return nil, err
	}

	databaseName := os.Getenv("DatabaseName")
	log.Println("databaseName", databaseName)
	Collections.Accounts = DB.Database(databaseName).Collection(os.Getenv("AccountCollection"))
	Collections.Loans = DB.Database(databaseName).Collection(os.Getenv("LoanCollection"))
	Collections.Repayments = DB.Database(databaseName).Collection(os.Getenv("RepayCollection"))
	Collections.Users = DB.Database(databaseName).Collection(os.Getenv("UserCollection"))
	return session, nil
}
