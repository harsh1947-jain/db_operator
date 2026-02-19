package controller

import (
	// "k8s.io/apimachinery/pkg/api/resource"
	demov1 "github.com/harsh1947-jain/db-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type DbConfig struct {
	NameofDB     string
	ImageofDB    string
	Mountingpath string
	UsernameEnv  string
	PasswordEnv  string
	DbEnv        string
	CPU          resource.Quantity
	Memory       resource.Quantity
	Storage      resource.Quantity
}

func getDbDefaults(db *demov1.Dboperator) DbConfig {
	var config DbConfig

	// Resource Mapping
	var cpu, mem, storage string
	switch db.Spec.Size {
	case "medium":
		cpu, mem, storage = "1", "2Gi", "2Gi"
	case "big":
		cpu, mem, storage = "2", "3Gi", "3Gi"
	default: // "small"
		cpu, mem, storage = "250m", "256Mi", "1Gi"
	}

	config.CPU = resource.MustParse(cpu)
	config.Memory = resource.MustParse(mem)
	config.Storage = resource.MustParse(storage)

	// Engine Mapping
	switch db.Spec.DbType {
	case "postgresql":
		config.NameofDB, config.Mountingpath = "postgresql", "/var/lib/postgresql/data"
		config.UsernameEnv, config.PasswordEnv, config.DbEnv = "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"
	case "mongodb":
		config.NameofDB, config.Mountingpath = "mongodb", "/var/lib/mongodb/data"
		config.UsernameEnv, config.PasswordEnv, config.DbEnv = "MONGO_INITDB_ROOT_USERNAME", "MONGO_INITDB_ROOT_PASSWORD", "MONGO_INITDB_DATABASE"
	case "redis":
		config.NameofDB, config.Mountingpath = "redis", "/var/lib/redis/data"
		config.UsernameEnv, config.PasswordEnv, config.DbEnv = "REDIS_USERNAME", "REDIS_PASSWORD", "REDIS_DB"
	}
	config.ImageofDB = "mirror.gcr.io/library/" + config.NameofDB + ":" + db.Spec.Version

	return config
}
