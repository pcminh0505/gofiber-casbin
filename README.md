# GoFiber Admin Client with Casbin RBAC Boilerplate

This boilerplate supports a web server with Authentication (Login / Logout) and User Management. `Admin` role can have all accessibility to the database, while `User` can only change their password.

## 🛠 Installation Guideline

Download and install [Golang](https://go.dev/doc/install) | [Docker](https://www.docker.com/products/docker-desktop/) | `Make` for Windows

For local setup testing:

1. Make sure you have installed `PostgreSQL`, create role (user) and have a local DB in your computer. [Mac](https://www.sqlshack.com/setting-up-a-postgresql-database-on-mac/) | [Windows](https://www.microfocus.com/documentation/idol/IDOL_12_0/MediaServer/Guides/html/English/Content/Getting_Started/Configure/_TRN_Set_up_PostgreSQL.htm) or similar tutorials

2. Install [Postman](https://www.postman.com/) (Web with [Postman Agent](https://www.postman.com/downloads/postman-agent/) for `localhost` request, or full support [Desktop](https://www.postman.com/downloads/)) to test calling the API

## 📦 Packages Dependency

- [casbin](https://github.com/casbin/casbin)
- [casbin/gorm-adapter](https://github.com/casbin/gorm-adapter)
- [go-fiber](https://github.com/gofiber/fiber)
- [go-fiber/jwt](https://github.com/gofiber/jwt)
- [go-fiber/swagger](https://github.com/gofiber/swagger)
- [go-gorm](https://github.com/go-gorm)
- [go-gorm/postgres](https://github.com/go-gorm/postgres)
- [govalidator](https://github.com/asaskevich/govalidator)

## ⚡️ Quick Start

**Generate ECDSA256 Private Key if not existed**

1. Install `OpenSSL` in your computer. Remember `Add to PATH` for `Windows` user
2. Run the below code (detailed command can be found in `Makefile`)

```
make generate-ecdsa
```

**Backend in Local machine**

1. Install the dependencies

```
go get
```

2. Make a copy of `env.template` and rename to `env`

3. Start the server

```
go run .
```

4. For API document, please visit `localhost:8000/swagger`
5. **IMPORTANT! ON FIRST RUN:** A default admin user will be created with the config in the `.env` file

## 📖 Generating Swagger API Document

1. Add comments to your API source code, See [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format).

2. Download `Swag` by using:

```
go install github.com/swaggo/swag/cmd/swag@latest
```

3. Run the [Swag](https://github.com/swaggo/swag) in your Go project root folder which contains `main.go` file, Swag will parse comments and generate required files (`docs` folder and `docs/doc.go`). All the packages have been imported to `main.go` for dependency parsing when init document.

```
swag init
```

4. (Optional) Use `swag fmt` format the SWAG comment. (Please upgrade to the latest version)

```
swag fmt
```

## 🗄 Project Structure

### ./api

**Folder with business logic only**. This directory doesn't care about _what database driver you're using_ or any third-party things.

- `./api/controllers` folder for functional controllers (used in routes)
- `./api/models` folder for describing business models and methods
- `./api/repositories` folder for describing queries for models
- `./api/routes` folder for describing routes
- `./api/utils` folder for utility functions

### ./config

**Folder with configuration**. This directory contains utility functions for backend configuration (Eg. _get environment variables_)

### ./docs

**Folder with API Documentation**. This directory contains config files for auto-generated API Docs by Swagger.

### ./middleware

**Folder with supported middleware**. This directory contains all middleware (Fiber built-in and customization in the future)

### ./infra

**Folder with infrastructure-level logic**. This directory contains all the platform-level logic that will build up the actual project, like _setting up the database_

- `./infra/database` folder with database setup function (PostgreSQL)
