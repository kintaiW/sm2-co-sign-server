package repository

import (
	"database/sql"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB 初始化数据库连接
func InitDB(dbPath string) error {
	var err error
	once.Do(func() {
		// 确保数据目录存在
		dir := dbPath[:len(dbPath)-len("/cosign.db")]
		if err = os.MkdirAll(dir, 0755); err != nil {
			return
		}

		db, err = sql.Open("sqlite", dbPath)
		if err != nil {
			return
		}

		// 设置连接池
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)

		// 启用 WAL 模式
		_, err = db.Exec("PRAGMA journal_mode = WAL")
		if err != nil {
			return
		}
		_, err = db.Exec("PRAGMA synchronous = NORMAL")
		if err != nil {
			return
		}
		_, err = db.Exec("PRAGMA foreign_keys = ON")
	})
	return err
}

// GetDB 获取数据库连接
func GetDB() *sql.DB {
	return db
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
