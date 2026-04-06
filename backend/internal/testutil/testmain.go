package testutil

import (
	"context"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testContainer     *postgres.PostgresContainer
	testDSN           string
	containerInitOnce sync.Once
	containerInitErr  error
)

func GetTestContainer(ctx context.Context) (string, error) {
	containerInitOnce.Do(func() {
		if runtime.GOOS == "windows" && os.Getenv("DOCKER_HOST") == "" {
			os.Setenv("DOCKER_HOST", "npipe:////./pipe/docker_engine")
		}

		container, err := postgres.Run(
			ctx,
			"postgres:16-alpine",
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("test"),
			postgres.WithPassword("test"),
			postgres.BasicWaitStrategies(),
		)
		if err != nil {
			containerInitErr = err
			return
		}

		testContainer = container
		testDSN, containerInitErr = container.ConnectionString(ctx, "sslmode=disable")
		if containerInitErr != nil {
			return
		}

		log.Printf("Test container started: %s", testDSN)
	})

	return testDSN, containerInitErr
}
