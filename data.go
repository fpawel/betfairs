package main

import (
	_ "github.com/lib/pq"
	"os"
	"fmt"
	"github.com/jmoiron/sqlx"
)



func openDB() *sqlx.DB{
	db := sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))

	db.MustExec(`
CREATE TABLE IF NOT EXISTS events (
  event_id INTEGER UNIQUE NOT NULL CHECK (event_id > 0),
  event_date DATE,
  competition TEXT,
  PRIMARY KEY (event_id)
);
CREATE TABLE IF NOT EXISTS markets (
  market_id INTEGER UNIQUE NOT NULL CHECK (market_id > 0),
  market_name TEXT,
  event_id INT NOT NULL CHECK (event_id > 0),
  FOREIGN KEY (event_id)
  REFERENCES events (event_id) ON DELETE CASCADE,
  PRIMARY KEY (market_id)
);
CREATE TABLE IF NOT EXISTS runners (
  runner_id INTEGER NOT NULL,
  runner_name TEXT,
  market_id INT NOT NULL CHECK (market_id > 0),
  FOREIGN KEY (market_id)
  REFERENCES markets (market_id) ON DELETE CASCADE,
  PRIMARY KEY (market_id, runner_id)
);
DELETE FROM events WHERE event_date < CURRENT_DATE - '7 day'::interval;`)
	return db
}

func dbAddEventsIDs(db *sqlx.DB, eventIDs []int){
	sqlStr := "INSERT INTO events(event_id, event_date) VALUES "
	vals := []interface{}{}

	for _, eventID := range eventIDs {
		sqlStr += fmt.Sprintf("(%d, CURRENT_DATE),", eventID)
		vals = append(vals, eventID)
	}
	sqlStr = sqlStr[0:len(sqlStr)-1]
	sqlStr += " ON CONFLICT (event_id) DO NOTHING "

	db.MustExec(sqlStr)
}

func dbAddEventCompetiton(db *sqlx.DB, eventID int, competiton string){

	db.MustExec(`
INSERT INTO events (event_id, event_date, competition)
    VALUES ($1, CURRENT_DATE, $2)
ON CONFLICT (event_id)
DO UPDATE SET competition = $2`, eventID, competiton)
}
