package db

import (
	"database/sql"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewConnection(connectionString string) (*sql.DB, error) {  
	//parse the string 
	
	db, err := sql.Open("pgx", connectionString) 
	if err != nil {
		return nil, err
	}
	
	var lastErr error
	for i := 0; i < 5; i++ {
		if err := db.Ping(); err == nil {
			return db, nil
		}
		lastErr = err
		if i < 4 { 
			time.Sleep(time.Second * 2)
		}
	}
	return nil, lastErr
}
