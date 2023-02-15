cp .env.example.json .env.json

change next fields:
APP_PORT
WEBDAV_USER
WEBDAV_PASSWORD
MINIO_ENDPOINT
MINIO_ACCESS_KEY_ID
MINIO_SECRET_ACCESS_KEY
MINIO_BUCKET

go mod tidy for first pull

go run main.go
