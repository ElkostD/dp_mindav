package config

import (
	"github.com/gin-gonic/gin"
	. "github.com/totoval/framework/config"
)

func init() {
	webdav := make(map[string]interface{})

	webdav["driver"] = Env("WEBDAV_DRIVER", "memory") // file, memory, minio

	webdav["filesystems"] = map[string]interface{}{
		"file": map[string]interface{}{
			"base_path": ".",
		},
		"minio": map[string]interface{}{
			"endpoint":          Env("MINIO_ENDPOINT", "play.min.io:9000"),
			"bucket":            Env("MINIO_BUCKET", "bucket_name"),
			"access_key_id":     Env("MINIO_ACCESS_KEY_ID", "access_key_id"),
			"secret_access_key": Env("MINIO_SECRET_ACCESS_KEY", "secret_access_key"),
			"use_ssl":           Env("MINIO_USE_SSL", true),
			"location":          "us-east-1",
		},
	}
	webdav["base_path"] = "." // for "file" filesystem
	webdav["supported_folder_depth"] = 10
	webdav["base_url"] = Env("WEBDAV_BASE_URL", "/webdav")
	webdav["memory_upload_mode"] = Env("MEMORY_UPLOAD_MODE", false) // If the host has a large memory, then set to `true` could improve upload performance.

	webdav["accounts"] = gin.Accounts{
		Env("WEBDAV_USER", "totoval").(string): Env("WEBDAV_PASSWORD", "passw0rd").(string),
	}

	Add("webdav", webdav)
}
