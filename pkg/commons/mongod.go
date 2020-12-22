package commons

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://docs.mongodb.com/manual/crud/

// Trainer details
type Trainer struct {
	Name string
	Age  int
	City string
}

type mongoDB struct {
	client *mongo.Client
}

// DB declarations
const (
	ChatDatabase      = "Chat"
	ChatCollection    = "Chat"
	MessageCollection = "Message"
)

// MongoDB APIs
type MongoDB interface {
	Client() *mongo.Client
	Disconnect(context.Context)

	Create(ctx context.Context, d string, c string, item interface{}) error
	Read(ctx context.Context, d string, c string, filter interface{}) *mongo.SingleResult
	Update(ctx context.Context, d string, c string, filter interface{}, update interface{}) error
	Delete(ctx context.Context, d string, c string, filter interface{}) *mongo.SingleResult

	Upsert(ctx context.Context, d string, c string, filter interface{}, update interface{}) error

	BatchRead(ctx context.Context, d string, c string, filter interface{}) ([]*bson.Raw, error)
	Query(ctx context.Context, d string, c string, sort *SortData, itemsRange *RangeData, filter *FilterData) (int64, []*bson.Raw, error)
}

// InitMongoDB function
func InitMongoDB(ctx context.Context) (MongoDB, error) {
	mongodb := &mongoDB{}
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// // Connect to MongoDB
	// client, err := mongo.Connect(context.TODO(), clientOptions)
	// if err != nil {
	// 	return nil, err
	// }

	// Connect to MongoDB
	mongoServer := os.Getenv("MONGODB_SERVER")
	mongoPort := os.Getenv("MONGODB_PORT")
	mongoUrl := "mongodb://localhost:27017"
	if mongoServer == "" || mongoPort == "" {
		fmt.Printf("No mongo config\n\tuse default\n")
		mongoServer = "localhost"
		mongoPort = "27017"
	} else {
		mongoUrl = fmt.Sprintf("mongodb://%s:%s", mongoServer, mongoPort)
		fmt.Printf("Mongo config\n\tconnect to %s\n", mongoUrl)
	}
	mongoUsername := os.Getenv("MONGODB_USERNAME")
	mongoPassword := os.Getenv("MONGODB_PASSWORD")
	if mongoUsername == "" || mongoPassword == "" {
		fmt.Printf("No auth config\n\tuse default\n")
		mongoUsername = "root"
		mongoPassword = "example"
	}
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions = clientOptions.SetAuth(options.Credential{Username: mongoUsername, Password: mongoPassword})
	if os.Getenv("MONGODB_TLS") == "true" {
		roots, _ := x509.SystemCertPool()
		selfCA := os.Getenv("SELF_CA")
		mongoCert := os.Getenv("MONGODB_CERT")
		if selfCA != "" {
			crt, _ := ioutil.ReadFile(selfCA)
			roots.AppendCertsFromPEM(crt)
		}
		if mongoCert != "" {
			crt, _ := ioutil.ReadFile(mongoCert)
			roots.AppendCertsFromPEM(crt)
		}
		clientOptions = clientOptions.SetTLSConfig(&tls.Config{
			// InsecureSkipVerify: true,
			RootCAs:    roots,
			ServerName: mongoServer,
		})
	}
	// TODO: credential
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error getting connect mongo client: %v", err)
	}
	mongodb.client = mongoClient

	return mongodb, nil
}

func (db *mongoDB) Client() *mongo.Client {
	return db.client
}

func (db *mongoDB) Disconnect(ctx context.Context) {
	db.client.Disconnect(ctx)
}

func (db *mongoDB) Create(ctx context.Context, d string, c string, item interface{}) error {
	collection := db.client.Database(d).Collection(c)
	insertResult, err := collection.InsertOne(ctx, item)
	fmt.Printf("%v %v\n", insertResult, err)
	return err
}

