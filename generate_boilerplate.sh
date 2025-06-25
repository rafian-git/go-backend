#!/usr/bin/env bash

# Script to generate boilerplate code for a microservice

# Variables
SERVICE_NAME=$1

# Replace hyphens with spaces, capitalize each word, then join them for CamelCase (e.g., SmsPusher)
CAPITALIZED_SERVICE_NAME="$(echo "$SERVICE_NAME" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2))}1' | sed 's/ //g')"

# Convert service name to camelCase (e.g., smsPusher)
CAMEL_CASE_SERVICE_NAME="$(echo "$CAPITALIZED_SERVICE_NAME" | awk '{print tolower(substr($0,1,1)) substr($0,2)}')"

# Convert service name to all uppercase and replace hyphens with underscores for ALL_CAPS format (e.g., SMS_PUSHER)
ALL_CAPITALIZED_SERVICE_NAME="$(echo "$SERVICE_NAME" | tr '-' '_' | awk '{print toupper($0)}')"

# Convert service name to snake_case (e.g., sms_pusher)
SNAKE_CASE_SERVICE_NAME="$(echo "$SERVICE_NAME" | tr '-' '_' | awk '{print tolower($0)}')"

# Convert service name to initials (e.g., om)
INITIALS="$(echo "$SERVICE_NAME" | awk -F'-' '{for(i=1;i<=NF;i++) printf "%s", substr($i,1,1)}' | tr '[:upper:]' '[:lower:]')"

MODULE_NAME="github.com/rafian-git/$SERVICE_NAME"
REPO_NAME="github.com/rafian-git"

ROOT_DIR=$(dirname $(pwd))
SERVICE_DIR="$ROOT_DIR/$SERVICE_NAME"
CONFIG_DIR="$SERVICE_DIR/config"
PKG_DIR="$SERVICE_DIR/pkg/grpc"
CLI_DIR="$SERVICE_DIR/pkg/client"
INTERNAL_DIR="$SERVICE_DIR/internal/$SERVICE_NAME"
CMD_DIR="$SERVICE_DIR/cmd"
SQL_DIR="$SERVICE_DIR/sql"

# Output the transformed service names
echo "Capitalized Service Name (PascalCase): $CAPITALIZED_SERVICE_NAME"
echo "Camel Case Service Name: $CAMEL_CASE_SERVICE_NAME"
echo "All Capitalized Service Name: $ALL_CAPITALIZED_SERVICE_NAME"

# Function to create folder structure
create_folders() {
  echo "Creating folder structure for $SERVICE_NAME..."
  mkdir -p $CONFIG_DIR
  mkdir -p $PKG_DIR
  mkdir -p $CLI_DIR
  mkdir -p $INTERNAL_DIR/models/pg
  mkdir -p $INTERNAL_DIR/service
  mkdir -p $INTERNAL_DIR/server
  mkdir -p $INTERNAL_DIR/config
  mkdir -p $CMD_DIR
  mkdir -p $SQL_DIR/migrations
}

# Function to generate configuration files
generate_config_files() {
  echo "Generating config files..."
  cat <<EOT > $CONFIG_DIR/config.dev.yaml
# Development environment configuration

time_limit: 3000

# gRPC configurations
grpc:
  service_name: $SNAKE_CASE_SERVICE_NAME
  network: tcp
  bind: 0.0.0.0:15071
  max_recv_msg_size: 4194304
  max_send_msg_size: 4194304
  read_buffer_size: 4194304
  write_buffer_size: 4194304

# Database configuration
db:
  url: 'postgresql://root:123456@localhost:5432/$SNAKE_CASE_SERVICE_NAME?sslmode=disable'
EOT
}



# Function to generate main.go
generate_main_go() {
  echo "Generating main.go for $SERVICE_NAME..."
  cat <<EOT > $CMD_DIR/main.go
package main

import (
  "context"
  "github.com/spf13/cobra"
  "github.com/rafian-git/go-backend/pkg/log"
  "github.com/rafian-git/go-backend/pkg/rabbitmq"
  "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/config"
  "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/models/pg"
  "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/service"
  "os"
)

func main() {
  var (
    ctx    = context.Background()
    conf   = config.New().Load()
    err    error
    logger = log.New()
  )

  logger = logger.Named("$SERVICE_NAME")

  // database
  postgres, err := pg.New(conf.DB, logger)
  if err != nil {
    panic(err)
  }
  defer func(postgres *pg.DB) {
    _ = postgres.DB.Close()
  }(postgres)

  // Message Queue
  qu, err := rabbitmq.New(conf.Queue, logger)
  if err != nil {
    panic(err)
  }

  // service
  ${CAMEL_CASE_SERVICE_NAME}Init := &service.Init{
    Db:        postgres,
    Cnf:       conf,
    Log:       logger,
    Qu:        qu,
    PubConfig: conf.PubSub,
  }

  var rootCmd = &cobra.Command{}
  rootCmd.AddCommand(serve(ctx, ${CAMEL_CASE_SERVICE_NAME}Init))

  if err = rootCmd.Execute(); err != nil {
    logger.Error(ctx, err.Error())
    os.Exit(1)
  }
}
EOT
}




