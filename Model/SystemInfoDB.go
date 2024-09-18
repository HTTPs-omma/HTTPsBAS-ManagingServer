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
	ID            int       `json:"id"`
	Uuid          string    `json:"uuid"`
	HostName      string    `json:"host_name"`
	OsName        string    `json:"os_name"`
	OsVersion     string    `json:"os_version"`
	Family        string    `json:"family"`
	Architecture  string    `json:"architecture"`
	KernelVersion string    `json:"kernel_version"`
	BootTime      time.Time `json:"boot_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewSystemInfoDB() (*SystemInfoDB, error) {
	sysDB := &SystemInfoDB{"SystemInfo"}
	err := sysDB.CreateTable()
	if err != nil {
		return nil, err
	}
	return sysDB, nil
}

func (s *SystemInfoDB) CreateTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,    -- 내부 ID, 자동 증가
			uuid TEXT NOT NULL unique,               -- UUIDv4
			HostName string,
			OsName string,
			OsVersion string,
			Family string,
			Architecture string,
			KernelVersion string,
			BootTime DATETIME DEFAULT CURRENT_TIMESTAMP,
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	sqlStmt = fmt.Sprintf(sqlStmt, s.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		fmt.Println(63)
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

func (s *SystemInfoDB) InsertRecord(data *DsystemInfoDB) error {
	// fmt.Println("client uuid : " + data.Uuid)
	db, err := getDBPtr()
	if err != nil {

		return err
	}
	defer db.Close()

	checkQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE uuid = ?`, s.dbName)
	var count int
	err = db.QueryRow(checkQuery, data.Uuid).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// 중복된 항목이 있으면 업데이트
		return s.UpdateRecord(data)
	}

	query := fmt.Sprintf(`INSERT INTO %s (uuid, HostName,
       OsName, OsVersion, Family, Architecture, KernelVersion,
       BootTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)

	defer stmt.Close()
	if err != nil {
		fmt.Println(114)
		return err
	}

	_, err = stmt.Exec(data.Uuid, data.HostName, data.OsName,
		data.OsVersion, data.Family, data.Architecture,
		data.KernelVersion, data.BootTime)
	//fmt.Println(rst.LastInsertId())
	//fmt.Println("debug=-===============")

	if err != nil {
		fmt.Println(126)
		return err
	}

	return nil
}

/*
selectRecords()를 통해 반환된 DsystemInfoDB 객체의 값을 수정한 후,
수정된 객체를 updateRecord 함수의 매개변수로 전달합시오
*/
func (s *SystemInfoDB) UpdateRecord(data *DsystemInfoDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := s.SelectAllRecords()
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
		fmt.Println(156)
		return err
	}

	return nil
}

func (s *SystemInfoDB) DeleteRecordByUUID(uuid string) error {
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

func (s *SystemInfoDB) DeleteAllRecord() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s`, s.dbName)
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) SelectAllRecords() ([]DsystemInfoDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s`, s.dbName)
	row, err := db.Query(query)
	defer row.Close()
	if err != nil {
		fmt.Println(206)
		return nil, err
	}

	rows := []DsystemInfoDB{}

	for row.Next() {
		var data DsystemInfoDB
		err = row.Scan(&data.ID, &data.Uuid, &data.HostName, &data.OsName,
			&data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion,
			&data.BootTime, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			fmt.Println(218)
			return nil, err
		}
		rows = append(rows, data)
	}

	return rows, nil
}

func (s *SystemInfoDB) SelectRecordByUUID(uuid string) ([]DsystemInfoDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s WHERE uuid = ?`, s.dbName)
	rows, err := db.Query(query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []DsystemInfoDB{}
	for rows.Next() {
		data := DsystemInfoDB{}
		err = rows.Scan(&data.ID, &data.Uuid, &data.HostName, &data.OsName,
			&data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion,
			&data.BootTime, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			fmt.Println(248)
			return nil, err
		}
		results = append(results, data)
	}

	return results, nil
}
