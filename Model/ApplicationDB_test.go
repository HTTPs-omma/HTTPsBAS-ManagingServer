package Model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationDB_createTable(t1 *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Test case 1"},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			appdb := NewApplicationDB()
			err := appdb.createTable()
			if err != nil {
				t1.Fatalf("Error creating table: %v", err)
			}

			// ========== 검증 =============
			dbPtr, err := getDBPtr()
			if err != nil {
				t1.Fatalf(appdb.dbName + " : DB 포인터를 가져올 수 없습니다. getDBPtr() 함수 오류\n" + err.Error())
			}

			query := fmt.Sprintf("select * from sqlite_master where name = '%s'", appdb.dbName)

			dsys := &sqlite_master{}

			rst := dbPtr.QueryRow(query).Scan(&dsys.Type, &dsys.name, &dsys.tbl_name, &dsys.rootpage, &dsys.sql)
			if rst != nil {
				t1.Fatalf(appdb.dbName + " : 생성된 테이블이 존재하지 않습니다.")
			}
			assert.Equal(t1, dsys.name, appdb.dbName)
		})
	}
}

func TestApplicationDB_CRUD(t1 *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{name: "Test case 1"},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
		})
	}
}

//func TestApplicationDB_CreateAll(t1 *testing.T) {
//	type fields struct {
//		dbName string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//	}{
//		{name: "Test case 1"},
//	}
//	for _, tt := range tests {
//		t1.Run(tt.name, func(t *testing.T) {
//			// ============ Create 테스트 ================
//			appdb := NewApplicationDB()
//			appdb.createTable()
//
//			wind32 := getApplicationList()
//			for _, wind := range wind32 {
//				data := DapplicationDB{}
//				data.Name = wind.Name
//				data.Description = wind.Description
//				data.Version = wind.Version
//				data.Vendor = wind.Vendor
//				data.InstallDate2 = wind.InstallDate2
//				data.InstallLocation = wind.InstallLocation
//				data.InstallSource = wind.InstallSource
//				data.Language = wind.Language
//				data.PackageCode = wind.PackageCode
//				data.PackageName = wind.PackageName
//				data.RegCompany = wind.RegCompany
//				data.RegOwner = wind.RegOwner
//				data.URLInfoAbout = wind.URLInfoAbout
//
//				err := appdb.insertRecord(data)
//				if err != nil {
//					t.Fatalf(appdb.dbName + " : insert 에러\n" + err.Error())
//				}
//			}
//		})
//	}
//}
