package Model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Chatgpt 이용함.
/*
MongoDB는 100개의 커넥션 풀을 지원한다.
- 최대 풀 크기 50, 최소 풀 크기 10으로 설정.
- 커넥션 최대 유휴 시간을 30초로 설정.
*/
func getNoSqlDbPtr() (*mongo.Database, error) {
	// MongoDB 클라이언트 옵션 설정
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017"). // MongoDB URI
		SetMaxPoolSize(50).                    // 최대 풀 크기
		SetMinPoolSize(10).                    // 최소 풀 크기
		SetMaxConnIdleTime(30 * time.Second)   // 최대 유휴 시간

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

	// 사용할 데이터베이스 선택 (예: "mydatabase")
	database := client.Database("mydatabase")

	// 프로시저 실행
	// 성공 ( success )
	// 실패 ( failure )
	// 경고 ( Warning )
	//S

	// 클라이언트 연결 반환
	return database, nil
}
