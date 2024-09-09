package Model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

/*
빠른 개발을 위해  chatgpt 를 사용한 개발 코드입니다.
검토자 : 허남정
*/


// OperationLogDocument 구조체 정의
type OperationLogDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	AgentUUID    string             `bson:"agentUUID"`
	ProcedureID  string             `bson:"procedureID"`
	InstructionUUID  string         `bson:"instructionUUID"`
	ConductAt    time.Time          `bson:"conductAt"`
	ExitCode     int                `bson:"exitCode"`
	Log          string             `bson:"log"`
	Command      string             `bson:"command"` // Command 필드로 변경
}

// OperationLogDB는 CRUD 작업을 수행하는 구조체입니다.
type OperationLogDB struct {
	Collection *mongo.Collection
}

// NewOperationLogDB는 OperationLogDB 인스턴스를 생성합니다.
func NewOperationLogDB() (*OperationLogDB, error) {
	db, err := getCollectionPtr()
	if err != nil {
		return nil, err
	}
	return &OperationLogDB{
		Collection: db.Collection("execLog"),
	}, nil
}

// Create는 새로운 OperationLogDocument 문서를 MongoDB에 삽입합니다.
func (repo *OperationLogDB) insertDocument(log OperationLogDocument) (*mongo.InsertOneResult, error) {
	// Command 필드를 기본값으로 설정 (필요에 따라 변경)
	result, err := repo.Collection.InsertOne(context.TODO(), log)
	fmt.Println(log)
	if err != nil {
		return nil, err
	}

	fmt.Println("Inserted document with ID:", result.InsertedID)
	return result, nil
}

// Read는 OperationLogDocument 문서를 조회합니다.
func (repo *OperationLogDB) selectDocumentById(id string) (*OperationLogDocument, error) {
	var OperationLogDocument OperationLogDocument
	filter := bson.M{"instructionUUID": id}

	err := repo.Collection.FindOne(context.TODO(), filter).Decode(&OperationLogDocument)
	if err != nil {
		return nil, err
	}

	return &OperationLogDocument, nil
}

// Update는 OperationLogDocument 문서를 수정합니다.
func (repo *OperationLogDB) UpdateDocumentByInstID(id string, updateData bson.M) (*mongo.UpdateResult, error) {
	filter := bson.M{"instructionUUID": id}
	update := bson.M{
		"$set": updateData,
	}

	result, err := repo.Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	fmt.Println("Updated document count:", result.ModifiedCount)
	return result, nil
}

// Delete는 OperationLogDocument 문서를 삭제합니다.
func (repo *OperationLogDB) DeleteDocumentByInstID(id string) (*mongo.DeleteResult, error) {
	filter := bson.M{"instructionUUID": id}

	result, err := repo.Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	fmt.Println("Deleted document count:", result.DeletedCount)
	return result, nil
}
