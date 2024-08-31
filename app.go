package main

import (
	"errors"
	"fmt"
	"github.com/HTTPs-omma/HSProtocol/HSProtocol"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"net"
	"sync"
)

var queue *HSQueue // 패키지 수준에서 전역 큐 생성

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	queue = NewHSQueue()

	if err != nil {
		panic("큐 생성 에러")
	}

	// tcp
	go TCPServer()

	// HTTP
	//go HTTPServer()

	// udp
	//go UDPServer()

	for {
		if queue.HasNext() {
			hs, err := queue.Dequeue()
			if err != nil {
				fmt.Errorf("queue.Dequeue() 에러 : ", err)
				continue
			}
			fmt.Println("uuid : ", hs.UUID)
			fmt.Println("command : ", hs.Command)
		}
	}

}

func TCPServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close() // 함수 호출 종료 후 Close

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			break
		}
		if n < 1 {
			continue
		}

		HSMgr := HSProtocol.NewHSProtocolManager()
		hs, err := HSMgr.Parsing(buffer[:n])
		if err != nil {
			fmt.Println("Error parsing:", err)
			continue
		}
		// 받은 데이터를 PacketQueue에 추가
		err = queue.Enqueue(*hs)

		if err != nil {
			fmt.Println("Error adding packet to queue:", err)
		} else {
			fmt.Printf("Received and queued packet: %x\n", buffer[:n])
		}
	}
}

// HTTP 서버 함수 (Fiber 사용)
func HTTPServer() {
	app := fiber.New()

	app.Get("/status", func(c *fiber.Ctx) error {

		return c.SendString("문자")
	})
	app.Get("/next", queueNextHandler)

	fmt.Println("HTTP server listening on port 8081")
	err := app.Listen(":80")
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

const MaxQueueSize = 10000 // 큐의 최대 크기

// HS 타입을 위한 큐 구조체 정의
type HSQueue struct {
	items []HSProtocol.HS // HS 타입의 요소를 저장할 슬라이스
	lock  sync.Mutex      // 큐 접근 동기화를 위한 뮤텍스
	cond  *sync.Cond      // 큐의 상태를 확인하기 위한 조건 변수
}

// 새로운 HSQueue 생성
func NewHSQueue() *HSQueue {
	q := &HSQueue{
		items: make([]HSProtocol.HS, 0, MaxQueueSize), // 큐의 크기를 설정
	}
	q.cond = sync.NewCond(&q.lock) // 조건 변수 초기화
	return q
}

// 큐에 요소 추가 (Enqueue)
func (q *HSQueue) Enqueue(item HSProtocol.HS) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if len(q.items) >= MaxQueueSize {
		return errors.New("queue is full")
	}

	q.items = append(q.items, item)
	q.cond.Signal() // 대기 중인 고루틴에 신호를 보냄
	return nil
}

// 큐에 첫 번째 값이 있는지 확인 (HasNext)
func (q *HSQueue) HasNext() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	// 큐가 비어 있지 않으면 true 반환
	return len(q.items) > 0
}

// 큐에서 요소 제거 (Dequeue)
func (q *HSQueue) Dequeue() (HSProtocol.HS, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.items) == 0 {
		q.cond.Wait() // 큐가 비어있으면 대기
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// 큐의 현재 크기 확인
func (q *HSQueue) Size() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.items)
}
