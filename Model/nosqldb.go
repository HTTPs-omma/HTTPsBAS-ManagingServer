package Model

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Chatgpt 이용함.
/*
MongoDB는 100개의 커넥션 풀을 지원한다.
- 최대 풀 크기 50, 최소 풀 크기 10으로 설정.
- 커넥션 최대 유휴 시간을 30초로 설정.
*/
func getCollectionPtr() (*mongo.Database, error) {
	MONGOID := os.Getenv("MONGODBID")
	MONGOPW := os.Getenv("MONGODBPW")
	SERVER_DOMAIN := os.Getenv("SERVER_DOMAIN")
	MONGOPORT := os.Getenv("MONGOPORT")

	clientOptions := options.Client().
		ApplyURI("mongodb://" + MONGOID + ":" + MONGOPW + "@" + SERVER_DOMAIN + ":" + MONGOPORT). // MongoDB URI
		SetMaxPoolSize(20).                                                                       // 최대 풀 크기
		SetMinPoolSize(10).                                                                       // 최소 풀 크기
		SetMaxConnIdleTime(60 * time.Second)                                                      // 최대 유휴 시간
		//fmt.Println("mongodb://" + MONGOID + ":" + MONGOPW + "@uskawjdu.iptime.org:17017/")
	// 클라이언트 생성
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Error creating MongoDB client: ", err)
		return nil, err
	}

	// 연결 설정
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// MongoDB에 연결
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB: ", err)
		return nil, err
	}

	database := client.Database("httpsAgent")

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// 클라이언트 연결 반환
	return database, nil
}
