package storage_test

import (
	"degrens/panel/internal/storage"
	"testing"
)

func TestInitStateTokenStorage(t *testing.T) {
	// STS should not be nil
	if storage.GetStateTokenStorage() == nil {
		t.Error("StateTokenStorage cannot be nil")
		return
	}
}

func TestStateTokenStorage_Add(t *testing.T) {
	// Should return nil
	if err := storage.GetStateTokenStorage().Add("test", "test"); err != nil {
		t.Error("StateTokenStorage.Add() must return nil", "err:", err)
		return
	}
}

func TestStateTokenStorage_Get(t *testing.T) {
	STS := storage.GetStateTokenStorage()
	STS.Clear()
	// STS should be empty
	if _, err := STS.Get("test"); err == nil {
		t.Error("STS should be empty")
		return
	}

	// Add a value
	err := STS.Add("test", "test")
	if err != nil {
		return
	}

	// Should return the value
	if value, err := STS.Get("test"); err != nil || value != "test" {
		t.Error("StateTokenStorage.Get() must return \"test\"")
		return
	}
}

func TestStateTokenStorage_Remove(t *testing.T) {
	STS := storage.GetStateTokenStorage()
	STS.Clear()

	// Removing an inexistent key should return an error
	if err := STS.Remove("test"); err == nil {
		t.Error("STS.Remove() must return an error if the key does not exist", "STS:", STS)
		return
	}

	// Add a value
	err := STS.Add("test", "test")
	if err != nil {
		return
	}

	// Should return nil
	if err := STS.Remove("test"); err != nil {
		t.Error("STS.Remove() must return nil", "err:", err)
		return
	}

	// STS should be empty
	if _, err := STS.Get("test"); err == nil {
		t.Error("STS should be empty")
		return
	}
}

func TestStateTokenStorage_Move(t *testing.T) {
	STS := storage.GetStateTokenStorage()
	STS.Clear()

	// Moving an inexistent key should return an error
	if err := STS.Move("test", "test2"); err == nil {
		t.Error("STS.Move() must return an error if the key does not exist", "STS:", STS)
		return
	}

	STS.Add("test", "test")
	STS.Add("test2", "test2")
	STS.Move("test", "test2")
	// key test should not exist
	if _, err := STS.Get("test"); err == nil {
		t.Error("STS should be empty")
		return
	}
	// key test2 should exist
	value, err := STS.Get("test2")
	if err != nil {
		t.Error("STS should have a value for key test2", "STS:", STS)
		return
	}
	if value != "test" {
		t.Error("STS.Get() must return \"test\"", "value:", value)
		return
	}
}

func TestStateTokenStorage_Clear(t *testing.T) {
	STS := storage.GetStateTokenStorage()
	STS.Clear()
	STS.Add("test", "test")
	STS.Clear()
	// STS should be empty
	if _, err := STS.Get("test"); err == nil {
		t.Error("STS should be empty")
		return
	}
}

func TestStateTokenStorage_String(t *testing.T) {
	STS := storage.GetStateTokenStorage()
	STS.Clear()
	STS.Add("test", "test")
	if STS.String() != "map[test:test]" {
		t.Error("STS.String() must return \"map[test:test]\"", "STS:", STS)
		return
	}
}
