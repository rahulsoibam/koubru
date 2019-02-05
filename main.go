package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	sendgrid "github.com/sendgrid/sendgrid-go"

	"github.com/rahulsoibam/koubru/authutils"
	"github.com/rahulsoibam/koubru/logger"
	koubrumiddleware "github.com/rahulsoibam/koubru/middleware"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/go-redis/redis"
	"github.com/sony/sonyflake"

	// Koubru Apps
	"github.com/rahulsoibam/koubru/apps/auth"
	"github.com/rahulsoibam/koubru/apps/categories"
	"github.com/rahulsoibam/koubru/apps/countries"
	"github.com/rahulsoibam/koubru/apps/explore"
	"github.com/rahulsoibam/koubru/apps/feed"
	"github.com/rahulsoibam/koubru/apps/opinions"
	"github.com/rahulsoibam/koubru/apps/search"
	"github.com/rahulsoibam/koubru/apps/topics"
	"github.com/rahulsoibam/koubru/apps/user"
	"github.com/rahulsoibam/koubru/apps/users"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var (
	db               *sql.DB
	authCache        *redis.Client
	cache            *redis.Client
	authDB           *sql.DB
	uploader         *s3manager.Uploader
	flake            *sonyflake.Sonyflake
	setupOnce        sync.Once
	sendgridClient   *sendgrid.Client
	argon2Params     *authutils.Params
	koubruMiddleware *koubrumiddleware.Middleware
	logg             *logger.Logger
)

// MaxUploadSize is the max upload size of videos (including accompanying form data) in bytes
const MaxUploadSize = 200 << 20

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	initializeLogger()
	initializeDB()
	initializeAuthDB()
	initializeAuthCache()
	initializeS3Uploader()
	initializeSonyflake()
	initializeArgon2Params()
	initializeSendgridClient()
	initializeKoubruMiddleware()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("."))
	})
	fa := feed.App{}
	r.Mount("/", fa.Routes())
	oa := opinions.App{
		DB:         db,
		Cache:      cache,
		Middleware: koubruMiddleware,
		Log:        logg,
		Flake:      flake,
		Uploader:   uploader,
	}
	r.Mount("/opinions", oa.Routes())
	aa := auth.App{
		AuthCache:  authCache,
		Middleware: koubruMiddleware,
		DB:         db,
		AuthDB:     authDB,
		Log:        logg,
		// SendgridClient: sendgridClient,
		Argon2Params: argon2Params,
	}
	r.Mount("/auth", aa.Routes())

	ua := user.App{
		DB:         db,
		Cache:      cache,
		Middleware: koubruMiddleware,
		Log:        logg,
	}
	r.Mount("/user", ua.Routes())

	usa := users.App{
		DB:         db,
		Cache:      cache,
		Middleware: koubruMiddleware,
		Log:        logg,
	}
	r.Mount("/users", usa.Routes())

	ca := categories.App{
		DB:         db,
		Middleware: koubruMiddleware,
		Log:        logg,
	}

	ta := topics.App{
		DB:         db,
		Middleware: koubruMiddleware,
		Log:        logg,
	}
	r.Mount("/topics", ta.Routes())

	r.Mount("/categories", ca.Routes())
	coa := countries.App{}
	r.Mount("/countries", coa.Routes())
	ea := explore.App{}
	r.Mount("/explore", ea.Routes())
	sa := search.App{}
	r.Mount("/search", sa.Routes())

	log.Fatal(http.ListenAndServe(os.Getenv("API_PORT"), r))

}

func initializeLogger() {
	logg = logger.NewLogger(os.Stderr, os.Stdout, os.Stdout, os.Stderr)
}

// Initialize sets up the database connection, s3 session and routes for the app
func initializeAuthDB() {
	var err error
	authDBString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("AUTH_DB_USERNAME"),
		os.Getenv("AUTH_DB_PASSWORD"),
		os.Getenv("AUTH_DB_HOST"),
		os.Getenv("AUTH_DB_PORT"),
		os.Getenv("AUTH_DB_NAME"))
	authDB, err = sql.Open("postgres", authDBString)
	if err != nil {
		log.Fatal("credsdb: ", err)
	}
	err = authDB.Ping()
	if err != nil {
		log.Fatal("authDB: ", err)
	}
}

func initializeAuthCache() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("AUTH_REDIS_ADDRESS"),
		Password: os.Getenv("AUTH_REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("redis: ", err)
	}
	authCache = client
}

func initializeCache() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CACHE_REDIS_ADDRESS"),
		Password: os.Getenv("CACHE_REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal("redis: ", err)
	}
	cache = client
}

func initializeDB() {
	var err error
	dnsString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	db, err = sql.Open("postgres", dnsString)
	if err != nil {
		log.Fatal("db: ", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("db: ", err)
	}
}

func initializeS3Uploader() {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	if err != nil {
		log.SetOutput(os.Stdout)
		log.Fatal("s3Uploader: ", err)
	}
	uploader = s3manager.NewUploader(sess)
	if err != nil {
		log.Fatal("s3Uploader: ", err)
	}
}

func initializeSonyflake() {
	// Snowflake inspired UUID Generator
	settings := sonyflake.Settings{}
	settings.StartTime = time.Now().UTC()
	flake = sonyflake.NewSonyflake(settings)
}

func initializeArgon2Params() {
	argon2Params = &authutils.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func initializeSendgridClient() {
	sendgridClient = sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
}

func initializeKoubruMiddleware() {
	koubruMiddleware = &koubrumiddleware.Middleware{
		AuthCache: authCache,
		DB:        db,
		Log:       logg,
	}
}
