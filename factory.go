package main

import (
	"gorm.io/gorm"
)

type MetadataFactory struct {
	db *gorm.DB
}

func NewMetadataFactory(db *gorm.DB) *MetadataFactory {
	return &MetadataFactory{db: db}
}

func (f *MetadataFactory) CreateRepository() MetadataRepository {
	return NewMetadataRepository(f.db)
}

func (f *MetadataFactory) CreateService() *MetadataService {
	return NewMetadataService(f.CreateRepository())
}
