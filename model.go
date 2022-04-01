package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dataBase *mongo.Database
)

type ScanURLParam struct {
	Url string `json:"url,omitempty"`
}

type Data struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type ScanURLAPIResultResultParam struct {
	Data Data `json:"data"`
}

type ResultDetail struct {
	Category string `json:"category"`
	Result   string `json:"result"`
}

type Results struct {
	URL         string       `json:"URL"`
	WebGuard    ResultDetail `json:"Web Security Guard"`
	SecureBrain ResultDetail `json:"SecureBrain"`
	Fortinet    ResultDetail `json:"Fortinet"`
	GoogleSafe  ResultDetail `json:"Google Safebrowsing"`
	Tencent     ResultDetail `json:"Tencent"`
}

type Attributes struct {
	Status  string  `json:"status"`
	Results Results `json:"results"`
}

type ReportData struct {
	Attributes Attributes `json:"attributes,omitempty"`
}

type GetReportAPIResultResultParam struct {
	Data ReportData `json:"data"`
}

func InitDB(dbURL string, dbName string) error {
	opt := options.Client().ApplyURI(dbURL)

	client, err := mongo.NewClient(opt)
	if err != nil {
		return err
	}
	if err := client.Connect(context.Background()); err != nil {
		return err
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		return err
	}
	dataBase = client.Database(dbName)
	return nil
}

func getAllData() (*[]Results, error) {
	results := &[]Results{}
	cur, err := dataBase.Collection(resultTable).Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cur.All(context.Background(), results)
	if err == nil {
		return results, nil
	}
	return nil, err
}

func getResult(url string) (*Results, error) {
	result := &Results{}
	filter := bson.M{
		"url": url,
	}
	err := dataBase.Collection(resultTable).FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func saveResult(param *GetReportAPIResultResultParam) error {
	_, err := dataBase.Collection(resultTable).InsertOne(context.Background(), param.Data.Attributes.Results)
	return err
}

