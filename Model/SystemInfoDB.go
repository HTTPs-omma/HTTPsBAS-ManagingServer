package Model

import (
	"errors"
	"fmt"
	"time"
)

type SystemInfoDB struct {
	dbName string
}

type DsystemInfoDB struct {
	id            int
	Uuid          string
	HostName      string
	OsName        string
	OsVersion     string
	Family        string
	Architecture  string
	KernelVersion string
	BootTime      time.Time
	createAt      time.Time
	updateAt      time.Time
}

func NewSystemInfoDB() *SystemInfoDB {
	sysDB := &SystemInfoDB{"SystemInfo"}
	return sysDB
}

func (s *SystemInfoDB) createTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,    -- 내부 ID, 자동 증가
			uuid TEXT NOT NULL unique,               -- UUIDv4
			AgentUUID VARCHAR(255),
			HostName string,
			OsName string,
			OsVersion string,
			Family string,
			Architecture string,
			KernelVersion string,
			BootTime DATETIME,
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	sqlStmt = fmt.Sprintf(sqlStmt, s.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlModifyTrigger := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS update_ModificationTime
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN	
			UPDATE %s SET
				updateAt = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`, s.dbName, s.dbName)

	_, err = db.Exec(sqlModifyTrigger)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) insertRecord(data *DsystemInfoDB) error {
	// 데이터 베이스에는 단 하나의 Row 만을 보장해야함
	isExist, err := s.existRecord()
	if err != nil {
		return err
	}
	if isExist == true {
		err = s.updateRecord(data)
		if err != nil {
			return err
		}
		return nil
	}

	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO %s (uuid, HostName,
       OsName, OsVersion, Family, Architecture, KernelVersion,
       BootTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	fmt.Println(query)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(data.Uuid, data.HostName, data.OsName,
		data.OsVersion, data.Family, data.Architecture,
		data.KernelVersion, data.BootTime)
	//fmt.Println(rst.LastInsertId())
	//fmt.Println("debug=-===============")

	if err != nil {
		return err
	}

	return nil
}

/*
selectRecords()를 통해 반환된 DsystemInfoDB 객체의 값을 수정한 후,
수정된 객체를 updateRecord 함수의 매개변수로 전달합시오
*/
func (s *SystemInfoDB) updateRecord(data *DsystemInfoDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := s.selectRecords()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("SystemInfo 테이블에 저장된 데이터가 없습니다.")
	}
	row := rows[0]
	data.Uuid = row.Uuid

	query := fmt.Sprintf(`UPDATE %s SET HostName = ?, OsName = ?, OsVersion = ?, Family = ?, Architecture = ?, KernelVersion = ?, BootTime = ?`, s.dbName)
	_, err = db.Exec(query, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) deleteRecord(uuid string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s WHERE Uuid = ?`, s.dbName)
	_, err = db.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) selectRecords() ([]DsystemInfoDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s`, s.dbName)
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var rows []DsystemInfoDB

	for row.Next() {
		var data DsystemInfoDB

		err = row.Scan(&data.id, &data.Uuid, &data.HostName, &data.OsName,
			&data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion,
			&data.BootTime, &data.createAt, &data.updateAt)
		if err != nil {
			return nil, err
		}
		rows = append(rows, data)
	}

	return rows, nil
}

/*
*
하나 이상의 row 행이 있는지 검사한다.
*/
func (s *SystemInfoDB) existRecord() (bool, error) {
	db, err := getDBPtr()
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s)`, s.dbName)
	var exists bool
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
