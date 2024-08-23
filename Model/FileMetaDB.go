package Model

//
//import (
//	"log"
//)
//
//import (
//	_ "github.com/mattn/go-sqlite3"
//)
//
//type FileMetaDB struct {
//	dbPath string
//}
//
//func NewFileMetaDB() (metaTable *FileMetaDB) {
//
//	metaTable = &FileMetaDB{
//		dbPath: "file:db.db?cache=shared",
//	}
//
//	metaTable.createTable()
//
//	return metaTable
//}
//
///*
//refer : https://github.com/steffenfritz/FileTrove?tab=readme-ov-file
//action : 인덱싱을 위한 데이터 테이블을 생성합니다.
//- 인덱싱 테이블을 만들 때,
//
//단일 바이너리 애플리케이션이 디렉터리 트리를 탐색하고 Siegfried를 사용하여 모든 일반 파일을 형식별로 식별합니다.
//이때 다음 정보를 제공합니다 :
//1. MIME type
//- 파일의 인터넷 미디어 타입을 나타내며, 파일의 형식을 식별합니다.
//- ex. 이미지 파일은 image/jpeg,
//2. PRONOM identifier
//3. Format version
//4. Identification proof and note
//5. filename extension
//
//os.Stat() is giving you the
//6. File size
//7. File creation time
//8. File modification time
//9. File access time
//
//and the same for directories
//Furthermore it creates and calculates
//
//10. UUIDv4s as unique identifiers (not stable across sessions)
//11. hash sums (md5, sha1, sha256, sha512 and blake2b-512)
//12. the entropy of each file (up to 1GB)
//13. and it extracts some EXIF metadata and
//14. you can add your own DublinCore Elements metadata to scans.
//15. A very powerful feature is FileTrove's ability to consume YARA-X rule files.
//16. FileTrove also checks if the file is in the NSRL.
//*/
//func (t *FileMetaDB) createTable() {
//	sqlStmt := `
//		CREATE TABLE IF NOT EXISTS FileMetadata (
//			id INTEGER PRIMARY KEY AUTOINCREMENT,    -- 내부 ID, 자동 증가
//			uuid TEXT NOT NULL unique,                      -- UUIDv4
//			path TEXT NOT NULL,                      -- 파일 경로
//			name TEXT NOT NULL,                      -- 파일명
//			type TEXT NOT NULL,                      -- 파일 유형 (파일 또는 디렉터리)
//			mime_type TEXT,                          -- MIME 타입
//			pronom_id TEXT,                          -- PRONOM 식별자
//			format_version TEXT,                     -- 포맷 버전
//			identification_proof TEXT,               -- 식별 증거 및 노트
//			file_extension TEXT,                     -- 파일 확장자
//			file_size INTEGER,                       -- 파일 크기 (바이트)
//			creation_time DATETIME,                  -- 파일 생성 시간
//			modification_time DATETIME,              -- 파일 수정 시간
//			access_time DATETIME,                    -- 파일 접근 시간
//			md5 TEXT,                                -- MD5 해시
//			sha256 TEXT,                             -- SHA-256 해시
//			entropy REAL,                            -- 엔트로피 (최대 1GB 파일)
//			exif_metadata TEXT,                      -- EXIF 메타데이터
//			dublin_core TEXT,                        -- Dublin Core 메타데이터
//			yara_x_result TEXT,                      -- YARA-X 룰 결과
//			nsrl_check BOOLEAN,                      -- NSRL 체크 여부
//			scan_date DATETIME DEFAULT CURRENT_TIMESTAMP -- 스캔 날짜
//		);
//	`
//
//	_, err := t.db.Exec(sqlStmt)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
