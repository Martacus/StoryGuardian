package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Link represents a bidirectional relationship between two entities.
type Link struct {
	ID           string    `json:"id"`
	FromEntityID string    `json:"fromEntityId"`
	ToEntityID   string    `json:"toEntityId"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// LinkService handles CRUD operations for links within a world folder.
// All links are stored in a single links.json file.
type LinkService struct {
	mu sync.RWMutex
}

// linksPath returns the path to links.json for a world folder.
func (s *LinkService) linksPath(folderPath string) string {
	return filepath.Join(folderPath, "links.json")
}

// readLinks reads all links from links.json. Returns an empty slice if the file
// doesn't exist.
func (s *LinkService) readLinks(folderPath string) ([]Link, error) {
	data, err := os.ReadFile(s.linksPath(folderPath))
	if os.IsNotExist(err) {
		return []Link{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read links file: %w", err)
	}

	var links []Link
	if err := json.Unmarshal(data, &links); err != nil {
		return nil, fmt.Errorf("parse links file: %w", err)
	}
	return links, nil
}

// writeLinks writes all links to links.json atomically.
func (s *LinkService) writeLinks(folderPath string, links []Link) error {
	return writeJSONAtomic(s.linksPath(folderPath), links)
}

// ListLinks returns all links in the world.
func (s *LinkService) ListLinks(folderPath string) ([]Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	return s.readLinks(folderPath)
}

// GetLinksForEntity returns all links where the given entity is either the
// fromEntityId or toEntityId.
func (s *LinkService) GetLinksForEntity(folderPath, entityId string) ([]Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if entityId == "" {
		return nil, fmt.Errorf("entity ID cannot be empty")
	}

	all, err := s.readLinks(folderPath)
	if err != nil {
		return nil, err
	}

	var result []Link
	for _, l := range all {
		if l.FromEntityID == entityId || l.ToEntityID == entityId {
			result = append(result, l)
		}
	}
	if result == nil {
		result = []Link{}
	}
	return result, nil
}

// GetLink returns a single link by ID.
func (s *LinkService) GetLink(folderPath, id string) (*Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("link ID cannot be empty")
	}

	all, err := s.readLinks(folderPath)
	if err != nil {
		return nil, err
	}

	for _, l := range all {
		if l.ID == id {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("link not found: %q", id)
}

// CreateLink creates a new link between two entities and appends it to
// links.json. Returns the created link.
func (s *LinkService) CreateLink(folderPath, fromEntityId, toEntityId, linkType string) (*Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if fromEntityId == "" || toEntityId == "" {
		return nil, fmt.Errorf("both entity IDs are required")
	}
	if strings.ContainsAny(fromEntityId, "/\\") || strings.ContainsAny(toEntityId, "/\\") {
		return nil, fmt.Errorf("invalid entity ID")
	}
	if fromEntityId == toEntityId {
		return nil, fmt.Errorf("cannot create a link from an entity to itself")
	}

	all, err := s.readLinks(folderPath)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	link := Link{
		ID:           uuid.New().String(),
		FromEntityID: fromEntityId,
		ToEntityID:   toEntityId,
		Type:         linkType,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	all = append(all, link)
	if err := s.writeLinks(folderPath, all); err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}
	return &link, nil
}

// UpdateLink updates an existing link's type and description. Preserves
// CreatedAt and bumps UpdatedAt.
func (s *LinkService) UpdateLink(folderPath string, link Link) (*Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}
	if link.ID == "" {
		return nil, fmt.Errorf("link ID cannot be empty")
	}

	all, err := s.readLinks(folderPath)
	if err != nil {
		return nil, err
	}

	found := false
	for i, l := range all {
		if l.ID == link.ID {
			link.FromEntityID = l.FromEntityID
			link.ToEntityID = l.ToEntityID
			link.CreatedAt = l.CreatedAt
			link.UpdatedAt = time.Now().UTC()
			all[i] = link
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("link not found: %q", link.ID)
	}

	if err := s.writeLinks(folderPath, all); err != nil {
		return nil, fmt.Errorf("update link: %w", err)
	}
	return &link, nil
}

// DeleteLink removes a link by ID from links.json.
func (s *LinkService) DeleteLink(folderPath, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}
	if id == "" {
		return fmt.Errorf("link ID cannot be empty")
	}

	all, err := s.readLinks(folderPath)
	if err != nil {
		return err
	}

	found := false
	filtered := make([]Link, 0, len(all))
	for _, l := range all {
		if l.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, l)
	}
	if !found {
		return fmt.Errorf("link not found: %q", id)
	}

	if err := s.writeLinks(folderPath, filtered); err != nil {
		return fmt.Errorf("delete link: %w", err)
	}
	return nil
}

// DeleteLinksForEntity removes all links that reference the given entity ID.
// Used for cleanup when an entity is deleted.
func (s *LinkService) DeleteLinksForEntity(folderPath, entityId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}
	if entityId == "" {
		return fmt.Errorf("entity ID cannot be empty")
	}

	all, err := s.readLinks(folderPath)
	if err != nil {
		return err
	}

	filtered := make([]Link, 0, len(all))
	for _, l := range all {
		if l.FromEntityID == entityId || l.ToEntityID == entityId {
			continue
		}
		filtered = append(filtered, l)
	}

	// Only write if something was removed.
	if len(filtered) == len(all) {
		return nil
	}

	if err := s.writeLinks(folderPath, filtered); err != nil {
		return fmt.Errorf("delete links for entity: %w", err)
	}
	return nil
}
