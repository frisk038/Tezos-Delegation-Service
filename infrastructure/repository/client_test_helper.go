package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testContainer represents a Docker container used for PostgreSQL testing.
type testContainer struct {
	testcontainers.Container
	uri string
}

// initContainer initializes and starts a PostgreSQL Docker container for testing.
func initContainer(t *testing.T) *testContainer {
	ctx := context.Background()
	var env = map[string]string{
		"POSTGRES_PASSWORD": "test",
		"POSTGRES_USER":     "test",
		"POSTGRES_DB":       "testdb",
	}

	// Create and start a PostgreSQL container.
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres",
			ExposedPorts: []string{"5432/tcp"},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	time.Sleep(time.Second)

	return &testContainer{
		Container: container,
		uri: fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			env["POSTGRES_USER"], env["POSTGRES_PASSWORD"], mappedPort.Port(), env["POSTGRES_DB"])}
}

// migrateDb performs database migration using Golang Migrate.
func migrateDb(t *testing.T, tc *testContainer) {
	m, err := migrate.New(
		"file://../../migration/deploy",
		tc.uri)
	require.NoError(t, err)

	err = m.Up()
	require.NoError(t, err)
}

// clears all data from table.
func clearTable(ctx context.Context, t *testing.T, conn *pgxpool.Pool) {
	_, err := conn.Exec(ctx, "TRUNCATE delegations")
	require.NoError(t, err)
}
