package network

import (
	"fmt"
	"net"
	"nro/src/pkg/protocol"
)

// TCPServer quản lý việc lắng nghe kết nối TCP.
type TCPServer struct {
	Port       int
	Controller *Controller
}

// NewTCPServer tạo một server mới.
func NewTCPServer(port int, controller *Controller) *TCPServer {
	return &TCPServer{
		Port:       port,
		Controller: controller,
	}
}

// Start bắt đầu lắng nghe kết nối.
func (s *TCPServer) Start() error {
	addr := fmt.Sprintf(":%d", s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Printf("Server started on port %d\n", s.Port)

	idCounter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		idCounter++
		session := protocol.NewSession(conn, idCounter)
		session.Handler = s.Controller // Gán Controller làm Handler
		fmt.Printf("New connection: %s (ID: %d)\n", conn.RemoteAddr().String(), idCounter)

		// Bắt đầu xử lý session
		session.Start()
	}
}