# Function to generate serve.go
generate_serve_go() {
  echo "Generating serve.go in $CMD_DIR..."
  cat <<EOT > $CMD_DIR/serve.go
package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/rafian-git/go-backend/pkg/log"
	"github.com/rafian-git/go-backend/pkg/migration"
	"github.com/rafian-git/go-backend/pkg/sigint"
	"github.com/rafian-git/go-backend/pkg/server"
	"github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/service"
	${INITIALS} "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/server"
	"github.com/rafian-git/${SERVICE_NAME}/pkg/grpc"
	"github.com/rafian-git/${SERVICE_NAME}/sql"
	"go.uber.org/zap"
	"net/http"
)

func run(ctx context.Context, init *service.Init) {
  conf := init.Cnf
  logger := log.New().Named("$SERVICE_NAME")

  // migrations
  migrateDirection, migrateOnly := conf.MigrationDirectionFlag()
  migrateDB, err := migration.SQLFromUrl(conf.DB.URL)
  if err != nil {
    panic(err)
  }

  migrations := sql.GetMigrations()
  err = migration.MigrateFromFS(migrateDB, migrateDirection, "$SNAKE_CASE_SERVICE_NAME", migrations)
  if err != nil {
    panic(err)
  }
  _ = migrateDB.Close()

  if migrateOnly {
    logger.Info(ctx, "Migration complete, exiting")
    return
  }

  // service
  var s service.Service
  s, err = service.New(init)
  if err != nil {
    panic(err)
  }

  // server initialization
  var srv *server.Server
  reg := prometheus.NewRegistry()
  if srv, err = server.New(logger, conf.GRPC, reg); err != nil {
    panic(err)
  }

  logger.Info(ctx, "server port will be: ", zap.String("network", conf.GRPC.Bind))
  meServer := ${INITIALS}.New(logger, s)
  grpc.Register${CAPITALIZED_SERVICE_NAME}Server(srv.GRPC(), meServer)

  // Create a HTTP server for prometheus.
  httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{
    EnableOpenMetrics: true,
  }), Addr: fmt.Sprintf("0.0.0.0:%d", 9092)}

  // Start your http server for prometheus.
  go func() {
    if err := httpServer.ListenAndServe(); err != nil {
      logger.Error(ctx, err.Error())
    }
  }()

  srv.Run()
  defer func(srv *server.Server) {
    _ = srv.Close()
  }(srv)

  sigint.Wait()
  logger.Info(ctx, "stopping server!!")
}

func serve(ctx context.Context, init *service.Init) *cobra.Command {
  return &cobra.Command{
    Use: "serve",
    Run: func(cmd *cobra.Command, args []string) {
      run(ctx, init)
    },
  }
}
EOT
}



# Function to generate config.go
generate_config_go() {
  echo "Generating config.go in $INTERNAL_DIR/config"

  cat <<EOT > "$INTERNAL_DIR/config/config.go"
package config

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/rafian-git/go-backend/pkg/bundb"
	"github.com/rafian-git/go-backend/pkg/migration"
	"github.com/rafian-git/go-backend/pkg/pubsub"
	"github.com/rafian-git/go-backend/pkg/rabbitmq"
	"github.com/rafian-git/go-backend/pkg/server"
	"os"
)

