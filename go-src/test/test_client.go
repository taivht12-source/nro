package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"time"
)

var (
	clientID   = flag.Int("id", 1, "Client ID for testing")
	serverAddr = flag.String("addr", "localhost:14445", "Server address")
)

func main() {
	flag.Parse()

	fmt.Printf("🧪 Test Client #%d starting...\n", *clientID)
	fmt.Printf("Connecting to %s\n", *serverAddr)

	// 1. Connect to server
	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect: %v", err))
	}
	defer conn.Close()

	fmt.Println("✓ Connected to server")

	// 2. Read Session ID (Handshake)
	sessionIDBuf := make([]byte, 1)
	_, err = conn.Read(sessionIDBuf)
	if err != nil {
		panic(fmt.Sprintf("Failed to read session ID: %v", err))
	}
	sessionID := int8(sessionIDBuf[0])
	fmt.Printf("✓ Received Session ID: %d\n", sessionID)

	// 3. Read encryption key
	keyBuf := make([]byte, 1)
	_, err = conn.Read(keyBuf)
	if err != nil {
		panic(fmt.Sprintf("Failed to read key: %v", err))
	}
	key := keyBuf[0]
	fmt.Printf("✓ Received Key: %d\n", key)

	// 4. Initialize cipher state
	curR := byte(0)
	curW := byte(0)
	keyBytes := []byte{key}

	// Helper functions for encryption/decryption
	readKey := func(b byte) byte {
		result := (keyBytes[curR%byte(len(keyBytes))] & 0xFF) ^ (b & 0xFF)
		curR++
		return result
	}

	writeKey := func(b byte) byte {
		result := (keyBytes[curW%byte(len(keyBytes))] & 0xFF) ^ (b & 0xFF)
		curW++
		return result
	}

	// 5. Send PLAYER_MOVE packets
	fmt.Println("\n📤 Sending movement packets...")

	// Simulate player movement
	movements := []struct {
		x, y int16
	}{
		{100, 100},
		{150, 150},
		{200, 200},
		{250, 250},
		{300, 300},
	}

	for i, move := range movements {
		time.Sleep(500 * time.Millisecond) // Wait between moves

		// Create PLAYER_MOVE message (-7)
		// Format: [CMD: 1 byte] [Status: 1 byte] [X: 2 bytes] [Y: 2 bytes]
		cmd := int8(-7)
		status := int8(0)

		// Build message
		msg := make([]byte, 6)
		msg[0] = byte(cmd)
		msg[1] = byte(status)
		binary.BigEndian.PutUint16(msg[2:4], uint16(move.x))
		binary.BigEndian.PutUint16(msg[4:6], uint16(move.y))

		// Encrypt message
		encrypted := make([]byte, len(msg))
		for j := 0; j < len(msg); j++ {
			encrypted[j] = writeKey(msg[j])
		}

		// Send message length (1 byte) + encrypted message
		lengthByte := byte(len(encrypted))
		encryptedLength := writeKey(lengthByte)

		// Write to socket
		_, err = conn.Write([]byte{encryptedLength})
		if err != nil {
			fmt.Printf("❌ Failed to send length: %v\n", err)
			continue
		}

		_, err = conn.Write(encrypted)
		if err != nil {
			fmt.Printf("❌ Failed to send message: %v\n", err)
			continue
		}

		fmt.Printf("  [%d] Sent PLAYER_MOVE: X=%d, Y=%d\n", i+1, move.x, move.y)
	}

	fmt.Println("\n✓ All movement packets sent")
	fmt.Println("Listening for broadcasts from other players...")

	// 6. Listen for broadcast messages
	go func() {
		for {
			// Read message length
			lengthBuf := make([]byte, 1)
			_, err := conn.Read(lengthBuf)
			if err != nil {
				fmt.Printf("\n❌ Connection closed: %v\n", err)
				return
			}

			msgLen := readKey(lengthBuf[0])
			if msgLen == 0 {
				continue
			}

			// Read message
			msgBuf := make([]byte, msgLen)
			_, err = conn.Read(msgBuf)
			if err != nil {
				fmt.Printf("\n❌ Failed to read message: %v\n", err)
				return
			}

			// Decrypt message
			decrypted := make([]byte, msgLen)
			for i := 0; i < int(msgLen); i++ {
				decrypted[i] = readKey(msgBuf[i])
			}

			// Parse message
			if len(decrypted) >= 1 {
				cmd := int8(decrypted[0])
				if cmd == -7 && len(decrypted) >= 6 {
					playerID := int8(decrypted[1])
					x := int16(binary.BigEndian.Uint16(decrypted[2:4]))
					y := int16(binary.BigEndian.Uint16(decrypted[4:6]))
					fmt.Printf("📥 Received PLAYER_MOVE broadcast: PlayerID=%d, X=%d, Y=%d\n", playerID, x, y)
				} else {
					fmt.Printf("📥 Received message: CMD=%d, Length=%d\n", cmd, len(decrypted))
				}
			}
		}
	}()

	// Keep client alive
	fmt.Println("\nPress Ctrl+C to exit...")
	select {}
}