func (db *mongoDB) Read(ctx context.Context, d string, c string, filter interface{}) *mongo.SingleResult {
	collection := db.client.Database(d).Collection(c)
	// value.Decode(struct)
	return collection.FindOne(ctx, filter)
}

func (db *mongoDB) Update(ctx context.Context, d string, c string, filter interface{}, update interface{}) error {
	collection := db.client.Database(d).Collection(c)
	// update = bson.D{
	//     {"$set", bson.D{
	//         {"key", "value"},
	//     }},
	// }
	// filter = bson.D{{"id", "1WEB"}}
	result := collection.FindOneAndUpdate(ctx, filter, update)
	return result.Err()
}

func (db *mongoDB) Upsert(ctx context.Context, d string, c string, filter interface{}, update interface{}) error {
	collection := db.client.Database(d).Collection(c)
	// update = bson.D{
	//     {"$set", bson.D{
	//         {"key", "value"},
	//     }},
	// }
	// filter = bson.D{{"id", "1WEB"}}
	opt := &options.UpdateOptions{}
	opt = opt.SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opt)
	return err
}

func (db *mongoDB) Delete(ctx context.Context, d string, c string, filter interface{}) *mongo.SingleResult {
	collection := db.client.Database(d).Collection(c)
	// update = bson.D{
	//     {"$set", bson.D{
	//         {"key", "value"},
	//     }},
	// }
	// filter = bson.D{{"id", "1WEB"}}
	return collection.FindOneAndDelete(ctx, filter)
	// return err
}

func (db *mongoDB) BatchRead(ctx context.Context, d string, c string, filter interface{}) ([]*bson.Raw, error) {
	collection := db.client.Database(d).Collection(c)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	raws := []*bson.Raw{}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		raw := cur.Current
		raws = append(raws, &raw)
	}
	// bson.Unmarshal(raws[1], struct)
	return raws, nil
}

// SortData ...
type SortData struct {
	Item string
	Desc bool // true as descending order
}

// RangeData ...
type RangeData struct {
	From int
	To   int
}

// FilterData ...
type FilterData struct {
	Item  string
	Value interface{}
}

func (db *mongoDB) Query(ctx context.Context, d string, c string, sort *SortData, itemsRange *RangeData, filter *FilterData) (int64, []*bson.Raw, error) {
	collection := db.client.Database(d).Collection(c)

	var cursor *mongo.Cursor
	var err error
	var count int64
	var raws []*bson.Raw

	findOptions := options.Find()
	// set range
	if itemsRange != nil {
		findOptions.SetSkip(int64(itemsRange.From))
		findOptions.SetLimit(int64(itemsRange.To + 1 - itemsRange.From))
	}

	// set sorter
	if sort != nil {
		order := 1
		if sort.Desc {
			order = -1
		}
		findOptions.SetSort(bson.D{{Key: sort.Item, Value: order}})
	}

	// set filter
	if filter != nil {
		if filter.Item == "" {
			cursor, err = collection.Find(
				ctx, filter.Value, findOptions,
			)
		} else {
			cursor, err = collection.Find(
				ctx, bson.D{
					{Key: filter.Item, Value: filter.Value},
				}, findOptions,
			)
		}
		if err != nil {
			return 0, nil, err
		}
		count, err = collection.CountDocuments(ctx, bson.D{
			{Key: filter.Item, Value: filter.Value},
		})
		if err != nil {
			return 0, nil, err
		}
	} else {
		cursor, err = collection.Find(ctx, bson.D{{}}, findOptions)
		if err != nil {
			return 0, nil, err
		}
		count, err = collection.CountDocuments(ctx, bson.D{{}})
		if err != nil {
			return 0, nil, err
		}
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		// raw := &bson.Raw{}
		// if err = cursor.Decode(&raw); err != nil {
		// 	return 0, nil, err
		// }
		raw := cursor.Current
		raws = append(raws, &raw)
	}

	return count, raws, nil
}
