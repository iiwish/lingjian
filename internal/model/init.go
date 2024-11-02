package model

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var DB *sqlx.DB

// InitDB 初始化数据库连接
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.mysql.username"),
		viper.GetString("database.mysql.password"),
		viper.GetString("database.mysql.host"),
		viper.GetString("database.mysql.port"),
		viper.GetString("database.mysql.dbname"),
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("connect DB failed: %v", err)
	}

	db.SetMaxIdleConns(viper.GetInt("database.mysql.max_idle_conns"))
	db.SetMaxOpenConns(viper.GetInt("database.mysql.max_open_conns"))
	db.SetConnMaxLifetime(time.Duration(viper.GetInt("database.mysql.conn_max_lifetime")) * time.Second)

	DB = db
	return nil
}
