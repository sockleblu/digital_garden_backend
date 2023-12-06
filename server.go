package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/sockleblu/digital_garden_backend/graph"
	"github.com/sockleblu/digital_garden_backend/graph/auth"
	"github.com/sockleblu/digital_garden_backend/graph/generated"
	"github.com/sockleblu/digital_garden_backend/graph/model"
)

const (
	dbport      = 5432
	graphqlPort = "1337"
)

var db *gorm.DB

func main() {
	digidenEnv := os.Getenv("DIGIDEN_ENV")
	if digidenEnv == "" {
		digidenEnv = "dev"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = graphqlPort
	}

	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		panic("No DATABASE_HOST env variable set")
	}

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		panic("No DATABASE_USER env variable set")
	}

	password := os.Getenv("DATABASE_PASS")
	if password == "" {
		panic("No DATABASE_PASS env variable set")
	}

	dbname := os.Getenv("DATABASE_NAME")
	if dbname == "" {
		panic("No DATABASE_NAME env variable set")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s database=%s sslmode=disable",
		host, dbport, user, password, dbname)

	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{}, &model.Article{}, &model.Tag{})

	router := chi.NewRouter()
	allowed_origins := []string{"http://localhost:3000", "http://localhost:1337", "http://kylekennedy.dev", "https://kylekennedy.dev"}

	router.Use(middleware.Logger)

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowed_origins,
		//AllowedOrigins: []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Origin", "Accept", "X-Requested-With"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           3600,
		Debug:            true,
	}))

	router.Use(auth.Middleware(db))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				if r.Method == "OPTIONS" {
					return true
				}

				return r.Host == "kylekennedy.dev"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	if digidenEnv != "prod" {
		router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	}
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	cfg := &tls.Config{}

	cert, err := tls.LoadX509KeyPair("/etc/ssl/kylekennedy.dev.crt", "/etc/ssl/kylekennedy.dev.key")

	if err != nil {
		log.Fatal(err)
	}

	cfg.Certificates = append(cfg.Certificates, cert)

	cfg.BuildNameToCertificate()

	server := http.Server{
		Addr:      ":1337",
		Handler:   router,
		TLSConfig: cfg,
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}
