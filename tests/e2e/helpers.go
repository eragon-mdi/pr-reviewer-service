package e2e

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUser     = "testuser"
	dbPassword = "testpass"
	dbName     = "testdb"
	dbPort     = "5432"
)

func startPostgres(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.5",
		ExposedPorts: []string{dbPort + "/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPassword,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(dbPort+"/tcp"),
		).WithDeadline(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start postgres: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, "", fmt.Errorf("failed to get host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, dbPort)
	if err != nil {
		container.Terminate(ctx)
		return nil, "", fmt.Errorf("failed to get mapped port: %w", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, host, mappedPort.Port(), dbName)

	return container, dsn, nil
}

func runMigrations(dsn string) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	migratePath := filepath.Join(projectRoot, "migrations", "sql")
	migrateCmd := exec.Command("migrate",
		"-path", migratePath,
		"-database", dsn,
		"up",
	)

	migrateCmd.Dir = projectRoot
	output, err := migrateCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("migration failed: %w, output: %s", err, string(output))
	}

	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found")
		}
		dir = parent
	}
}

func extractPort(dsn string) string {
	// dsn format: postgres://user:pass@host:port/db
	parts := strings.Split(dsn, "@")
	if len(parts) != 2 {
		return "5432"
	}
	hostPort := strings.Split(parts[1], "/")[0]
	portParts := strings.Split(hostPort, ":")
	if len(portParts) == 2 {
		return portParts[1]
	}
	return "5432"
}
