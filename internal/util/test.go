package util

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/ttagiyeva/entain/internal/config"
	"github.com/ttagiyeva/entain/internal/database"
)

func CreateTestContainer(ctx context.Context, s *suite.Suite) *database.Postgres {
	cfg := config.Config{
		DB: config.DB{
			Host:     "localhost",
			Port:     5432,
			User:     "root",
			Password: "root",
			Name:     "postgres",
		},
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DB.User,
			"POSTGRES_PASSWORD": cfg.DB.Password,
			"POSTGRES_DB":       cfg.DB.Name,
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	endpoint, err := container.Endpoint(ctx, "")
	s.Require().NoError(err)

	portStr := strings.Split(endpoint, ":")[1]

	port, err := strconv.Atoi(portStr)
	s.Require().NoError(err)

	cfg.DB.Port = uint16(port)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in database connection", r)
		}
	}()

	db := database.NewPostgres()

	// 10 iterations to wait for the database to be ready.
	for i := 0; i < 10; i++ {
		err = db.Connect(ctx, &cfg)
		if err != nil {
			time.Sleep(time.Millisecond * 500)

			continue
		}

		s.Require().NoError(err)

		break
	}

	return db
}
