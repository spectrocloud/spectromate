package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"

	_ "go.uber.org/automaxprocs"
	"spectrocloud.com/docs-slack-bot/commands"
	"spectrocloud.com/docs-slack-bot/internal"
)

const (
	db_driver string = "postgres"
)

var (
	// dbName              string
	// dbUser              string
	// dbPassword          string
	// dbHost              string
	// dbPort              int64
	globalTraceLevel string
	// globalDb            *sqlx.DB
	globalHost          string
	globalPort          string
	globalHostURL       string = globalHost + ":" + globalPort
	globalSigningSecret string
)

func init() {
	globalTraceLevel = strings.ToUpper(internal.Getenv("TRACE", "INFO"))
	internal.InitLogger(globalTraceLevel)
	globalSigningSecret = internal.Getenv("SLACK_SIGNING_SECRET", "")
	// initDB := strings.ToLower(internal.Getenv("DB_INIT", "false"))
	port := internal.Getenv("PORT", "3000")
	host := internal.Getenv("HOST", "0.0.0.0")
	globalHost = host
	globalPort = port
	globalHostURL = host + ":" + port
	// dbName = internal.Getenv("DB_NAME", "counter")
	// dbUser = internal.Getenv("DB_USER", "postgres")
	// dbHost = internal.Getenv("DB_HOST", "0.0.0.0")
	// dbEncryption := internal.Getenv("DB_ENCRYPTION", "disable")
	// dbPassword = internal.Getenv("DB_PASSWORD", "password")
	// dbPort = internal.StringToInt64(internal.Getenv("DB_PORT", "5432"))
	// db, err := sqlx.Open(db_driver, fmt.Sprintf(
	// 	"host=%s port=%d dbname=%s user=%s password=%s connect_timeout=5 sslmode=%s",
	// 	dbHost,
	// 	dbPort,
	// 	dbName,
	// 	dbUser,
	// 	dbPassword,
	// 	dbEncryption,
	// ))
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("Error connecting to database")
	// }

	// db.SetConnMaxIdleTime(45 * time.Second)
	// db.SetMaxIdleConns(3)
	// db.SetConnMaxLifetime(1 * time.Minute)

	// log.Debug().Msg("Checking database connection...")
	// err = db.Ping()
	// if err != nil {
	// 	log.Debug().Msg("Database is not available")
	// 	log.Fatal().Err(err).Msg("Error connecting to database")
	// }

	// if initDB == "true" {
	// 	log.Debug().Msg("Initializing database")
	// 	err = internal.InitDB(context.Background(), db)
	// 	if err != nil {
	// 		log.Fatal().Err(err).Msg("Error initializing database")
	// 	}
	// }

	// globalDb = db
}

func main() {
	ctx := context.Background()
	healthRoute := commands.NewHealthHandlerContext(ctx)
	helpRoute := commands.NewHelpHandlerContext(ctx, globalSigningSecret)

	http.HandleFunc(internal.ApiPrefix+"health", healthRoute.HealthHTTPHandler)
	http.HandleFunc(internal.ApiPrefix+"help", helpRoute.HelpHTTPHandler)

	log.Info().Msgf("Server is configured for port %s and listing on %s", globalPort, globalHostURL)
	// log.Info().Msgf("Database is configured for %s:%d", dbHost, dbPort)
	log.Info().Msgf("Trace level set to: %s", globalTraceLevel)
	// log.Info().Msgf("Authorization is set to: %v", globalAuthorization)
	log.Info().Msg("Starting server...")
	err := http.ListenAndServe(globalHostURL, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("There's an error with the server")
	}

}
