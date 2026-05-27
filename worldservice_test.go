package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateWorldUsesSelectedEmptyFolder(t *testing.T) {
	root := t.TempDir()
	worldDir := filepath.Join(root, "chosen")
	if err := os.Mkdir(worldDir, 0755); err != nil {
		t.Fatalf("mkdir selected folder: %v", err)
	}

	info, err := (&WorldService{}).CreateWorld("Emberfall", worldDir)
	if err != nil {
		t.Fatalf("CreateWorld returned error: %v", err)
	}

	if info.Path != worldDir {
		t.Fatalf("expected world path %q, got %q", worldDir, info.Path)
	}
	if _, err := os.Stat(filepath.Join(worldDir, worldMetaFileName)); err != nil {
		t.Fatalf("expected world.json in selected folder: %v", err)
	}
}

func TestCreateWorldCreatesChildFolderWhenSelectedFolderIsNotEmpty(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "library")
	if err := os.Mkdir(parent, 0755); err != nil {
		t.Fatalf("mkdir selected folder: %v", err)
	}
	if err := os.WriteFile(filepath.Join(parent, "notes.txt"), []byte("occupied"), 0644); err != nil {
		t.Fatalf("seed selected folder: %v", err)
	}

	info, err := (&WorldService{}).CreateWorld("Emberfall", parent)
	if err != nil {
		t.Fatalf("CreateWorld returned error: %v", err)
	}

	want := filepath.Join(parent, "Emberfall")
	if info.Path != want {
		t.Fatalf("expected world path %q, got %q", want, info.Path)
	}
	if _, err := os.Stat(filepath.Join(want, worldMetaFileName)); err != nil {
		t.Fatalf("expected world.json in child folder: %v", err)
	}
}

func TestCreateWorldErrorsWhenChildWorldFolderAlreadyExists(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "library")
	if err := os.Mkdir(parent, 0755); err != nil {
		t.Fatalf("mkdir selected folder: %v", err)
	}
	if err := os.WriteFile(filepath.Join(parent, "notes.txt"), []byte("occupied"), 0644); err != nil {
		t.Fatalf("seed selected folder: %v", err)
	}
	if err := os.Mkdir(filepath.Join(parent, "Emberfall"), 0755); err != nil {
		t.Fatalf("mkdir existing child folder: %v", err)
	}

	if _, err := (&WorldService{}).CreateWorld("Emberfall", parent); err == nil {
		t.Fatal("expected error for existing child world folder, got nil")
	}
}
