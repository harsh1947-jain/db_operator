package controller

import (
	demov1 "github.com/harsh1947-jain/db-operator/api/v1"
)

// DbConfig holds the mapping for different database types
type DbConfig struct {
	NameofDB     string
	ImageofDB    string
	Mountingpath string
	usernameenv  string
	passwordenv  string
	dbenv        string
}

func getDbDefaults(db *demov1.Dboperator) DbConfig {
	var config DbConfig

	if db.Spec.DbType == "postgresql" {
		config.NameofDB = "postgresql"
		config.ImageofDB = "postgres:" + db.Spec.Version
		config.Mountingpath = "/var/lib/postgresql/data"
		config.usernameenv = "POSTGRES_USER"
		config.passwordenv = "POSTGRES_PASSWORD"
		config.dbenv = "POSTGRES_DB"
	} else if db.Spec.DbType == "mongodb" {
		config.NameofDB = "mongodb"
		config.ImageofDB = "docker.io/library/mongo:" + db.Spec.Version
		config.Mountingpath = "/var/lib/mongodb/data"
		config.usernameenv = "MONGO_INITDB_ROOT_USERNAME"
		config.passwordenv = "MONGO_INITDB_ROOT_PASSWORD"
		config.dbenv = "MONGO_INITDB_DATABASE"
	} else if db.Spec.DbType == "redis" {
		config.NameofDB = "redis"
		config.ImageofDB = "docker.io/library/redis:" + db.Spec.Version
		config.Mountingpath = "/var/lib/redis/data"
		config.usernameenv = "REDIS_USERNAME"
		config.passwordenv = "REDIS_PASSWORD"
		config.dbenv = "REDIS_DB"
	}

	return config
}
