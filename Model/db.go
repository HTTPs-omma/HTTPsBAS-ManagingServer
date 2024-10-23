package Model

import (
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

/*
sqlite3 refer : https://pkg.go.dev/modernc.org/sqlite#hdr-Connecting_to_a_database

db파일 : https://medium.com/@SlackBeck/golang-database-sql-패키지-삽질기-2편-sqlite-메모리-데이터베이스-c356dbd77e12
- cache=shared 옵션을 넣어주지 않으면 DB커넥션들은 하나의 데이터베이스를 공유하지 않는다.
- 즉, 데이터베이스 커넥션이 열려 있는 상태에서 새로운 커넥션을 열게 되면 기존 데이터베이스가 아닌 새로운 빈 데이터베이스를 할당받는다.
*/
func getDBPtr() (*sql.DB, error) {
	dbPath := "file:db.db?cache=shared"
	//dbPath := "db.db"
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return db, errors.New("open db failed")
	}

	db.SetMaxOpenConns(50)                 // 최대 오픈 커넥션 수
	db.SetMaxIdleConns(50)                 // 최대 유휴 커넥션 수;
	db.SetConnMaxLifetime(5 * time.Minute) // 연결이 닫히기 전에 열려 있는 최대 시간을 설정할 수 있습니다.
	// 최대 5분 까지 유지

	return db, nil
}
