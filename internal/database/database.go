package database

import (
	"context"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	InitDB(models ...any) error

	GetDB(ctx context.Context) *gorm.DB
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error

	ValStore
	ValkeyService
}

type service struct {
	db *gorm.DB
	valkeyclient valkey.Client
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	dbpassword   = os.Getenv("DB_PASSWORD")
	dbusername   = os.Getenv("DB_USERNAME")
	dbport       = os.Getenv("DB_PORT")
	dbhost       = os.Getenv("DB_HOST")
	dbtableprefix = os.Getenv("DB_TABLE_PREFIX")
	isSingularTable, _ = strconv.ParseBool(os.Getenv("DB_SINGULAR_TABLE"))
	maxOpenConns, _  = strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _  = strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	dbMaxLifetime, _ = strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFETIME"))
	valHost          = os.Getenv("VALKEY_CLIENT_HOST")
	valPort          = os.Getenv("VALKEY_CLIENT_PORT")
	valPass          = os.Getenv("VALKEY_CLIENT_PASSWORD")
	dbInstance *service
	ctxTxKey = "Tx"
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbusername, dbpassword, dbhost, dbport, dbname)
	db,err := gorm.Open(mysql.Open(dsn),&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable:isSingularTable,
			TablePrefix: dbtableprefix,
		},
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Logger.Fatal().Err(err).Msgf("Failed to connect to database due to %v", err)
	}
	sqlDB, err := db.DB()
		if err != nil {
			log.Logger.Fatal().Err(err).Msg("failed to get sql db")
			os.Exit(2)
		}
		sqlDB.SetMaxOpenConns(maxOpenConns)
		sqlDB.SetMaxIdleConns(maxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(dbMaxLifetime) * time.Second)
		valClient, err := valkey.NewClient(valkey.ClientOption{
				InitAddress: []string{valHost + ":" + valPort},
				Password:    valPass,
			})
		if err != nil {
			log.Logger.Fatal().Err(err).Msg("failed to get valkey client")
			os.Exit(2)
		}

	dbInstance = &service{
		db: db,
		valkeyclient: valClient,
	}
	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stats := make(map[string]string)
	sqlDB, err := s.db.DB()
	if err != nil {
		stats["mariadb_status"] = "down"
		stats["mariadb_error"] = fmt.Sprintf("mariadb down: %v", err)
		log.Logger.Fatal().Err(err).Msg("mariadb down")
		return stats
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		stats["mariadb_status"] = "down"
		stats["mariadb_error"] = fmt.Sprintf("mariadb down: %v", err)
		log.Logger.Fatal().Err(err).Msg("mariadb down")
		return stats
	}
	stats["mariadb_status"] = "up"
	stats["mariadb_message"] = "It's healthy"

	dbStats := sqlDB.Stats()
	stats["mariadb_open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["mariadb_in_use"] = strconv.Itoa(dbStats.InUse)
	stats["mariadb_idle"] = strconv.Itoa(dbStats.Idle)
	stats["mariadb_wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["mariadb_wait_duration"] = dbStats.WaitDuration.String()
	stats["mariadb_max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["mariadb_max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["mariadb_message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["mariadb_message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["mariadb_message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["mariadb_message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}
	valResult := s.valkeyclient.Do(ctx, s.valkeyclient.B().Ping().Build())
	if valResult.Error() != nil {
		stats["valkey_status"] = "down"
		stats["valkey_error"] = fmt.Sprintf("valkey down: %v", valResult.Error())
		log.Logger.Fatal().Err(valResult.Error()).Msg("valkey down")
		return stats
	}
	stats["valkey_status"] = "up"
	stats["valkey_message"] = "It's healthy"
	valStatus := parseValkeyInfo(valResult.String())
	maps.Copy(stats, valStatus)
	return stats
}
func parseValkeyInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	sqlDB, _ := s.db.DB()
	log.Logger.Info().Msgf("Disconnected from %s database", dbname)
	s.valkeyclient.Close()
	log.Logger.Info().Msgf("Disconnected from %s valkey", valHost)
	return sqlDB.Close()
}

// GetDB return tx
// If you need to create a Transaction, you must call DB(ctx) and Transaction(ctx,fn)
func (s *service) GetDB(ctx context.Context) *gorm.DB {
	if ctx != nil {
		if tx, ok := ctx.Value(ctxTxKey).(*gorm.DB); ok {
			return tx
		}
		return s.db.WithContext(ctx)
	}
	return s.db
}

func (s *service) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, ctxTxKey, tx)
		return fn(ctx)
	})
}

func(s *service) InitDB(models ...any) error{
	if len(models) == 0 {
		return nil
	}
	return s.db.AutoMigrate(models...)
}
