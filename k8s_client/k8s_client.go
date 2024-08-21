package k8s_client

import (
	"k8s.io/client-go/rest"
)

const (
	namespace = "default"

	postgresImage   = "postgres:16.3-alpine3.20"
	postgresStorage = "128Mi"
	postgresPort    = 5432

	postgresUser     = "username"
	postgresPassword = "password"
	postgresDB       = "kaas"
)

var clientConfig *rest.Config

func Initialize() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	} else {
		clientConfig = config
	}
}
