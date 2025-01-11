package main

import (
	"fmt"
	"regexp"
)

var (
	// Regex pattern for my_key validation
	// Allows alphanumeric characters, underscore, hyphen, and dot
	keyPattern = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
)

type MetadataService struct {
	repo MetadataRepository
}

func NewMetadataService(repo MetadataRepository) *MetadataService {
	return &MetadataService{repo: repo}
}

func (s *MetadataService) validateKey(key string) error {
	if !keyPattern.MatchString(key) {
		return fmt.Errorf("invalid key format: key can only contain alphanumeric characters, underscore, hyphen, and dot")
	}
	return nil
}

func (s *MetadataService) Create(entry *MetadataEntry) error {
	if err := s.validateKey(entry.MyKey); err != nil {
		return err
	}
	return s.repo.Create(entry)
}

func (s *MetadataService) Get(key string) (*MetadataEntry, error) {
	if err := s.validateKey(key); err != nil {
		return nil, err
	}
	return s.repo.Get(key)
}

func (s *MetadataService) Update(entry *MetadataEntry) error {
	if err := s.validateKey(entry.MyKey); err != nil {
		return err
	}
	return s.repo.Update(entry)
}

func (s *MetadataService) PatchUpdate(key string, partialValue map[string]any) error {
	if err := s.validateKey(key); err != nil {
		return err
	}
	return s.repo.PatchUpdate(key, partialValue)
}

func (s *MetadataService) Delete(key string) error {
	if err := s.validateKey(key); err != nil {
		return err
	}
	return s.repo.Delete(key)
}

func (s *MetadataService) DeleteAll() error {
	return s.repo.DeleteAll()
}

func (s *MetadataService) Upsert(entry *MetadataEntry) error {
	if err := s.validateKey(entry.MyKey); err != nil {
		return err
	}
	return s.repo.Upsert(entry)
}
