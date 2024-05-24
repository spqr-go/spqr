package cmd

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// coffeeCmd represents the coffee command
var coffeeCmd = &cobra.Command{
	Use:   "coffee",
	Short: "Creates with SPQR a Spring Boot  project with Gradle and hexagonal architecture",
	Long:  `This command creates a new Spring Boot project with Gradle and hexagonal architecture.`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Lest's get started with your Spring Boot 🍃 project using SPQR 🦅 ! 🚀\n\n")
		fmt.Print("Enter project name: ")
		projectName, _ := reader.ReadString('\n')
		projectName = strings.TrimSpace(projectName)

		fmt.Print("Enter package name (e.g., com.example.project): ")
		packageName, _ := reader.ReadString('\n')
		packageName = strings.TrimSpace(packageName)

		// Valores por defecto para la base de datos
		dbUser := "default_user"
		dbPassword := "default_password"
		dbName := "default_db"
		dbHost := "localhost"
		dbPort := "5432"

		createProjectScaffoldingCoffee(projectName, packageName, dbUser, dbPassword, dbName, dbHost, dbPort)
	},
}

func createProjectScaffoldingCoffee(projectName, packageName, dbUser, dbPassword, dbName, dbHost, dbPort string) {
	baseDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	projectPath := filepath.Join(baseDir, projectName)
	err = os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating project directory %s: %v\n", projectPath, err)
		return
	}

	url := fmt.Sprintf("https://start.spring.io/starter.zip?type=gradle-project&language=java&bootVersion=3.3.0&baseDir=%s&groupId=%s&artifactId=%s&name=%s&description=%s&packageName=%s&packaging=jar&javaVersion=17&dependencies=web,data-jpa,postgresql", projectName, packageName, projectName, projectName, projectName, packageName)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading project template: %v\n", err)
		return
	}
	defer resp.Body.Close()

	zipPath := filepath.Join(projectPath, "project.zip")
	out, err := os.Create(zipPath)
	if err != nil {
		fmt.Printf("Error creating zip file: %v\n", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to zip file: %v\n", err)
		return
	}

	err = unzip(zipPath, projectPath)
	if err != nil {
		fmt.Printf("Error unzipping project: %v\n", err)
		return
	}

	err = os.Remove(zipPath)
	if err != nil {
		fmt.Printf("Error removing zip file: %v\n", err)
		return
	}

	// Moverse dentro del directorio del proyecto generado
	projectInnerPath := filepath.Join(projectPath, projectName)

	// Crear directorio src/main/resources si no existe
	resourcesDir := filepath.Join(projectInnerPath, "src", "main", "resources")
	err = os.MkdirAll(resourcesDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating resources directory %s: %v\n", resourcesDir, err)
		return
	}

	applicationPropertiesPath := filepath.Join(resourcesDir, "application.properties")
	applicationPropertiesContent := fmt.Sprintf(`spring.datasource.url=jdbc:postgresql://%s:%s/%s
spring.datasource.username=%s
spring.datasource.password=%s
spring.jpa.hibernate.ddl-auto=update`, dbHost, dbPort, dbName, dbUser, dbPassword)
	err = os.WriteFile(applicationPropertiesPath, []byte(applicationPropertiesContent), 0644)
	if err != nil {
		fmt.Printf("Error creating application.properties file: %v\n", err)
		return
	}

	sourceDir := filepath.Join(projectInnerPath, "src", "main", "java", strings.ReplaceAll(packageName, ".", "/"))

	// Definir estructura de directorios
	dirs := []string{
		filepath.Join(sourceDir, "adapters", "in", "consumer"),
		filepath.Join(sourceDir, "adapters", "in", "controllers", "models"),
		filepath.Join(sourceDir, "adapters", "out", "jdbc"),
		filepath.Join(sourceDir, "adapters", "out", "jpa", "models"),
		filepath.Join(sourceDir, "adapters", "out", "jpa", "repository"),
		filepath.Join(sourceDir, "adapters", "producer"),
		filepath.Join(sourceDir, "adapters", "rest"),
		filepath.Join(sourceDir, "configs", "handlers"),
		filepath.Join(sourceDir, "configs", "properties"),
		filepath.Join(sourceDir, "core", "domain", "dtos"),
		filepath.Join(sourceDir, "core", "domain", "enums"),
		filepath.Join(sourceDir, "core", "domain", "exceptions"),
		filepath.Join(sourceDir, "core", "domain", "models"),
		filepath.Join(sourceDir, "core", "domain", "parser"),
		filepath.Join(sourceDir, "core", "ports", "in"),
		filepath.Join(sourceDir, "core", "ports", "out"),
		filepath.Join(sourceDir, "core", "usecases"),
		filepath.Join(projectInnerPath, "src", "main", "resources", "db", "migration"),
		filepath.Join(projectInnerPath, "src", "test", "java"),
	}

	// Crear directorios
	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Crear archivo build.gradle con configuración y dependencias necesarias
	buildGradleContent := `plugins {
    id 'java'
    id 'org.springframework.boot' version '3.1.6'
    id 'io.spring.dependency-management' version '1.1.4'
    id 'jacoco'
    id 'org.barfuin.gradle.jacocolog' version '3.1.0'
    id "io.sentry.jvm.gradle" version "3.14.0"
}

group = '` + packageName + `'
version = '0.0.1-SNAPSHOT'

java {
    sourceCompatibility = '17'
}

configurations {
    compileOnly {
        extendsFrom annotationProcessor
    }
}

repositories {
    mavenCentral()
}

ext {
    set('springCloudVersion', "2022.0.4")
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-data-jpa'
    implementation 'org.springframework.boot:spring-boot-starter-web'
    compileOnly 'org.projectlombok:lombok'
    annotationProcessor 'org.projectlombok:lombok'
    testImplementation 'org.springframework.boot:spring-boot-starter-test'
    runtimeOnly 'org.postgresql:postgresql'
}

dependencyManagement {
    imports {
        mavenBom "org.springframework.cloud:spring-cloud-dependencies:${springCloudVersion}"
    }
}

tasks.named('test') {
    useJUnitPlatform()
}

test {
    finalizedBy jacocoTestReport
    finalizedBy jacocoTestCoverageVerification
}

jacocoTestReport {
    dependsOn test
    reports {
        xml.required = true
    }
    afterEvaluate {
        classDirectories.setFrom(files(classDirectories.files.collect {
            fileTree(dir: it, exclude: ["**/configs/**", "**/*Application*", "**/exceptions/**"])
        }))
    }
}

jacocoTestCoverageVerification {
    dependsOn jacocoTestReport
    violationRules {
        rule {
            limit {
                minimum = 0.80
            }
        }
    }
}`

	buildGradlePath := filepath.Join(projectInnerPath, "build.gradle")
	err = os.WriteFile(buildGradlePath, []byte(buildGradleContent), 0644)
	if err != nil {
		fmt.Printf("Error creating build.gradle file: %v\n", err)
		return
	}

	// Crear archivo settings.gradle
	settingsGradleContent := `rootProject.name = '` + projectName + `'`
	settingsGradlePath := filepath.Join(projectInnerPath, "settings.gradle")
	err = os.WriteFile(settingsGradlePath, []byte(settingsGradleContent), 0644)
	if err != nil {
		fmt.Printf("Error creating settings.gradle file: %v\n", err)
		return
	}

	// Crear Dockerfile
	dockerfileContent := `FROM openjdk:17-jdk-slim
VOLUME /tmp
COPY build/libs/*.jar app.jar
ENTRYPOINT ["java","-Djava.security.egd=file:/dev/./urandom","-jar","/app.jar"]`
	dockerfilePath := filepath.Join(projectInnerPath, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
	if err != nil {
		fmt.Printf("Error creating Dockerfile: %v\n", err)
		return
	}

	fmt.Println("Project scaffolding created successfully.")
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(coffeeCmd)
}
