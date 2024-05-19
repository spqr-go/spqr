package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Go project with predefined scaffolding",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter project name: ")
		projectName, _ := reader.ReadString('\n')
		projectName = strings.TrimSpace(projectName)

		fmt.Print("Select database (1 for Postgres 🐘, 2 for Mariadb 🦭): ")
		dbOption, _ := reader.ReadString('\n')
		dbOption = strings.TrimSpace(dbOption)

		var db string
		switch dbOption {
		case "1":
			db = "Postgres 🐘"
		case "2":
			db = "Mariadb 🦭"
		default:
			fmt.Println("Unsupported database option")
			return
		}

		fmt.Print("Enter DB_USER: ")
		dbUser, _ := reader.ReadString('\n')
		dbUser = strings.TrimSpace(dbUser)

		fmt.Print("Enter DB_PASSWORD: ")
		dbPassword, _ := reader.ReadString('\n')
		dbPassword = strings.TrimSpace(dbPassword)

		fmt.Print("Enter DB_NAME: ")
		dbName, _ := reader.ReadString('\n')
		dbName = strings.TrimSpace(dbName)

		dbHost := strings.TrimSpace("localhost")

		fmt.Print("Enter DB_PORT: ")
		dbPort, _ := reader.ReadString('\n')
		dbPort = strings.TrimSpace(dbPort)

		createProjectScaffolding(projectName, db, dbUser, dbPassword, dbName, dbHost, dbPort)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
func createProjectScaffolding(projectName string, db string, dbUser string, dbPassword string, dbName string, dbHost string, dbPort string) {

	// Create project directory
	err := os.Mkdir(projectName, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating project directory %s: %v\n", projectName, err)
		return
	}

	// Initialize Go module
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectName
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error initializing Go module: %v\n", err)
		return
	}

	// Install Gin as a dependency
	cmd = exec.Command("go", "get", "github.com/gin-gonic/gin")
	cmd.Dir = projectName
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error installing Gin: %v\n", err)
		return
	}

	baseDirs := []string{
		"cmd/api",
		"internal/adapters/in",
		"internal/adapters/in/models/request",
		"internal/adapters/in/models/response",
		"internal/adapters/out/repositories",
		"internal/configs",
		"internal/core/auth",
		"internal/core/domain",
		"internal/core/ports/in",
		"internal/core/ports/out",
		"internal/core/usecases",
		"internal/routing",
	}

	for _, dir := range baseDirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", path, err)
			return
		}
	}

	files := map[string]string{
		filepath.Join(projectName, "cmd/api/spqr.go"):                                        "package main\n\nfunc main() {\n\t// TODO: Implement\n}\n",
		filepath.Join(projectName, "internal/configs/config.go"):                             "package configs\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/core/domain/domain.go"):                         "package domain\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/adapters/in/models/request/request_model.go"):   "package request\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/adapters/in/models/response/response_model.go"): "package response\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/adapters/in/adapter_in.go"):                     "package in\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/adapters/out/repositories/repository.go"):       "package repositories\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/core/auth/auth.go"):                             "package auth\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/core/usecases/usecases.go"):                     "package usecases\n\n// TODO: Implement\n",
		filepath.Join(projectName, "internal/routing/router.go"):                             "package routing\n\n// TODO: Implement\n",
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("Error creating file %s: %v\n", path, err)
			return
		}
	}

	// Create docker-compose file
	var dbImage string
	var dbEnv string
	switch db {
	case "Postgres 🐘":
		dbImage = "postgres:latest"
		dbEnv = fmt.Sprintf("POSTGRES_PASSWORD: %s\n      POSTGRES_DB: %s\n      POSTGRES_USER: %s", dbPassword, dbName, dbUser)
	case "Mariadb 🦭":
		dbImage = "mariadb:latest"
		dbEnv = fmt.Sprintf("MYSQL_ROOT_PASSWORD: %s\n      MYSQL_DATABASE: %s\n      MYSQL_USER: %s\n      MYSQL_PASSWORD: %s", dbPassword, dbName, dbUser, dbPassword)
	default:
		fmt.Printf("Unsupported database: %s\n", db)
		return
	}

	dockerComposeContent := fmt.Sprintf(`version: '3.8'
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: builder
    ports:
      - "8080:8080"
    depends_on:
      db:
        condition: service_healthy
    environment:
      DB_USER: %s
      DB_PASSWORD: %s
      DB_NAME: %s
      DB_HOST: %s
      DB_PORT: %s
  db:
    image: %s
    environment:
      %s
    volumes:
      - db-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U %s"]
      timeout: 20s
      retries: 10
volumes:
  db-data:
`, dbUser, dbPassword, dbName, dbHost, dbPort, dbImage, dbEnv, dbUser)

	dockerComposePath := filepath.Join(projectName, "docker-compose.yml")
	if err := os.WriteFile(dockerComposePath, []byte(dockerComposeContent), 0644); err != nil {
		fmt.Printf("Error creating docker-compose file: %v\n", err)
		return
	}

	// Copy Dockerfile
	dockerfileSrc := "cmd/Dockerfile"
	dockerfileDst := filepath.Join(projectName, "Dockerfile")
	dockerfileContent, err := ioutil.ReadFile(dockerfileSrc)
	if err != nil {
		fmt.Printf("Error reading Dockerfile: %v\n", err)
		return
	}
	if err := os.WriteFile(dockerfileDst, dockerfileContent, 0644); err != nil {
		fmt.Printf("Error writing Dockerfile: %v\n", err)
		return
	}

	fmt.Printf("Project %s created successfully!\n", projectName)
}
