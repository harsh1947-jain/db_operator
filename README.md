# db-operator

A Kubernetes Operator that automates the deployment and management of database instances. Define your database as a Custom Resource, and the operator handles creating and managing the underlying StatefulSets with proper configuration.

Supports **PostgreSQL**, **MongoDB**, and **Redis** out of the box.

## Description

db-operator follows the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) using [Kubebuilder](https://book.kubebuilder.io/introduction.html). It watches for `Dboperator` Custom Resources and reconciles them into fully configured StatefulSets with persistent storage, database-specific environment variables, and proper volume mounts.

Instead of manually crafting StatefulSets, PVCs, and environment configs for each database, you declare what you want:

```yaml
apiVersion: demo.db.com/v1
kind: Dboperator
metadata:
  name: my-postgres
spec:
  dbtype: "your-db-type"
  username: "your-username"
  password: "your-password"
  dbname: "your-dbname"
  version: "enter version"
  storage: "enter storage"
```

The operator creates a StatefulSet with the correct container image, volume mounts, and environment variables for the specified database type.

## Supported Databases

| Database   | Image                          | Data Mount Path                 | Config via Env Vars                                            |
|------------|--------------------------------|---------------------------------|----------------------------------------------------------------|
| PostgreSQL | `postgres:<version>`           | `/var/lib/postgresql/data`      | `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`            |
| MongoDB    | `docker.io/library/mongo:<version>` | `/var/lib/mongodb/data`    | `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD`, `MONGO_INITDB_DATABASE` |
| Redis      | `docker.io/library/redis:<version>` | `/var/lib/redis/data`      | `REDIS_USERNAME`, `REDIS_PASSWORD`, `REDIS_DB`                 |

## Custom Resource Reference

### Spec Fields

| Field      | Required | Description                                      | Example        |
|------------|----------|--------------------------------------------------|----------------|
| `dbtype`   | Yes      | Database type: `postgresql`, `mongodb`, or `redis` | `"postgresql"` |
| `username` | Yes      | Database username                                | `"admin"`      |
| `password` | Yes      | Database password                                | `"secret"`     |
| `dbname`   | No       | Database name                                    | `"myapp"`      |
| `size`     | No       | Small/medium/big                                  | `"small"`        |
| `version`  | No       | Database image tag                               | `"15"`         |

### Status Fields

| Field        | Description                          |
|--------------|--------------------------------------|
| `phase`      | Current phase (e.g., `Creating`, `Running`) |
| `hostname`   | Database service hostname            |
| `port`       | Database listening port              |
| `conditions` | Standard Kubernetes conditions       |

## Getting Started

### Prerequisites

- Go 1.20+
- kubectl with access to a Kubernetes cluster
- Docker or Podman

You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster. Your controller will automatically use the current context in your kubeconfig file.

### Run Locally (Development)

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run the controller (runs in the foreground):

```sh
make run
```

3. Apply a sample resource:

```sh
kubectl apply -f config/samples/demo_v1_dboperator.yaml
```

### Deploy to a Cluster

1. Build and push the container image:

```sh
make docker-build docker-push IMG=<some-registry>/db-operator:tag
```

2. Deploy the controller:

```sh
make deploy IMG=<some-registry>/db-operator:tag
```

3. Apply your database resources:

```sh
kubectl apply -f config/samples/
```

### Cleanup

Remove the controller from the cluster:

```sh
make undeploy
```

Remove the CRDs from the cluster:

```sh
make uninstall
```

## Development

### Build

```sh
make build
```

### Run Tests

```sh
make test
```

### Regenerate Manifests

After editing API definitions in `api/v1/dboperator_types.go`, regenerate CRDs and RBAC manifests:

```sh
make manifests
make generate
```

Run `make --help` for all available targets.

## Project Structure

```
api/v1/                  - CRD type definitions (DboperatorSpec, DboperatorStatus)
cmd/main.go              - Operator entry point
internal/controller/     - Reconciler logic and database defaults
config/crd/              - Generated CRD manifests
config/rbac/             - RBAC roles and bindings
config/manager/          - Controller manager deployment
config/samples/          - Example Dboperator CRs
```

## License

Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
