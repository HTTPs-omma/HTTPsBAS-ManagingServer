package Model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type sqlite_master struct {
	Type     string
	name     string
	tbl_name string
	rootpage string
	sql      string
}

func TestSystemInfoDB_createTable(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create DB test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			s := NewSystemInfoDB()
			s.createTable()
			if err != nil {
				t.Fatalf(s.dbName + " : DB를 생성할 수 없습니다. \n" + err.Error())
			}

			// ========== 검증 =============
			dbPtr, err := getDBPtr()
			if err != nil {
				t.Fatalf(s.dbName + " : DB 포인터를 가져올 수 없습니다. getDBPtr() 함수 오류\n" + err.Error())
			}

			query := fmt.Sprintf("select * from sqlite_master where name = '%s'", s.dbName)

			dsys := &sqlite_master{}

			rst := dbPtr.QueryRow(query).Scan(&dsys.Type, &dsys.name, &dsys.tbl_name, &dsys.rootpage, &dsys.sql)
			if rst != nil {
				t.Fatalf(s.dbName + " : 생성된 테이블이 존재하지 않습니다.")
			}
			assert.Equal(t, dsys.name, s.dbName)

		})
	}
}

func TestNewSystemInfoDB_selectRecord(t *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Select Record in DB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			s := NewSystemInfoDB()

			data, err := s.selectRecords()

			if err != nil {
				t.Fatalf(s.dbName + " : select 오류. Query 재 확인 \n" + err.Error())
			}

			assert.Equal(t, 1, len(data), "Row를 하나 이상 채우고 재시도를 하시오")
			fmt.Println("sql select result : \n", data)

		})
	}
}