// Config of entire application.
type Config struct {
	Url                 string               \`json:"url" yaml:"url" toml:"url" mapstructure:"url"\`
	GRPC                *server.Config       \`json:"grpc" yaml:"grpc" toml:"grpc" mapstructure:"grpc"\`
	DB                  *bundb.Config        \`json:"db" yaml:"db" toml:"db" mapstructure:"db"\` // nolint
	MigrateDirection    migration.Direction  \`json:"migrate"\`
	PubSub              *pubsub.Config       \`json:"pubsub" yaml:"pubsub" toml:"pubsub" mapstructure:"pubsub"\`
	Queue               *rabbitmq.Config     \`json:"rabbit_mq" yaml:"rabbit_mq" toml:"rabbit_mq" mapstructure:"rabbit_mq"\`
}

// New default configurations.
func New() (conf *Config) {
	conf = new(Config)
	conf.GRPC = server.NewConfig()
	conf.PubSub = pubsub.NewConfig()
	return
}

// Load the Config from configuration files. This method panics on error.
func (c *Config) Load() *Config {

	consulPath := os.Getenv("${ALL_CAPITALIZED_SERVICE_NAME}_CONSUL_PATH")
	consulURL := os.Getenv("${ALL_CAPITALIZED_SERVICE_NAME}_CONSUL_URL")

	viper.AddRemoteProvider("consul", consulURL, consulPath)
	viper.SetConfigType("yaml") // Need to explicitly set this to json

	err := viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.Unmarshal(c)

	migrate := c.MigrateDirection
	c.MigrateDirection = migration.Direction(migrate)
	if err = c.MigrateDirection.Check(); err != nil {
		panic(err)
	}

	return c
}

// MigrationDirectionFlag returns migration direction and migrateOnly flag
func (c *Config) MigrationDirectionFlag() (
	migrateDirection migration.Direction, migrateOnly bool) {
	if c.MigrateDirection == "" {
		return migration.DirectionUp, false
	}

	return c.MigrateDirection, true
}

EOT
}

# Function to generate server.go
generate_server_go() {
  echo "Generating server.go in $INTERNAL_DIR..."
  cat <<EOT > $INTERNAL_DIR/server/server.go
package server

import (
  "github.com/rafian-git/go-backend/pkg/log"
  "github.com/gogo/protobuf/types"
  "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/service"
  pb "github.com/rafian-git/${SERVICE_NAME}/pkg/grpc"
  "context"
)

type ${CAPITALIZED_SERVICE_NAME}Server struct {
  log     *log.Logger
  service service.Service
}

func New(logger *log.Logger, service service.Service) (srv *${CAPITALIZED_SERVICE_NAME}Server) {
  srv = new(${CAPITALIZED_SERVICE_NAME}Server)
  srv.service = service
  srv.log = logger.Named("server")
  return
}

// HealthHandler handles the health-check API call
func (v *${CAPITALIZED_SERVICE_NAME}Server) HealthCheck(ctx context.Context, empty *types.Empty) (*pb.HealthCheckResponse, error) {
    v.log.Info(ctx, "Health check handler triggered")
    return v.service.Health().HealthCheck(ctx)
}
EOT
}




# Function to generate service.go
generate_service_go() {
  echo "Generating service.go in $INTERNAL_DIR..."
  cat <<EOT > $INTERNAL_DIR/service/service.go
package service

import (
    "github.com/rafian-git/go-backend/pkg/log"
    "github.com/rafian-git/go-backend/pkg/pubsub"
    "github.com/rafian-git/go-backend/pkg/rabbitmq"
    "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/config"
    "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/models"
)

type Service interface {
    Health() HealthService
}

type ${CAPITALIZED_SERVICE_NAME}Service struct {
    db           model.DB
    config       *config.Config
    log          *log.Logger
    queue        rabbitmq.Queue
    PubSubConfig *pubsub.Config
}


type Init struct {
	Db        model.DB
	Cnf       *config.Config
	Log       *log.Logger
	Qu        rabbitmq.Queue
	PubConfig *pubsub.Config
}

func New(init *Init) (*${CAPITALIZED_SERVICE_NAME}Service, error) {
    return &${CAPITALIZED_SERVICE_NAME}Service{
      db:               init.Db,
      config:           init.Cnf,
      log:              init.Log.Named("service"),
      queue:            init.Qu,
      PubSubConfig:     init.PubConfig,
    }, nil
}
EOT
}


# Function to generate health.go
generate_service_health_go() {
  echo "Generating health.go in $INTERNAL_DIR..."
  cat <<EOT > $INTERNAL_DIR/service/health.go
package service

import (
	"context"
	pb "github.com/rafian-git/${SERVICE_NAME}/pkg/grpc"

)

type HealthService interface {
  HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error)
}

type HealthReceiver struct {
	*${CAPITALIZED_SERVICE_NAME}Service
}

func (ms *${CAPITALIZED_SERVICE_NAME}Service) Health() HealthService {
	return &HealthReceiver{
		ms,
	}
}


func (v *HealthReceiver) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
    v.log.Info(ctx, "Service health check handler")
    err := v.db.Health().Ping(ctx)
    if err != nil {
      return nil, err
    }
    return &pb.HealthCheckResponse{
        Status: "ok",
    }, nil
}

EOT

}


# Function to generate service.go
generate_migration_up() {
  echo "Generating 000000_health_schema.up.sql in $SQL_DIR..."
  cat <<EOT > $SQL_DIR/migrations/000000_health_schema.up.sql
CREATE TABLE IF NOT EXISTS health (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50) NOT NULL,

    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NULL,
    deleted_at  TIMESTAMPTZ DEFAULT NULL
);
EOT

}

generate_migration_down() {
  echo "Generating 000000_health_schema.down.sql in $SQL_DIR..."
  cat <<EOT > $SQL_DIR/migrations/000000_health_schema.down.sql
DROP  TABLE health;
EOT

}

generate_sql_go() {
  echo "Generating service.go in $SQL_DIR..."
  cat <<EOT > $SQL_DIR/sql.go
package sql

import (
	"embed"
	"io/fs"
)

//go:embed migrations
var migrations embed.FS

func GetMigrations() fs.FS {
	return migrations
}


EOT

}


# Function to generate service.go
generate_template_file() {
  echo "Generating service.go in $INTERNAL_DIR..."
  cat <<EOT > $INTERNAL_DIR/config/config.go

EOT

}



# Function to generate models/pg/pg.go
generate_pg_go() {
  echo "Generating pg.go in $INTERNAL_DIR..."
cat <<EOT > $INTERNAL_DIR/models/pg/pg.go

package pg

import (
  "context"
  "database/sql"

  "github.com/uptrace/bun"
  "github.com/rafian-git/go-backend/pkg/bundb"
  "github.com/rafian-git/go-backend/pkg/log"
  model "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/models"
)

// compiler time check
var _ model.DB = (*DB)(nil)

// DB is a database representation
type DB struct {
  // *bun.IDB
  *bun.DB // underlying go-pg DB wrapper instance
  log     *log.Logger
}

// Tx represents transactions
type Tx struct {
  *bun.Tx
  log *log.Logger
}

// compiler time check
var _ model.Repository = (*Tx)(nil)

// New DB with given configurations and logger.
func New(conf *bundb.Config, logger *log.Logger) (db *DB, err error) {

  var pg *bundb.DB
  if pg, err = bundb.New(conf); err != nil {
    return
  }

  db = new(DB)
  db.DB = pg.DB // embed
  db.registerTables()
  db.log = logger.Named("db_model")
  db.log.Info(context.Background(), "db initialization done")
  return
}

// register all tables for relations
func (db *DB) registerTables() {

}

// InTx runs given function in SQL-transaction.
func (db *DB) InTx(ctx context.Context, txFunc model.TxFunc) (err error) {

  err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) (err error) {
    return txFunc(ctx, &Tx{
      &tx,
      db.log,
    })
  })
  return
}
EOT
}

# Function to generate models/pg/pg.go
generate_pg_health_go() {
  echo "Generating pg.go in $INTERNAL_DIR..."
cat <<EOT > $INTERNAL_DIR/models/pg/health.go

package pg

import (
	"context"
  "github.com/rafian-git/go-backend/pkg/log"
	"github.com/uptrace/bun"
	model "github.com/rafian-git/${SERVICE_NAME}/internal/${SERVICE_NAME}/models"
)

func (db *DB) Health() model.HealthRepository {
	return &Health{
  		db.DB,
  		db.log,
  }
}

func (db *Tx) Health() model.HealthRepository {
	return &Health{
  		db.Tx,
  		db.log,
  	}
}

type Health struct {
	bun.IDB
	*log.Logger
}

func (b *Health) Ping(ctx context.Context) error {
	var result int
  	// Run a simple query to check the connection
  	err := b.IDB.NewRaw("SELECT 1").Scan(ctx, &result)
  	if err != nil {
  		b.Logger.Error(ctx, err.Error())
  		return err
  	}
  	return nil
}

EOT
}

# Function to generate models/model.go
generate_model_go() {
echo "Generating model.go in $INTERNAL_DIR..."
cat <<EOT > $INTERNAL_DIR/models/model.go

package model

import (
  "context"
)

type Repository interface {
  Health() HealthRepository
}

// TxFunc represents function to run in an SQL-transaction.
type TxFunc func(ctx context.Context, tx Repository) (err error)

type DB interface {
  // Repository access without transactions
  Repository
  // InTx runs given function in transaction
  InTx(context.Context, TxFunc) error
}


EOT
}

# Function to generate models/health.go
generate_model_health_go() {
  echo "Generating health.go in $INTERNAL_DIR..."

  # Ensure the directory exists
  mkdir -p "$INTERNAL_DIR/models"

  # Use cat to write the content to health.go
  cat <<EOT > "$INTERNAL_DIR/models/health.go"
package model

import (
  "context"
  "github.com/uptrace/bun"
  "time"
)

type HealthRepository interface {
  Ping(ctx context.Context) error
}

type Health struct {
  bun.BaseModel \`bun:"table:health"\`

  ID        int64     \`json:"id" bun:"id,pk,autoincrement"\`
  Status    string    \`json:"status" bun:"status,default:'HEALTHY'"\`
  CreatedAt time.Time \`json:"created_at" bun:"created_at,default:current_timestamp"\`
  UpdatedAt time.Time \`json:"updated_at" bun:"updated_at,nullzero"\`
  DeletedAt time.Time \`json:"-" bun:"deleted_at,nullzero"\`  // Excluded from JSON
}
EOT

  echo "health.go has been generated successfully."
}



# Function to generate service.go
generate_client_go() {
  echo "Generating client.go in $PKG_DIR..."
  cat <<EOT > $CLI_DIR/client.go

package client

import (
	"context"
	"github.com/gogo/protobuf/types"
	pb "github.com/rafian-git/${SERVICE_NAME}/pkg/grpc"
	"google.golang.org/grpc"
)

type Clients interface {
	 HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error)
}

type ${CAPITALIZED_SERVICE_NAME} struct {
	client pb.${CAPITALIZED_SERVICE_NAME}Client
}

func New${CAPITALIZED_SERVICE_NAME}Client(target string) (Clients, error) {

	m := &${CAPITALIZED_SERVICE_NAME}{}
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	m.client = pb.New${CAPITALIZED_SERVICE_NAME}Client(conn)

	return m, err
}

func (m *${CAPITALIZED_SERVICE_NAME}) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
	return m.client.HealthCheck(ctx, &types.Empty{})
}


EOT

}

# Function to generate proto
generate_proto() {
# Create proto file
echo "Generating proto file..."
cat <<EOT > $PKG_DIR/$SERVICE_NAME.proto
syntax = "proto3";

package $SNAKE_CASE_SERVICE_NAME;

option go_package = "github.com/rafian-git/$SERVICE_NAME/pkg/grpc;grpc";

import "google/protobuf/empty.proto";
import "google/api/annotations.proto";

service ${CAPITALIZED_SERVICE_NAME} {
  // Health check method
  rpc HealthCheck (google.protobuf.Empty) returns (HealthCheckResponse) {
    option (google.api.http) = {
      post: "/api/v1/${SERVICE_NAME}/health-check"
      body: "*"
    };
  }
}

// Response for health check
message HealthCheckResponse {
  string status = 1;
}
EOT
}


# Function to generate generate.go
generate_generate_go() {
  echo "Generating generate.go in $PKG_DIR..."
cat <<EOT > $PKG_DIR/generate.go
package grpc

//go:generate sh generate.sh
EOT

}


# Function to generate generate.sh
generate_generate_sh() {
  echo "Generating protobuf compile script..."
  cat <<EOT > "$PKG_DIR/generate.sh"
#!/usr/bin/env bash

GRPC_VER=v1
GRPC_PKG=github.com/grpc-ecosystem/grpc-gateway
GRPC_PATH=\$(go list -m -f "{{.Dir}}" \${GRPC_PKG}@\${GRPC_VER})

GAPI_VER=v1.4.1
GAPI_PKG=github.com/gogo/googleapis
GAPI_PATH=\$(go list -m -f "{{.Dir}}" \${GAPI_PKG}@\${GAPI_VER})

GOGO_VER=v1.3.2
GOGO_PKG=github.com/gogo/protobuf
GOGO_PATH=\$(go list -m -f "{{.Dir}}" \${GOGO_PKG}@\${GOGO_VER})

TECHETRON_ROOT=\$(go list -m -f "{{.Dir}}")
BACKEND_ROOT=\$(dirname "\$TECHETRON_ROOT")"/backend"

APIERROR=\${BACKEND_ROOT}/pkg/apierror

GOGO_ANY=Mgoogle/protobuf/any.proto=\${GOGO_PKG}/types
GOGO_DURATION=Mgoogle/protobuf/duration.proto=\${GOGO_PKG}/types
GOGO_STRUCT=Mgoogle/protobuf/struct.proto=\${GOGO_PKG}/types
GOGO_TIMESTAMP=Mgoogle/protobuf/timestamp.proto=\${GOGO_PKG}/types
GOGO_WRAPPERS=Mgoogle/protobuf/wrappers.proto=\${GOGO_PKG}/types
GOGO_EMPTY=Mgoogle/protobuf/empty.proto=\${GOGO_PKG}/types
GOGO_GAPI=Mgoogle/api/annotations.proto=\${GAPI_PKG}/google/api
GOGO_FLDMSK=Mgoogle/protobuf/field_mask.proto=\${GOGO_PKG}/types

FULL=\${GOGO_ANY},\${GOGO_DURATION},\${GOGO_STRUCT},\${GOGO_TIMESTAMP}
FULL=\${FULL},\${GOGO_WRAPPERS},\${GOGO_EMPTY},\${GOGO_GAPI},\${GOGO_FLDMSK}

protoc -I . \\
    -I \${GOGO_PATH} \\
    -I \${GRPC_PATH} \\
    -I \${GAPI_PATH} \\
    -I \${APIERROR} \\
    --include_imports \\
    --gogofast_out=\${FULL},paths=source_relative,plugins=grpc:. \\
    --swagger_out=logtostderr=true,json_names_for_fields=true:../../docs/swagger/ \\
    --descriptor_set_out=\${BACKEND_ROOT}/integration/etc/envoy/descriptors/${SERVICE_NAME}.desc \\
    *.proto
EOT


}




# Function to generate run_local.sh
generate_run_local_sh() {
  # Confirming ROOT_DIR is correctly set to the current script's root
  echo "Generating run_local.sh in $ROOT_DIR..."

  # Check if ROOT_DIR is accessible and exists
  if [ ! -d "$ROOT_DIR" ]; then
    echo "Error: $ROOT_DIR does not exist or is not a directory."
    return 1
  fi

  # Ensure ALL_CAPITALIZED_SERVICE_NAME is set
  if [ -z "$ALL_CAPITALIZED_SERVICE_NAME" ]; then
    echo "Error: ALL_CAPITALIZED_SERVICE_NAME is not set."
    return 1
  fi

  # Writing the run_local.sh script
  cat <<EOT > "$SERVICE_DIR/run_local.sh"
#!/bin/bash

echo "Starting server"

export ${ALL_CAPITALIZED_SERVICE_NAME}_CONSUL_PATH="${SERVICE_NAME}"
export ${ALL_CAPITALIZED_SERVICE_NAME}_CONSUL_URL="localhost:8500"

# Uploading the configuration to Consul
curl --request PUT --data-binary @config/config.dev.yaml http://127.0.0.1:8500/v1/kv/${SERVICE_NAME}

# Running the Go server with any passed arguments
go run cmd/*.go "\$1"

EOT




  echo "run_local.sh generated successfully."
}



# Function to generate run_local.sh
generate_readme_md() {
  echo "Generating README.MD in $SERVICE_DIR..."
cat <<EOT > "$SERVICE_DIR/README.md"

## Getting started

To make it easy for you to get started with GitLab, here's a list of recommended next steps.

Already a pro? Just edit this README.md and make it your own. Want to make it easy? [Use the template at the bottom](#editing-this-readme)!

EOT

}


# Function to generate Dockerfile
generate_dockerfile() {
  echo "Generating Dockerfile in $SERVICE_DIR..."

  cat <<EOT > "$SERVICE_DIR/Dockerfile"
# Default to Go 1.19
ARG GO_VERSION=1.19

# Start from golang v\${GO_VERSION}-alpine base image as builder stage
FROM golang:\${GO_VERSION}-alpine AS builder

# Add Maintainer Info
LABEL maintainer="Abu Hanifa <a8u.han1fa&gmail.com>"

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \\
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \\
    echo 'nobody:x:65534:' > /user/group

# Install CA certificates and timezone data
RUN apk add --no-cache ca-certificates tzdata

# Set the working directory outside \$GOPATH to enable support for Go modules
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Build the Go app
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 go build -a -installsuffix 'static' -o /app cmd/*.go

######## Start a new stage from scratch #######
# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the builder stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import CA certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Import the compiled executable from the builder stage
COPY --from=builder /app /app

# Expose port 8080 as we are running the executable as an unprivileged user
EXPOSE 8080

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]

EOT

  echo "Dockerfile generated successfully in $ROOT_DIR"
}


# Function to generate .gitignore
generate_gitignore() {
  echo "Generating .gitignore in $SERVICE_DIR..."

  cat <<EOT > "$SERVICE_DIR/.gitignore"
# Ignore IDE-specific files (e.g., JetBrains .idea directory)
.idea

# Ignore vendor directory for Go modules
vendor/

# Add any additional files or directories you want to ignore below

EOT

  echo ".gitignore generated successfully in $ROOT_DIR"
}



# Function to generate .gitlab.ci.yml
generate_gitlab_ci_yml(){
  echo "Generating .gitlab-ci.yml in $SERVICE_DIR..."

  cat <<EOT > "$SERVICE_DIR/.gitlab-ci.yml"
workflow:
  rules:
    - if: \$CI_PIPELINE_SOURCE == 'merge_request_event' || \$CI_COMMIT_BRANCH == 'main'

stages:
  - build_deploy
  - build
  - deploy
  - deploy_prod

before_script:
  - export SOURCE_BRANCH=\$CI_MERGE_REQUEST_SOURCE_BRANCH_NAME
  - SOURCE_BRANCH=\${SOURCE_BRANCH:-\$CI_DEFAULT_BRANCH}
  - eval \$(ssh-agent -s)
  - export SSH_DECODE_PRIVATE_KEY="\$(echo \$SSH_PRIVATE_KEY | base64 -d)"
  - ssh-add <(echo "\$SSH_DECODE_PRIVATE_KEY")
  - echo -e "Host gitlab.techetronventures.com\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config
  - git clone git@gitlab.techetronventures.com:core/backend.git
  - git clone git@gitlab.techetronventures.com:core/portfolio.git
  - git clone -b \$SOURCE_BRANCH git@gitlab.techetronventures.com:core/me.git

after_script:
  - >
    if [ "\$CI_JOB_STATUS" == "success" ]; then
      CI_JOB_STATUS="passed  :gopher:"
    else
      CI_JOB_STATUS="failed  :x:"
    fi
  - >
    curl -i -X POST -H "application/json" \$SLACK_HOOK -d "{\\"text\\": \\" \$SERVICE pipeline's \$CI_JOB_NAME job has been \$CI_JOB_STATUS\$JOB_RETURN\n\$CI_PIPELINE_URL\\"}"

Start Build & Deploy:
  stage: build_deploy
  image: gcr.io/stock-x-342909/docker-golang

  script:
    - echo "Start build & deploy for \$SOURCE_BRANCH"
  rules:
    - if: \$CI_PIPELINE_SOURCE == "merge_request_event"
      when: manual
      allow_failure: false

Build:
  stage: build
  image: gcr.io/stock-x-342909/docker-golang
  services:
    - docker:19.03.12-dind
  variables:
    DOCKER_HOST: tcp://docker:2375/
    DOCKER_DRIVER: overlay2
    DOCKER_TLS_CERTDIR: ""
    SERVICE_IMAGE: gcr.io/\$GCP_PROJECT_NAME/\$SERVICE

  script:
    - echo "Compiling & building the code for \$SOURCE_BRANCH"
    - git config --global user.email "abu.hanifa@techetronventures.com"
    - git config --global url.git@gitlab.techetronventures.com:core.insteadOf https://github.com/rafian-git
    - echo \$SERVICE_ACCOUNT_KEY > key.json
    - docker login -u _json_key --password-stdin https://gcr.io < key.json
    - go version
    - cd \$SERVICE
    - export IMAGE_TAG=\$(git rev-parse --short HEAD)
    - echo "IMAGE_TAG - \$IMAGE_TAG"
    - cd ../backend
    - source ~/.profile
    - go get github.com/gogo/googleapis@v1.4.1
    - go get github.com/grpc-ecosystem/grpc-gateway@latest
    - make build-all service=\$SERVICE
    - make docker-build service=\$SERVICE
    - ./build-envoy.sh \$SERVICE-\$IMAGE_TAG
    - echo "Compile & build complete."

Deploy-Stage:
  stage: deploy
  image: gcr.io/stock-x-342909/gcloud
  environment: stage

  script:
    - cd \$SERVICE
    - export IMAGE_TAG=\$(git rev-parse --short HEAD)
    - echo "IMAGE_TAG - \$IMAGE_TAG"
    - curl -X PUT --data-binary @config/config.stage.yaml --location http://consul-stage.trek.com.bd/v1/kv/\$SERVICE
    - cd ..
    - echo "\$SERVICE_ACCOUNT_KEY" > key.json
    - gcloud auth activate-service-account --key-file=key.json
    - gcloud config set project \$GCP_PROJECT_NAME
    - gcloud config set container/cluster \$GKE_CLUSTER_NAME
    - gcloud config set compute/zone asia-east1-a
    - gcloud container clusters get-credentials \$GKE_CLUSTER_NAME --zone asia-east1-a --project \$GCP_PROJECT_NAME
    - kubectl version --short
    - cd backend/deployment/\$SERVICE
    - envsubst  <deploy.yaml | kubectl apply -f -
    - kubectl rollout status deploy/\$SERVICE -n \$SERVICE --timeout=60s
    - envsubst  <deploy-account-creation.yaml | kubectl apply -f -
    - kubectl rollout status deploy/account-listener -n \$SERVICE --timeout=60s

Deploy-Prod:
  stage: deploy_prod
  image: gcr.io/stock-x-342909/gcloud
  environment: production
  when: manual
  only:
    - main

  script:
    - cd \$SERVICE
    - export IMAGE_TAG=\$(git rev-parse --short HEAD)
    - echo "IMAGE_TAG - \$IMAGE_TAG"
    - curl -X PUT --data-binary @config/config.prod.yaml --location http://consul-prod.trek.com.bd/v1/kv/\$SERVICE
    - cd ..
    - echo "\$SERVICE_ACCOUNT_KEY" > key.json
    - gcloud auth activate-service-account --key-file=key.json
    - gcloud config set project \$GCP_PROJECT_NAME
    - gcloud config set container/cluster \$GKE_CLUSTER_PROD
    - gcloud config set compute/zone asia-east1-a
    - gcloud container clusters get-credentials \$GKE_CLUSTER_PROD --zone asia-east1-a --project \$GCP_PROJECT_NAME
    - kubectl version --short
    - cd backend/deployment/\$SERVICE
    - envsubst  <deploy.yaml | kubectl apply -f -
    - kubectl rollout status deploy/\$SERVICE -n \$SERVICE --timeout=60s
    - envsubst  <deploy-account-creation.yaml | kubectl apply -f -
    - kubectl rollout status deploy/account-listener -n \$SERVICE --timeout=60s
EOT

  echo ".gitlab-ci.yml generated successfully in $ROOT_DIR"
}

# Make the file executable
make_executable(){
  chmod +x "$PKG_DIR/generate.sh"
  chmod +x "$SERVICE_DIR/run_local.sh"
}


# Function to initialize Go module if not already initialized
initialize_go_mod() {   # Module name (passed as an argument to the function)

  echo "Checking for go.mod in $SERVICE_DIR..."

  # Navigate to the service directory
  cd "$SERVICE_DIR" || exit

  # Check if go.mod already exists
  if [ ! -f "go.mod" ]; then
    echo "go.mod not found. Initializing Go module..."
    go mod init "$MODULE_NAME"
    git init
    export GOSUMDB=off && go mod tidy
    echo "Go module initialized for $MODULE_NAME."
  else
    echo "go.mod already exists. Skipping go mod init."
  fi
}

append_go_mod(){

  # Check if go.mod exists
  if [ -f "go.mod" ]; then
    # Append the replace directive to the go.mod file
    echo "replace $REPO_NAME/go-backend => ../go-backend" >> go.mod
    echo "replace directive added to go.mod"
  else
    echo "go.mod file not found!"
  fi
}

generate_proto_pb_build() {
   mkdir docs && mkdir docs/swagger
   go mod tidy
   echo "generating protobuf"
   go generate ./...
   go get -u github.com/gogo/googleapis/google/api
   #indirect dependency version same for bun
   go get -u github.com/uptrace/bun/dialect/pgdialect

   go mod tidy
   go build cmd/*.go
}




# Execute script
create_folders
generate_config_files

generate_main_go
generate_serve_go

generate_config_go
generate_server_go
generate_service_go
generate_service_health_go

generate_migration_up
generate_migration_down
generate_sql_go

generate_pg_go
generate_pg_health_go

generate_model_go
generate_model_health_go

generate_client_go
generate_proto
generate_generate_go
generate_generate_sh

generate_dockerfile
generate_readme_md
generate_run_local_sh
generate_gitlab_ci_yml
generate_gitignore

make_executable
initialize_go_mod
append_go_mod

generate_proto_pb_build


echo "Boilerplate generation complete for $SERVICE_NAME."