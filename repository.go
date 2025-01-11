package main

import (
	"encoding/json"
	"gorm.io/gorm"
)

type MetadataRepository interface {
	Create(entry *MetadataEntry) error
	Get(key string) (*MetadataEntry, error)
	Update(entry *MetadataEntry) error
	PatchUpdate(key string, partialValue map[string]any) error
	Delete(key string) error
	DeleteAll() error
	Upsert(entry *MetadataEntry) error
}

type metadataRepository struct {
	db *gorm.DB
}

func NewMetadataRepository(db *gorm.DB) MetadataRepository {
	return &metadataRepository{db: db}
}

func (r *metadataRepository) Create(entry *MetadataEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(entry).Error
	})
}

func (r *metadataRepository) Get(key string) (*MetadataEntry, error) {
	var entry MetadataEntry
	err := r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Where("my_key = ?", key).First(&entry).Error
	})
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *metadataRepository) Update(entry *MetadataEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&MetadataEntry{}).Where("my_key = ?", entry.MyKey).Update("my_value", entry.MyValue).Error
	})
}

func (r *metadataRepository) PatchUpdate(key string, partialValue map[string]any) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing MetadataEntry
		if err := tx.Where("my_key = ?", key).First(&existing).Error; err != nil {
			return err
		}

		// Get existing value as map
		existingMap, err := existing.GetValueAsMap()
		if err != nil {
			return err
		}

		// Merge the partial update with existing value
		for k, v := range partialValue {
			existingMap[k] = v
		}

		// Convert back to JSON
		updatedJSON, err := json.Marshal(existingMap)
		if err != nil {
			return err
		}

		return tx.Model(&MetadataEntry{}).Where("my_key = ?", existing.MyKey).Update("my_value", updatedJSON).Error
	})
}

func (r *metadataRepository) Delete(key string) error {
	var result *gorm.DB
	err := r.db.Transaction(func(tx *gorm.DB) error {
		result = tx.Where("my_key = ?", key).Delete(&MetadataEntry{})
		return result.Error
	})

	if err != nil {
		return err
	}

	// If no rows were affected, return record not found error
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *metadataRepository) DeleteAll() error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&MetadataEntry{}).Error
	})
}

func (r *metadataRepository) Upsert(entry *MetadataEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Save(entry).Error
	})
}
