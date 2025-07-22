package tests

import (
	"redis/command"
	"redis/internal/db"
	"testing"
)

func TestSetAndGetCommand(t *testing.T) {
	s := &db.Store{Mpp: make(map[string]string)}
	command.SetCommand("foo", "bar", s)
	val := command.GetCommand("foo", s)
	if val != "bar" {
		t.Errorf("expected 'bar', got '%s'", val)
	}
}

func TestINCRcommand(t *testing.T) {
	s := &db.Store{Mpp: make(map[string]string)}
	val, err := command.INCRcommand("counter", s)
	if err != nil || val != 1 {
		t.Errorf("expected 1, got %d, err: %v", val, err)
	}
	val, err = command.INCRcommand("counter", s)
	if err != nil || val != 2 {
		t.Errorf("expected 2, got %d, err: %v", val, err)
	}
}

func TestINCRBYcommand(t *testing.T) {
	s := &db.Store{Mpp: make(map[string]string)}
	val, err := command.INCRBYcommand("counter", "5", s)
	if err != nil || val != 5 {
		t.Errorf("expected 5, got %d, err: %v", val, err)
	}
	val, err = command.INCRBYcommand("counter", "3", s)
	if err != nil || val != 8 {
		t.Errorf("expected 8, got %d, err: %v", val, err)
	}
}

func TestDeletecommand(t *testing.T) {
	s := &db.Store{Mpp: make(map[string]string)}
	command.SetCommand("foo", "bar", s)
	command.SetCommand("baz", "qux", s)
	count := command.Deletecommand([]string{"foo", "baz"}, s)
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
	val := command.GetCommand("foo", s)
	if val != "" {
		t.Errorf("expected '', got '%s'", val)
	}
}
