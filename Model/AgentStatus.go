package Model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
)

// Binary 값을 AgentStatus로 변환하는 메서드를 구현합니다.
func BinaryToAgentStatus(i uint8) HSProtocol.AGENTSTATUS {
	switch i {
	case 0b00:
		return HSProtocol.NEW
	case 0b01:
		return HSProtocol.WAIT
	case 0b10:
		return HSProtocol.RUN
	case 0b11:
		return HSProtocol.DELETED
	default:
		return HSProtocol.NEW
	}
}

// Binary 값을 AgentStatus로 변환하는 메서드를 구현합니다.
func BinaryToProtocol(i uint8) HSProtocol.PROTOCOL {
	switch i {
	case 0b0001:
		return HSProtocol.TCP
	case 0b0010:
		return HSProtocol.UDP
	case 0b0011:
		return HSProtocol.HTTP
	case 0b0100:
		return HSProtocol.HTTPS
	default:
		return HSProtocol.UNKNOWN
	}
}

type AgentStatusDB struct {
	dbName string
}

type AgentStatusRecord struct {
	ID        int
	UUID      string
	Status    HSProtocol.AGENTSTATUS
	Protocol  HSProtocol.PROTOCOL
	NickName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAgentStatusDB creates a new instance of AgentStatusDB with the default table name.
func NewAgentStatusDB() (*AgentStatusDB, error) {
	db := &AgentStatusDB{dbName: "AgentStatus"}
	err := db.CreateTable()
	if err != nil {
		return nil, err
	}
	return db, nil
}
func (s *AgentStatusDB) CreateTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL UNIQUE,
			nickName TEXT,
			status int,
			protocol int DEFAULT 0,
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	sqlStmt = fmt.Sprintf(sqlStmt, s.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlTrigger := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS update_ModificationTime
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN	
			UPDATE %s SET
				updateAt = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`, s.dbName, s.dbName)

	_, err = db.Exec(sqlTrigger)
	if err != nil {
		return err
	}

	return nil
}

// InsertRecord 메서드에서 NickName 필드를 추가하여 레코드 삽입
func (s *AgentStatusDB) InsertRecord(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	checkQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE uuid = ?`, s.dbName)
	var count int
	err = db.QueryRow(checkQuery, data.UUID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		// 중복된 항목이 있으면 업데이트
		return s.UpdateRecordNoNick(data)
	}

	query := fmt.Sprintf(`INSERT INTO %s (uuid, nickName, status, protocol) VALUES (?, ?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.UUID, data.NickName, data.Status, data.Protocol)
	if err != nil {
		return err
	}

	return nil
}

// SelectAllRecords 메서드에서 NickName 필드 추가
func (s *AgentStatusDB) SelectAllRecords() ([]AgentStatusRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, uuid, nickName, status, protocol, createAt, updateAt FROM %s`, s.dbName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []AgentStatusRecord{}
	for rows.Next() {
		var record AgentStatusRecord
		err := rows.Scan(&record.ID, &record.UUID, &record.NickName, &record.Status, &record.Protocol, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// SelectRecordByUUID 메서드에서 NickName 필드 추가
func (s *AgentStatusDB) SelectRecordByUUID(uuid string) ([]AgentStatusRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, uuid, nickName, status, protocol, createAt, updateAt FROM %s WHERE uuid = ?`, s.dbName)
	row := db.QueryRow(query, uuid)

	var records []AgentStatusRecord
	var record AgentStatusRecord
	err = row.Scan(&record.ID, &record.UUID, &record.NickName, &record.Status, &record.Protocol, &record.CreatedAt, &record.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no record is found
		}
		return nil, err
	}

	return append(records, record), nil
}

// UpdateRecord 메서드에서 NickName 필드를 업데이트하도록 수정
func (s *AgentStatusDB) UpdateRecord(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET status = ?, nickName = ? WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, data.Status, data.NickName, data.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) UpdateRecordNoNick(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET status = ? WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, data.Status, data.UUID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecord deletes a record from the AgentStatus table based on its UUID.
func (s *AgentStatusDB) DeleteRecord(uuid string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *AgentStatusDB) DeleteAllRecord() error {
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

func (s *AgentStatusDB) ExistRecord() (bool, error) {
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
