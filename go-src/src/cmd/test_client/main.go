package main

import (
	"fmt"
	"net"
	"nro-go/pkg/protocol"
	"time"
)

func main() {
	fmt.Println("Connecting to server...")
	conn, err := net.Dial("tcp", "localhost:14445")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	session := protocol.NewSession(conn, 0)

	// Start reading loop in a goroutine to receive handshake
	go func() {
		for {
			// Read Command
			var cmdByte [1]byte
			_, err := conn.Read(cmdByte[:])
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			cmd := int8(cmdByte[0])

			if session.SentKey {
				// Decode command if key is set
				// Note: In this simple test client, we might not fully implement the session read loop
				// exactly like the server, but we need to decode to see what we got.
				// For simplicity, let's just print raw byte if we haven't implemented full client-side decoding yet.
				// But wait, the server sends -27 (GET_SESSION_ID) BEFORE encryption starts.
			}

			fmt.Printf("Received CMD: %d\n", cmd)

			if cmd == -27 {
				fmt.Println("Received Handshake (GET_SESSION_ID)!")
				// Read the rest of the message (Size + Data)
				// Size is short (2 bytes) because encryption is not yet active for this message
				// Actually, the server sends -27 using SendMessage.
				// Server logic:
				// msg := NewMessage(-27)
				// msg.WriteByte(1)
				// s.SendMessage(msg)
				//
				// Server writeLoop:
				// 1. Write CMD (-27)
				// 2. Write Size (Short 1) -> 0x00 0x01
				// 3. Write Data (1 byte) -> 0x01

				// So we expect: -27, 0, 1, 1

				var sizeShort int16
				// We need to read 2 bytes for size
				buf := make([]byte, 2)
				conn.Read(buf)
				// Manual decode big endian
				sizeShort = int16(buf[0])<<8 | int16(buf[1])

				fmt.Printf("Size: %d\n", sizeShort)

				data := make([]byte, sizeShort)
				conn.Read(data)
				fmt.Printf("Data: %v\n", data)

				fmt.Println("Handshake successful!")
				return
			}
		}
	}()

	// Keep main alive
	time.Sleep(5 * time.Second)
}
