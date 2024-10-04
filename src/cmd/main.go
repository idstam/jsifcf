package main

import (
	"fmt"
	"jsifcf/internal"
)

func main() {
	fmt.Println("asf")

	db := &internal.SqliteDB{}
	logger := internal.Logger{}
	logger.Init(internal.Trace)
	db.Init(":memory:", logger)

	session, _ := db.AddSession("foobar")
	internal.ScanPath(db, session, "/home/src", internal.MD5)

}
