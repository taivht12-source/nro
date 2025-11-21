package protocol

import "fmt"

// Packet defines the interface for all network packets.
type Packet interface {
	Decode(msg *Message) error
}

// LoginRequest represents a login command (-1).
type LoginRequest struct {
	Username string
	Password string
}

func (p *LoginRequest) Decode(msg *Message) error {
	var err error
	p.Username, err = msg.ReadUTF()
	if err != nil {
		return fmt.Errorf("read username: %w", err)
	}
	p.Password, err = msg.ReadUTF()
	if err != nil {
		return fmt.Errorf("read password: %w", err)
	}
	return nil
}

// PlayerMoveRequest represents a movement command (-7).
type PlayerMoveRequest struct {
	X int16
	Y int16
}

func (p *PlayerMoveRequest) Decode(msg *Message) error {
	var err error
	p.X, err = msg.ReadShort()
	if err != nil {
		return fmt.Errorf("read x: %w", err)
	}
	p.Y, err = msg.ReadShort()
	if err != nil {
		return fmt.Errorf("read y: %w", err)
	}
	return nil
}

// UseSkillRequest represents a skill usage command (-11).
type UseSkillRequest struct {
	SkillIndex int8
	TargetID   int32
}

func (p *UseSkillRequest) Decode(msg *Message) error {
	var err error
	p.SkillIndex, err = msg.ReadByte()
	if err != nil {
		return fmt.Errorf("read skill index: %w", err)
	}
	// TargetID is Int in controller.go
	targetID, err := msg.ReadInt()
	if err != nil {
		return fmt.Errorf("read target id: %w", err)
	}
	p.TargetID = int32(targetID)
	return nil
}

// NPCInteractRequest represents an interaction with NPC (-20).
type NPCInteractRequest struct {
	NPCID int32
}

func (p *NPCInteractRequest) Decode(msg *Message) error {
	// NPCID is Int in controller_npc.go
	npcID, err := msg.ReadInt()
	if err != nil {
		return fmt.Errorf("read npc id: %w", err)
	}
	p.NPCID = int32(npcID)
	return nil
}

// TaskAcceptRequest represents accepting a task (-21).
type TaskAcceptRequest struct {
	TaskID int32
}

func (p *TaskAcceptRequest) Decode(msg *Message) error {
	taskID, err := msg.ReadInt()
	if err != nil {
		return fmt.Errorf("read task id: %w", err)
	}
	p.TaskID = int32(taskID)
	return nil
}

// TaskCompleteRequest represents completing a task (-22).
type TaskCompleteRequest struct {
	TaskID int32
}

func (p *TaskCompleteRequest) Decode(msg *Message) error {
	taskID, err := msg.ReadInt()
	if err != nil {
		return fmt.Errorf("read task id: %w", err)
	}
	p.TaskID = int32(taskID)
	return nil
}

// SkillSelectRequest represents selecting a skill (-25).
type SkillSelectRequest struct {
	SkillID int16
}

func (p *SkillSelectRequest) Decode(msg *Message) error {
	var err error
	p.SkillID, err = msg.ReadShort()
	if err != nil {
		return fmt.Errorf("read skill id: %w", err)
	}
	return nil
}
