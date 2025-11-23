package e2e

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

var (
	serverPort     string
	postgresContainer testcontainers.Container
	serverProcess  *exec.Cmd
	serverDSN      string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Убираем смайлики из testcontainers
	os.Setenv("TESTCONTAINERS_LOG_LEVEL", "WARN")

	// Find free port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Printf("Failed to find free port: %v\n", err)
		os.Exit(1)
	}
	serverPort = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	baseURLStr := "http://localhost:" + serverPort
	listener.Close()

	// Инициализируем baseURL в requests.go
	setBaseURL(baseURLStr)

	// Start PostgreSQL
	container, dsn, err := startPostgres(ctx)
	if err != nil {
		fmt.Printf("Failed to start postgres: %v\n", err)
		os.Exit(1)
	}
	postgresContainer = container
	serverDSN = dsn

	// Run migrations
	if err := runMigrations(dsn); err != nil {
		fmt.Printf("Failed to run migrations: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	// Start server
	if err := startServer(dsn); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		cleanup()
		os.Exit(1)
	}

	// Wait for server to be ready
	time.Sleep(3 * time.Second)
	if !waitForServer() {
		fmt.Println("Server failed to start")
		cleanup()
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}

func startServer(dsn string) error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return err
	}

	port := extractPort(dsn)
	host := "localhost"

	// Set environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("STORAGES_POSTGRES_HOST=%s", host))
	env = append(env, fmt.Sprintf("STORAGES_POSTGRES_PORT=%s", port))
	env = append(env, "STORAGES_POSTGRES_USER=testuser")
	env = append(env, "STORAGES_POSTGRES_PASS=testpass")
	env = append(env, "STORAGES_POSTGRES_NAME=testdb")
	env = append(env, "STORAGES_POSTGRES_SSLM=disable")
	env = append(env, "SERVERS_REST_ADDR=0.0.0.0")
	env = append(env, fmt.Sprintf("SERVERS_REST_PORT=%s", serverPort))
	env = append(env, "LOGGER_LEVEL=info")
	env = append(env, "LOGGER_ENCODING=json")
	env = append(env, "BUSSINES_LOGIC_ALLOWED_REUSE_TO_REASIGN=true")
	env = append(env, "BUSSINES_LOGIC_ALLOWE_STATUSES_TO_REASIGN=active")
	env = append(env, "BUSSINES_LOGIC_ALLOWED_ROLES_TO_REASIGN=default")
	env = append(env, "SERVERS_REST_READ_TIMEOUT=5s")
	env = append(env, "SERVERS_REST_WRITE_TIMEOUT=5s")
	env = append(env, "SERVERS_REST_READ_HEADER_TIMEOUT=5s")
	env = append(env, "SERVERS_REST_IDLE_TIMEOUT=5s")
	env = append(env, "SERVERS_REST_HEALTH_CHECK_ROUTE=health")

	serverProcess = exec.Command("go", "run", "./cmd/pr-reviewer-service/main.go")
	serverProcess.Dir = projectRoot
	serverProcess.Env = env
	serverProcess.Stdout = os.Stdout
	serverProcess.Stderr = os.Stderr

	return serverProcess.Start()
}

func waitForServer() bool {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(getBaseURL() + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func cleanup() {
	if serverProcess != nil && serverProcess.Process != nil {
		// Try graceful shutdown first
		if err := serverProcess.Process.Signal(os.Interrupt); err == nil {
			done := make(chan error, 1)
			go func() {
				done <- serverProcess.Wait()
			}()
			select {
			case <-done:
			case <-time.After(2 * time.Second):
				// Force kill if graceful shutdown didn't work
				_ = serverProcess.Process.Kill()
				<-done
			}
		} else {
			// If signal failed, just kill
			_ = serverProcess.Process.Kill()
			_ = serverProcess.Wait()
		}
		serverProcess = nil
	}

	if postgresContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = postgresContainer.Terminate(ctx)
		cancel()
		postgresContainer = nil
	}
}

