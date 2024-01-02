package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"structs"
	"go.etcd.io/bbolt"
)

type DataBase struct {
	NameDataBase string
	Stock        stock
	Orders       order
	Users        user
	OutStock     OutStock
}

var MainDB *DataBase

var (
	ErrDataBaseOutStock = errors.New("Sorry, out of stock")
)

func (db *DataBase) AddCommint(id int, Container, kind string, commint structs.UserCommint) error {
	err := db.Users.GetUser(commint.Username, nil)
	if err == ErrUsereNotFound {
		return ErrUsereNotFound
	}
	return db.Stock.addCommint(id, Container, kind, commint)
}

func (db *DataBase) Buy(order Orders) error {
	models, err := db.Stock.GetModelsInKind(order.IdModel, 0, order.Container, order.Kind)
	if err != nil {
		return err
	}
	if len(models) == 0 {
		return errors.New("no data")
	}
	user := structs.User{}
	err = db.Users.GetUser(order.Username, &user)

	if err == ErrUsereNotFound {
		return err
	}

	if user.Phone == "" {
		return errors.New("updata your number")
	}
	if user.Name == "" {
		return errors.New("updata your Name")
	}
	if user.UserAddr.City == "" {

		return errors.New("updata your address")
	}
	if len(user.UserVisa.Number) == 0 {
		return errors.New("add visa")
	}

	if _, ok := models[0].Sizes[order.SizeName]; !ok {
		return ErrModelSizeNotFound
	}
	order.Username = user.Name
	order.UserAddr = fmt.Sprintf("%s/%s/%s" , user.UserAddr.City , user.UserAddr.Street ,  user.UserAddr.Area)
	order.UserPhone = user.Phone  
	db.Orders.Add(order)
	models[0].Sizes[order.SizeName][order.Color]--
	if models[0].Sizes[order.SizeName][order.Color] <= 0 {
		db.OutStock.add(&modelsOutStock{IdModel: order.IdModel, Container: order.Container, Kind: order.Kind, Size: order.SizeName, Color: order.Color})
		delete(models[0].Sizes[order.SizeName], order.Color)
		db.Stock.UpdataSizeFromModel(order.IdModel, order.Container, order.Kind, order.SizeName, models[0].Sizes[order.SizeName])
	} else {
		db.Stock.UpdataSizeFromModel(order.IdModel, order.Container, order.Kind, order.SizeName, models[0].Sizes[order.SizeName])
	}
	if len(models[0].Sizes[order.SizeName]) == 0 {
		db.Stock.DeleteSizeFromModel(order.IdModel, order.Container, order.Kind, order.SizeName)
	}

	return nil
}

func (db *DataBase) Close() {
	for _, v1 := range db.Stock.database {
		for _, data := range v1 {
			data.Close()
		}
	}
	db.Orders.dataBase.Close()
	db.Users.dataBase.Close()
}

func OpenDirDataBase(name string) *DataBase {

	db := DataBase{}

	db.Stock.pathStock = name + "/Stock"
	err := os.Mkdir(name, 0770)
	if err != nil && !os.IsExist(err) {
		log.Fatal("[database] " + err.Error())
	}
	err = os.Mkdir(db.Stock.pathStock, 0770)
	if err != nil && !os.IsExist(err) {
		log.Fatal("[database] " + err.Error())
	}

	if err != nil && !os.IsExist(err) {
		log.Fatal("[database] " + err.Error())
	}

	db.Stock.database = make(map[string]map[string]*bbolt.DB)

	db.NameDataBase = name

	if err := db.Stock.init(); err != nil {
		log.Fatal(err)
	}

	db.Users.dataBase, err = bbolt.Open(name+"/user"+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.OutStock.dataBase, err = bbolt.Open(name+"/outstock"+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.OutStock.dataBase.Batch(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("outstock"))
		return err
	}); err != nil {
		log.Fatal(err)
	}
	//db.OutStock.dataBase,.name = ""
	db.Users.dataBase.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	// init Orders
	db.Orders.dataBase, err = bbolt.Open(name+"/orders"+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Orders.dataBase.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("orders"))
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})
	MainDB = &db
	return &db
}
