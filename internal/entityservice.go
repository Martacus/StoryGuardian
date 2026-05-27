package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Entity represents a single worldbuilding entity stored as entities/{id}.json.
type Entity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Categories  []string  `json:"categories"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// EntityService handles CRUD operations for entities within a world folder.
type EntityService struct {
	mu sync.RWMutex
}

// entitiesDir returns the entities/ subdirectory path for a world folder.
func (s *EntityService) entitiesDir(folderPath string) string {
	return filepath.Join(folderPath, "entities")
}

// entityPath returns the path for a single entity JSON file.
func (s *EntityService) entityPath(folderPath, id string) string {
	return filepath.Join(s.entitiesDir(folderPath), id+".json")
}

// ListEntities reads all entity JSON files from entities/ and returns them
// sorted by name. Returns an empty slice (not an error) if the directory
// doesn't exist yet.
func (s *EntityService) ListEntities(folderPath string) ([]Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	dir := s.entitiesDir(folderPath)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return []Entity{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read entities dir: %w", err)
	}

	var entities []Entity
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		var e Entity
		if err := json.Unmarshal(data, &e); err != nil {
			continue
		}
		entities = append(entities, e)
	}

	sort.Slice(entities, func(i, j int) bool {
		return strings.ToLower(entities[i].Name) < strings.ToLower(entities[j].Name)
	})

	return entities, nil
}

// GetEntity reads a single entity by ID from entities/{id}.json.
func (s *EntityService) GetEntity(folderPath, id string) (*Entity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}
	if strings.ContainsAny(id, "/\\") {
		return nil, fmt.Errorf("invalid entity ID: %q", id)
	}

	data, err := os.ReadFile(s.entityPath(folderPath, id))
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("entity not found: %q", id)
	}
	if err != nil {
		return nil, fmt.Errorf("read entity: %w", err)
	}

	var e Entity
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("parse entity: %w", err)
	}
	return &e, nil
}

// CreateEntity creates a new entity with the given name and type, generates a
// UUID, and writes it atomically to entities/{id}.json. Returns the created entity.
func (s *EntityService) CreateEntity(folderPath, name, entityType string) (*Entity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("entity name cannot be empty")
	}

	now := time.Now().UTC()
	e := Entity{
		ID:         uuid.New().String(),
		Name:       name,
		Type:       entityType,
		Categories: []string{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := writeJSONAtomic(s.entityPath(folderPath, e.ID), e); err != nil {
		return nil, fmt.Errorf("create entity: %w", err)
	}
	return &e, nil
}

// UpdateEntity overwrites an existing entity on disk. It preserves CreatedAt
// from the existing file and bumps UpdatedAt. Returns the updated entity.
func (s *EntityService) UpdateEntity(folderPath string, entity Entity) (*Entity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if entity.ID == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}
	if strings.ContainsAny(entity.ID, "/\\") {
		return nil, fmt.Errorf("invalid entity ID: %q", entity.ID)
	}
	if entity.Name == "" {
		return nil, fmt.Errorf("entity name cannot be empty")
	}

	// Read existing to preserve CreatedAt.
	path := s.entityPath(folderPath, entity.ID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("entity not found: %q", entity.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("read entity: %w", err)
	}

	var existing Entity
	if err := json.Unmarshal(data, &existing); err != nil {
		return nil, fmt.Errorf("parse entity: %w", err)
	}

	entity.CreatedAt = existing.CreatedAt
	entity.UpdatedAt = time.Now().UTC()

	if entity.Categories == nil {
		entity.Categories = []string{}
	}

	if err := writeJSONAtomic(path, entity); err != nil {
		return nil, fmt.Errorf("update entity: %w", err)
	}
	return &entity, nil
}

// DeleteEntity removes an entity file from disk.
func (s *EntityService) DeleteEntity(folderPath, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("entity ID cannot be empty")
	}
	if strings.ContainsAny(id, "/\\") {
		return fmt.Errorf("invalid entity ID: %q", id)
	}

	path := s.entityPath(folderPath, id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("entity not found: %q", id)
		}
		return fmt.Errorf("delete entity: %w", err)
	}
	return nil
}
