package database

import (
	"io"
	syslog "log"
	"os"
	"strings"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	toolsConfig "github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	mycasbin "github.com/go-admin-team/go-admin-core/sdk/pkg/casbin"
	toolsDB "github.com/go-admin-team/go-admin-core/tools/database"

	// . "github.com/go-admin-team/go-admin-core/tools/gorm/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/global"
)

// Setup 配置数据库
func Setup() {
	for k := range toolsConfig.DatabasesConfig {
		db := setupSimpleDatabase(k, toolsConfig.DatabasesConfig[k])
		if k == "default" {
			e := mycasbin.Setup(db, "sys_")
			sdk.Runtime.SetCasbin("*", e)
		}
	}
}

func setupSimpleDatabase(host string, c *toolsConfig.Database) *gorm.DB {
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	log.Infof("%s => %s", host, pkg.Green(c.Source))
	registers := make([]toolsDB.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = toolsDB.NewResolverConfigure(
			c.Registers[i].Sources,
			c.Registers[i].Replicas,
			c.Registers[i].Policy,
			c.Registers[i].Tables)
	}
	resolverConfig := toolsDB.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifeTime, registers)
	db, err := resolverConfig.Init(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// Logger: New(
		// 	logger.Config{
		// 		SlowThreshold: time.Second,
		// 		Colorful:      true,
		// 		LogLevel: logger.LogLevel(
		// 			log.DefaultLogger.Options().Level.LevelForGorm()),
		// 	},
		// ),
	}, opens[c.Driver])

	SetLogger(db, toolsConfig.LoggerConfig.Path, 500)

	if err != nil {
		log.Fatal(pkg.Red(c.Driver+" connect error :"), err)
	} else {
		log.Info(pkg.Green(c.Driver + " connect success !"))
	}

	sdk.Runtime.SetDb(host, db)
	return db
}

func SetLogger(db *gorm.DB, LogPath string, SlowThresholdMs int) {
	var writer io.Writer = os.Stdout
	if LogPath != "" && LogPath != "default" {
		logFile := strings.TrimSuffix(LogPath, "/") + "/sql.log"
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("failed to set logger :" + err.Error())
		}
		writer = file
	}

	db.Logger = logger.New(
		syslog.New(writer, "\r\n", syslog.Ltime),
		logger.Config{
			SlowThreshold:             time.Millisecond * time.Duration(SlowThresholdMs), // Slow SQL threshold
			LogLevel:                  logger.Info,                                       // Log level
			IgnoreRecordNotFoundError: true,                                              // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                              // Disable color
		},
	)
}
