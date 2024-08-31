package Model

import (
	"fmt"
	"time"
)

/*
agent 의 상태 Status 는 총 3가지로 나뉜다.

Running : 동작중인 상태
waiting : 대기중인 상태
Stopping : 정지후 사라지기 전에 상태
*/
type AgentStatus int

const (
	Running  AgentStatus = iota // 동작 중인 상태
	Waiting                     // 대기 중인 상태
	Stopping                    // 정지 후 사라지기 전의 상태
)

/*
해당 코드를 빠른 코드 작성을 위해서 Chatgpt 가 작성후 허남정 연구원이 검토하는 형태로 만들었습니다.
*/

// AgentStatus를 문자열로 변환하는 메서드를 구현합니다.
func (s AgentStatus) String() string {
	switch s {
	case Running:
		return "Running"
	case Waiting:
		return "Waiting"
	case Stopping:
		return "Stopping"
	default:
		return "Unknown"
	}
}

type AgentStatusDB struct {
	dbName string
}

type AgentStatusRecord struct {
	ID        int
	UUID      string
	Status    AgentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAgentStatusDB creates a new instance of AgentStatusDB with the default table name.
func NewAgentStatusDB() *AgentStatusDB {
	return &AgentStatusDB{dbName: "AgentStatus"}
}

// CreateTable creates the AgentStatus table if it does not exist.
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
			status int,
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

// InsertRecord inserts a new record into the AgentStatus table.
func (s *AgentStatusDB) InsertRecord(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO %s (uuid, status) VALUES (?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.UUID, data.Status)
	if err != nil {
		return err
	}

	return nil
}

// SelectRecords retrieves all records from the AgentStatus table.
func (s *AgentStatusDB) SelectRecords() ([]AgentStatusRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, uuid, status, createAt, updateAt FROM %s`, s.dbName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []AgentStatusRecord
	for rows.Next() {
		var record AgentStatusRecord
		err := rows.Scan(&record.ID, &record.UUID, &record.Status, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// UpdateRecord updates the status of a record identified by its UUID.
func (s *AgentStatusDB) UpdateRecord(data *AgentStatusRecord) error {
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

// ExistRecord checks if at least one record exists in the AgentStatus table.
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
