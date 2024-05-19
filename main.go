package main

import (
	"context"
	// "encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var dbName, uri string = "dairyDB", "mongodb://localhost:27017/"
	if err := godotenv.Load(); err != nil {
		log.Println("Set your 'MONGODB_URI' environment variable. " + "No .env file found\nUsing the default 'mongodb://localhost:27017'")
		uri = uri + dbName
	} else {
		dbName = os.Getenv("DB_NAME")
		uri = os.Getenv("MONGODB_URI") + dbName

	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to MongoDB! ", dbName)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	// validator := bson.M{
	// 	"bsonType": "object",
	// 	"required": []string{"phone", "name"},
	// 	"properties": bson.M{
	// 		"phone": bson.M{
	// 			"bsonType":    "string",
	// 			"description": "must be a string and is required",
	// 		},
	// 		"name": bson.M{
	// 			"bsonType":    "string",
	// 			"description": "the endpoint IP address",
	// 		},
	// 	},
	// }
	db := client.Database(dbName)
	collection := db.Collection("numbers")
	// insertCommand := bson.D{{"insert", "numbers"}}
	// command := bson.D{{"validate", insertCommand}}

	// bson.D{{
	// 	collMod: "contacts",
	// 	validator: { $jsonSchema: {
	// 		bsonType: "object",
	// 		required: [ "name" ],
	// 		properties: {
	// 			phone: {
	// 				bsonType: "string",
	// 				description: "phone must be a string and is required"
	// 			},
	// 			name: {
	// 				bsonType: "string",
	// 				description: "name must be a string and is required"
	// 			}
	// 		}
	// 	} },
	// 	validationLevel: "strict"
	// }}

	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		panic(err)
	}
	id := res.InsertedID
	fmt.Println(res, id)
	// var result bson.M
	// err = coll.FindOne(context.TODO(), bson.D{{"title", title}}).
	// 	Decode(&result)
	// if err == mongo.ErrNoDocuments {
	// 	fmt.Printf("No document was found with the title %s\n", title)
	// 	return
	// }
	// if err != nil {
	// 	panic(err)
	// }

	// jsonData, err := json.MarshalIndent(result, "", "    ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", jsonData)
}
