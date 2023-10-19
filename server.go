package main

import (
	"fmt"
	"log"
	"net/http"
	"crypto/tls"
	"os"
	"slices"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

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
	// pgx (the driver used for Gorm -> psql) enables prepared statement cache by default. To disable...
	//db, err = gorm.Open(postgres.New(postgres.Config{
	//	DSN:                  psqlInfo,
	//	PreferSimpleProtocol: true,
	//}), &gorm.Config{})

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.LogMode(true)
	db.AutoMigrate(&model.User{}, &model.Article{}, &model.Tag{})
	//db.Model(&model.Article{}).AddForeignKey("article_id", "tags(id)", "RESTRICT", "RESTRICT")

	router := chi.NewRouter()
	allowed_domains := []string{"http://localhost:3000", "https://kylekennedy.dev"}

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   allowed_domains,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Use(auth.Middleware(db))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{DB: db}}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				// return r.Host == "kylekennedy.dev"
				return slices.Contains(allowed_domains, r.Host)
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, router))

	cfg := &tls.Config{}

	cert, err := tls.LoadX509KeyPair("/etc/ssl/kylekennedy.dev.crt", "/etc/ssl/kylekennedy.dev.key")

	if err != nil {
    		log.Fatal(err)
	}

	cfg.Certificates = append(cfg.Certificates, cert)

	cfg.BuildNameToCertificate()

	server := http.Server{
    		Addr:      "kylekennedy.local:1337",
    		Handler:   router,
    		TLSConfig: cfg,
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
	//log.Fatal(http.ListenAndServeTLS(":"+port, "/etc/ssl/kylekennedy.local.crt", "/etc/ssl/kylekennedy.local.key", router))
}
