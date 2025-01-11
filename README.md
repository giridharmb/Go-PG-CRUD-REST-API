### DB CRUD Operations | REST API : gorm.io / golang-fiber

```bash
docker-compose up -d
```

```bash
go mod init db_crud

go mod tidy

go get -u -v ./...

go mod vendor

go build -o db_crud

./db_crud
```

```bash
docker exec -it metadata_postgres psql -U postgres -d metadata_db
psql (15.10 (Debian 15.10-1.pgdg120+1))
Type "help" for help.

metadata_db=# SELECT * FROM metadata_table;
      my_key      |                                  my_value
------------------+----------------------------------------------------------------------------
 config1          | {"debug_mode": false, "environment": "production", "max_connections": 100}
 user_preferences | {"theme": "dark", "language": "en-US", "notifications": true}
 app_settings     | {"timeout": 30, "cache_ttl": 3600, "retry_attempts": 3}
 feature_flags    | {"new_ui": true, "beta_features": false, "maintenance_mode": false}
 api_config       | {"version": "v2", "rate_limit": 1000, "allowed_origins": ["example.com"]}
(5 rows)
```

```bash
# Connect to the database
docker exec -it metadata_postgres psql -U postgres -d metadata_db

# Once connected, you can run these queries:
# List all records
SELECT * FROM metadata_table;

# Query specific key
SELECT * FROM metadata_table WHERE my_key = 'config1';

# Query JSON content
SELECT my_key, my_value->>'environment' as env 
FROM metadata_table 
WHERE my_key = 'config1';
```

```bash
# Stop containers
docker-compose down

# To also remove the volume (this will delete all data):
docker-compose down -v
```

```bash
# Get config1
curl http://localhost:3000/api/metadata/config1

# Get user preferences
curl http://localhost:3000/api/metadata/user_preferences

# Get app settings
curl http://localhost:3000/api/metadata/app_settings

# Create new configuration
curl -X POST http://localhost:3000/api/metadata \
  -H "Content-Type: application/json" \
  -d '{
    "my_key": "email_settings",
    "my_value": {
      "smtp_host": "smtp.example.com",
      "smtp_port": 587,
      "enable_tls": true
    }
}'

# Update config1 completely
curl -X PUT http://localhost:3000/api/metadata/config1 \
  -H "Content-Type: application/json" \
  -d '{
    "my_value": {
      "environment": "staging",
      "max_connections": 50,
      "debug_mode": true
    }
}'

# Update only debug_mode in config1
curl -X PATCH http://localhost:3000/api/metadata/config1 \
  -H "Content-Type: application/json" \
  -d '{
    "debug_mode": true
}'

# Update theme in user_preferences
curl -X PATCH http://localhost:3000/api/metadata/user_preferences \
  -H "Content-Type: application/json" \
  -d '{
    "theme": "light"
}'

# Upsert (will create if not exists, update if exists)
curl -X PUT http://localhost:3000/api/metadata \
  -H "Content-Type: application/json" \
  -d '{
    "my_key": "cache_settings",
    "my_value": {
      "ttl": 3600,
      "max_size": "1GB"
    }
}'

# Delete feature flags
curl -X DELETE http://localhost:3000/api/metadata/feature_flags

# Delete api config
curl -X DELETE http://localhost:3000/api/metadata/api_config

# Delete all records
curl -X DELETE http://localhost:3000/api/metadata
```

#### To Populate The Table Again

```bash
# Stop containers and remove volumes
docker-compose down -v

# Start containers again
docker-compose up -d

docker exec -it metadata_postgres psql -U postgres -d metadata_db -c "SELECT * FROM metadata_table;"
```
