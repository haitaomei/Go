package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type UserType struct {
	Name     string
	Email    string
	Id       int
	Inactive bool
}

/*

docker run --rm --name mongodb-local \
-p 27017:27017 -p 27018:27018 -p 27019:27019 \
-e MONGO_INITDB_ROOT_USERNAME=admin \
-e MONGO_INITDB_ROOT_PASSWORD=adminpwd \
mongo

*/

func main() {
	os.Setenv("MONGO_SERVER", "localhost")
	os.Setenv("MONGO_PORT", "27017")
	os.Setenv("MONGO_USER_NAME", "admin")
	os.Setenv("MONGO_PASSWORD", "adminpwd")

	mongo := &MongoUtil{
		DBName:         "test",
		CollectionName: "Users",
	}
	mongo.InitMongoConnection()
	defer mongo.Disconnect()

	//// Insert one
	err := mongo.InsertOne(&UserType{"testUser", "testUser@ht.com", 1, false})
	if err != nil {
		fmt.Println("InsertOne fail")
	}

	// Insert many
	bill := &UserType{Name: "Bill", Email: "bill@ht.com"}
	billJr := &UserType{Name: "Bill", Email: "bill2@ht.com"}
	henry := &UserType{Name: "Henry", Email: "bill@ht.com"}
	err = mongo.InsertMany([]interface{}{bill, billJr, henry})
	if err != nil {
		fmt.Println("InsertMany fail")
	}

	/* ------------------------------------------------------------
		bson.D{{"email", "bill2@ht.com"}},
		Note, here must use ALL lower case !!!
		Unless specify `bson:"Email"` in the structure
	------------------------------------------------------------*/

	// Update one
	_, err = mongo.collection.UpdateOne(context.TODO(),
		bson.D{{"email", "bill2@ht.com"}},
		bson.D{{"$set", bson.D{{"inactive", true}}}})
	if err != nil {
		log.Fatal(err)
	}

	// Find one
	var aUser UserType
	err = mongo.collection.FindOne(context.TODO(), bson.D{{"name", "testUser"}}).Decode(&aUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Find a user", &aUser)

	// Find many
	findOptions := options.Find()
	findOptions.SetLimit(2)
	var results []*UserType
	cur, err := mongo.collection.Find(context.TODO(), bson.D{{"name", "Bill"}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var item UserType
		err := cur.Decode(&item)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &item)
	}

	for _, user := range results {
		fmt.Println(*user)
	}

	mongo.DeleteAll()
}

// MongoDb Util
type MongoUtil struct {
	DBName         string
	CollectionName string

	client     *mongo.Client
	collection *mongo.Collection
}

func (u *MongoUtil) InitMongoConnection() {
	var mongoConnectStr string
	var mongoServer string
	var mongoPort string
	var mongoUser string
	var mongoPassword string
	ok := false
	if mongoServer, ok = os.LookupEnv("MONGO_SERVER"); !ok {
		fmt.Errorf("MongoDB server is not specified")
	}
	if mongoPort, ok = os.LookupEnv("MONGO_PORT"); !ok {
		fmt.Errorf("MongoDB port is not specified")
	}
	if mongoUser, ok = os.LookupEnv("MONGO_USER_NAME"); !ok {
		fmt.Errorf("MongoDB username is not specified")
	}
	if mongoPassword, ok = os.LookupEnv("MONGO_PASSWORD"); !ok {
		fmt.Errorf("MongoDB passowrd is not specified")
	}

	mongoConnectStr = "mongodb://" + mongoUser + ":" + mongoPassword + "@" + mongoServer + ":" + mongoPort + "/"

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoConnectStr,
	))
	if err != nil {
		log.Fatal("error connecting to mongoDB", err)
	}
	u.client = client
	log.Println("Connected to MongoDB!")

	u.collection = u.client.Database(u.DBName).Collection(u.CollectionName)
}

func (u *MongoUtil) Disconnect() {
	err := u.client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to MongoDB closed.")
}

func (u *MongoUtil) InsertOne(document interface{}) error {
	ret, err := u.collection.InsertOne(context.TODO(), document)
	if err != nil {
		return err
	}
	fmt.Println("Inserted a single document: ", ret.InsertedID)
	return nil
}

func (u *MongoUtil) InsertMany(document []interface{}) error {
	ret, err := u.collection.InsertMany(context.TODO(), document)
	if err != nil {
		return err
	}
	fmt.Println("Inserted multiple documents: ", ret.InsertedIDs)
	return nil
}

func (u *MongoUtil) DeleteAll() {
	deleteResult, err := u.collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the %s collection\n", deleteResult.DeletedCount, u.CollectionName)
}
