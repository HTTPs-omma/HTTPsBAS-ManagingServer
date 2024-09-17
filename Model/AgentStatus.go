package Model

import (
	"fmt"
	"time"
)

// Protocol 유형을 정의합니다.
type Protocol uint8

const (
	TCP             Protocol = 0b0001 // TCP 프로토콜
	UDP             Protocol = 0b0010 // UDP 프로토콜
	HTTP            Protocol = 0b0011 // HTTP 프로토콜
	HTTPS           Protocol = 0b0100 // HTTPS 프로토콜
	UnknownProtocol Protocol = 0b0000 // 알 수 없는 프로토콜
)

// // Protocol을 문자열로 변환하는 메서드를 구현합니다.
func (p Protocol) String() string {
	switch p {
	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	case HTTP:
		return "HTTP"
	case HTTPS:
		return "HTTPS"
	default:
		return "Unknown"
	}
}

// // AgentStatus 유형을 정의합니다.
type AgentStatus int

const (
	NEW  AgentStatus = iota // 동작 중인 상태
	WAIT                    // 대기 중인 상태
	RUN                     // 정지 후 사라지기 전의 상태
	DELETED
	UNKNOWN
)

// AgentStatus를 문자열로 변환하는 메서드를 구현합니다.
func (s AgentStatus) String() string {
	switch s {
	case NEW:
		return "NEW"
	case RUN:
		return "Running"
	case WAIT:
		return "Waiting"
	case DELETED:
		return "DELETED"
	default:
		return "Unknown"
	}
}

// Binary 값을 AgentStatus로 변환하는 메서드를 구현합니다.
func BinaryToAgentStatus(i uint8) AgentStatus {
	switch i {
	case 0b00:
		return NEW
	case 0b01:
		return WAIT
	case 0b10:
		return RUN
	case 0b11:
		return DELETED
	default:
		return UNKNOWN
	}
}

// Binary 값을 AgentStatus로 변환하는 메서드를 구현합니다.
func BinaryToProtocol(i uint8) Protocol {
	switch i {
	case 0b0001:
		return TCP
	case 0b0010:
		return UDP
	case 0b0011:
		return HTTP
	case 0b0100:
		return HTTPS
	default:
		return UnknownProtocol
	}
}

type AgentStatusDB struct {
	dbName string
}

type AgentStatusRecord struct {
	ID        int
	UUID      string
	Status    AgentStatus
	Protocol  Protocol
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
			protocol int Default 0,
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

	query := fmt.Sprintf(`INSERT INTO %s (uuid, status, protocol) VALUES (?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(data.UUID, data.Status, data.Protocol)
	if err != nil {
		return err
	}

	return nil
}

// SelectRecords retrieves all records from the AgentStatus table.
func (s *AgentStatusDB) SelectAllRecords() ([]AgentStatusRecord, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT id, uuid, status, protocol, createAt, updateAt FROM %s`, s.dbName)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []AgentStatusRecord{}
	for rows.Next() {
		var record AgentStatusRecord
		err := rows.Scan(&record.ID, &record.UUID, &record.Status, &record.Protocol, &record.CreatedAt, &record.UpdatedAt)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// UpdateRecord updates the status and protocol of a record identified by its UUID.
func (s *AgentStatusDB) UpdateRecord(data *AgentStatusRecord) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`UPDATE %s SET status = ?, protocol = ? WHERE uuid = ?`, s.dbName)
	_, err = db.Exec(query, data.Status, data.Protocol, data.UUID)
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

	query := fmt.Sprintf(`DELETE FROM %s`)
	_, err = db.Exec(query)
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
