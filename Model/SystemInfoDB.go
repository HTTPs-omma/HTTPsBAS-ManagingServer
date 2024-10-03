package Model

import (
	"encoding/json"
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
	IP            string    `json:"ip"`  // IP 추가
	MAC           string    `json:"mac"` // MAC 추가
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
			BootTime DATETIME,
			IP string,                               -- IP 추가
			MAC string,                              -- MAC 추가
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

func (s *SystemInfoDB) InsertRecord(data *DsystemInfoDB) error {
	// 데이터 베이스에는 단 하나의 Row 만을 보장해야함
	isExist, err := s.ExistRecordByUUID(data.Uuid)
	if err != nil {
		return err
	}
	if isExist == true {
		err = s.UpdateRecord(data)
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

	query := fmt.Sprintf(`INSERT INTO %s (uuid, HostName, OsName, OsVersion, Family, Architecture, KernelVersion, BootTime, IP, MAC) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(data.Uuid, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime, data.IP, data.MAC)

	if err != nil {
		return err
	}

	return nil
}

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

	query := fmt.Sprintf(`UPDATE %s SET HostName = ?, OsName = ?, OsVersion = ?, Family = ?, Architecture = ?, KernelVersion = ?, BootTime = ?, IP = ?, MAC = ?`, s.dbName)
	_, err = db.Exec(query, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime, data.IP, data.MAC)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) UpdateIP(IP string, Uuid string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET IP = ? WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, IP, Uuid)
	if err != nil {
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
		return nil, err
	}

	rows := []DsystemInfoDB{}

	for row.Next() {
		var data DsystemInfoDB

		err = row.Scan(&data.ID, &data.Uuid, &data.HostName, &data.OsName, &data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion, &data.BootTime, &data.IP, &data.MAC, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
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
		err = rows.Scan(&data.ID, &data.Uuid, &data.HostName, &data.OsName, &data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion, &data.BootTime, &data.IP, &data.MAC, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}

	return results, nil
}

func (s *SystemInfoDB) ExistRecord() (bool, error) {
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

func (s *SystemInfoDB) ExistRecordByUUID(uuid string) (bool, error) {
	db, err := getDBPtr()
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE uuid = ?)`, s.dbName)
	var exists bool
	err = db.QueryRow(query, uuid).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *SystemInfoDB) Unmarshal(data []byte) (*DsystemInfoDB, error) {

	var DsysInfo DsystemInfoDB
	err := json.Unmarshal(data, &DsysInfo)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return &DsysInfo, err
	}

	err = s.InsertRecord(&DsysInfo)
	if err != nil {
		return &DsysInfo, err
	}

	return &DsysInfo, err
}
