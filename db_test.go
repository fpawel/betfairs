package main

import "testing"

func TestOpenDB(t *testing.T){
	db := openDB()
	dbAddEventsIDs(db, []int{1,2,3,4,5})
	db.Close()
}
