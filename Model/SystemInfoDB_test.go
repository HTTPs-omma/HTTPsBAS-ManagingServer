package Model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"

	"agent/Extension"
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

func TestNewSystemInfoDB_insertRecord(t *testing.T) {
	type fields struct {
		dbName string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "insert DB test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			s := NewSystemInfoDB()

			sys, err := Extension.NewSysutils()
			if err != nil {
				t.Fatalf("Sysutils 생성 오류.")
			}

			charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
			var uuid string = ""
			for i := 0; i < 36; i++ {
				randI := rand.Intn(len(charset))
				uuid = uuid + string(charset[randI])
			}
			//fmt.Println("uuid : ", uuid)

			data := &DsystemInfoDB{
				Uuid:          uuid,
				HostName:      sys.GetHostName(),
				OsName:        sys.GetOsName(),
				OsVersion:     sys.GetOsVersion(),
				Family:        sys.GetFamily(),
				Architecture:  sys.GetArchitecture(),
				KernelVersion: sys.GetKernelVersion(),
				BootTime:      sys.GetBootTime(),
			}
			err = s.insertRecord(data)
			if err != nil {
				t.Fatalf(s.dbName + " : insert Record 오류. Query 재 확인 \n" + err.Error())
			}

			// ========== 검증 =============
			dbPtr, err := getDBPtr()
			if err != nil {
				t.Fatalf(s.dbName + " : DB 포인터를 가져올 수 없습니다. getDBPtr() 함수 오류\n" + err.Error())
			}

			query := fmt.Sprintf("select * from %s", s.dbName)
			dsys := &DsystemInfoDB{}
			row := dbPtr.QueryRow(query)
			if err != nil {
				t.Fatalf(s.dbName + " : select Qeury 오류\n" + err.Error())
			}

			rst := row.Scan(&dsys.id, &dsys.Uuid, &dsys.HostName, &dsys.OsName,
				&dsys.OsVersion, &dsys.Family, &dsys.Architecture, &dsys.KernelVersion,
				&dsys.BootTime, &dsys.createAt, &dsys.updateAt)
			if rst != nil {
				t.Fatalf(s.dbName + " : DsystemInfoDB 내용이 실제 DB 컬럼 내용과 일치하지 않습니다.\n" + err.Error())
			}

			assert.Equal(t, dsys.HostName, data.HostName)
			assert.Equal(t, dsys.Architecture, data.Architecture)
			assert.Equal(t, dsys.Family, data.Family)
			//assert.Equal(t, dsys.OsVersion, data.OsVersion)
			assert.Equal(t, dsys.BootTime, data.BootTime)
			assert.Equal(t, dsys.KernelVersion, data.KernelVersion)
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
