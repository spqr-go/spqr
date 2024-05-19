# SPQR

SPQR is a tool for creating Go projects using hexagonal architecture easily and simply. It provides a structured way to set up your Go project and helps you adhere to hexagonal architecture principles.

## Features

- Setup Go projects with predefined scaffolding following hexagonal architecture.
- Support for Postgres and MariaDB databases.
- Generation of necessary files to quickly start a Go project.
- Integration with Docker and Docker Compose for easy dependency management.

## Requirements

- Go 1.16 or higher
- Docker
- Docker Compose

## Installation

To install SPQR, clone this repository and build the project:

```bash
git clone https://github.com/your_username/spqr.git
cd spqr
go build
```


## Usage
To create a new project with SPQR, use the create command:


```bash
./spqr create
```
You will be prompted to enter the following information:
```bash
Project name
Database type (Postgres or MariaDB)
Database user
Database password
Database name
Database port
```

Example
```bash
Copy code
./spqr create
Enter project name: my_project
Select database (1 for Postgres 🐘, 2 for Mariadb 🦭): 1
Enter DB_USER: user
Enter DB_PASSWORD: password
Enter DB_NAME: mydb
Enter DB_PORT: 5432
```
## Project Structure
The generated project will have the following structure:

```bash
my_project/
├── cmd/
│   └── api/
│       └── spqr.go
├── internal/
│   ├── adapters/
│   │   ├── in/
│   │   │   ├── models/
│   │   │   │   ├── request/
│   │   │   │   └── response/
│   │   └── out/
│   │       └── repositories/
│   ├── configs/
│   ├── core/
│   │   ├── auth/
│   │   ├── domain/
│   │   ├── ports/
│   │   │   ├── in/
│   │   │   └── out/
│   │   └── usecases/
│   └── routing/
├── Dockerfile
└── docker-compose.yml
```

## Contributions
Contributions are welcome! If you want to contribute to this project, please follow these steps:

Fork the project.
Create a new branch (git checkout -b feature/new-feature).
Make your changes and commit them (git commit -am 'Add new feature').
Push to the branch (git push origin feature/new-feature).
Open a Pull Request.
License
This project is licensed under the MIT License. See the LICENSE file for more details.

## Contact
- Author: Galileo Luna
- Email: galileoluna1997@gmail.com

Thank you for using SPQR!


Let me know if you need any more help!






