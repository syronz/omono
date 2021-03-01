package corstartoff

import (
	"log"
	"omono/internal/core"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB initiate the db connection by getting help from gorm
func ConnectDB(engine *core.Engine, printQueries bool) {
	var err error

	var logLevel logger.LogLevel
	switch printQueries {
	case false:
		logLevel = logger.Silent
	case true:
		logLevel = logger.Info
	}

	// engine.Sqldb.SetMaxIdleConns(20)
	// engine.Sqldb.SetMaxOpenConns(1000)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 10 * time.Second, // Slow SQL threshold
			LogLevel:      logLevel,         // Log level
			Colorful:      true,             // Disable color
		},
	)

	if engine.Envs[core.DatabaseDataType] == "mysql" {
		engine.DB, err = gorm.Open(mysql.Open(engine.Envs[core.DatabaseDataWriteDSN]),
			&gorm.Config{
				Logger: newLogger,
			})

		if err != nil {
			log.Fatalln(err.Error())
		}

		mDB, err := engine.DB.DB()
		if err != nil {
			log.Fatalln(err.Error())
		}

		mDB.SetMaxIdleConns(10)
		mDB.SetMaxOpenConns(100)
		// mDB.SetConnMaxLifetime(time.Minute * 30)
		mDB.SetConnMaxIdleTime(time.Hour)

		engine.ReadDB, err = gorm.Open(mysql.Open(engine.Envs[core.DatabaseDataReadDSN]),
			&gorm.Config{
				Logger: newLogger,
			})

		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	// engine.DB.LogMode(engine.Envs.ToBool(core.DatabaseDataLog))

	// if printQueries {
	// 	engine.DB.LogMode(true)
	// }
}

// ConnectActivityDB initiate the db connection by getting help from gorm
func ConnectActivityDB(engine *core.Engine) {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	// engine.ActivityDB, err = gorm.Open(engine.Envs[core.DatabaseActivityType],
	// 	engine.Envs[core.DatabaseActivityDSN])
	engine.ActivityDB, err = gorm.Open(mysql.Open(engine.Envs[core.DatabaseActivityDSN]),
		&gorm.Config{
			Logger: newLogger,
		})
	if err != nil {
		log.Fatalln(err.Error())
	}

	// engine.ActivityDB.LogMode(engine.Envs.ToBool(core.DatabaseActivityLog))

}
