package adapters

import (
	"fmt"
	"config"
	_ "github.com/jinzhu/gorm/dialects/postgres"      // Blank import for init() invocation
	"github.com/jinzhu/gorm"
	"logger"
)

var dbcfg *gorm.DB

func Init() {
	var err error
	conf := config.GetConfig()

	// Connect to DBCFG
	dbcfg, err = connectToDB(conf.GetString("dbcfg.db_user"), conf.GetString("dbcfg.db_password"),
		conf.GetString("dbcfg.db_host"), conf.GetString("dbcfg.db_port"),
		conf.GetString("dbcfg.db_name"), conf.GetString("dbcfg.db_type"),
		conf.GetString("dbcfg.db_ssl_mode"))
	if err != nil {
		logger.Log.Fatal("Connection to DBCFG failed with error: " + err.Error())
	}
}

// GetDbcfg for getting DBCFG database
func GetDbcfg() *gorm.DB {
	return dbcfg
}

func connectToDB(user, password, host, port, dbName, dbType, sslMode string) (db *gorm.DB, err error){
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbName, sslMode)
	db, err = gorm.Open(dbType, connectionString)
	if err != nil {
		return nil, err
	}
	err = db.DB().Ping() // check the database connectivity
	if err != nil {
		return nil, err
	}
	return db, nil
}
