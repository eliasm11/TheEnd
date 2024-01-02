package main

import (
	"db"
	"srever"
)

func main() {
	db.OpenDirDataBase("DataBase")
	srever.InitSever().Run()

	db.MainDB.Close()

}
