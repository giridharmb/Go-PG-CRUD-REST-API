package main

import (
	"encoding/json"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MetadataEntry struct {
	MyKey   string          `gorm:"column:my_key;primaryKey" json:"my_key"`
	MyValue json.RawMessage `gorm:"column:my_value;type:jsonb" json:"my_value"`
}

// Helper method to get the value as a map
func (m *MetadataEntry) GetValueAsMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(m.MyValue, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (MetadataEntry) TableName() string {
	return "metadata_table"
}

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=metadata_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&MetadataEntry{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
