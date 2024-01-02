package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"go.etcd.io/bbolt"
)

var (
	ErrOrdersUsername  = errors.New("username is empty")
	ErrOrdersColor     = errors.New("color is empty")
	ErrOrdersSize      = errors.New("size is empty")
	ErrOrdersId        = errors.New("id is zero or less then zore")
	ErrOrdersContainer = errors.New("Container is empty")
	ErrOrdersType      = errors.New("Type is empty")
	ErrOrderNotFound   = errors.New("Order id not found")
)

type Orders struct {
	Id        uint   `json:"id"`
	IdModel   int    `json:"idModel"`
	Color     string `json:"color"`
	SizeName  string `json:"size"`
	Container string `json:"Container"`
	Kind      string `json:"Kind"`
	Username  string `json:"username"`
	UserAddr  string `json:"addr"`
	UserPhone string `json:"phone"`

}

type order struct {
	dataBase *bbolt.DB
}

func (o *order) Delete(id int) error {

	return o.dataBase.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("orders"))
		if b == nil {
			return ErrOrderNotFound
		}
		return b.Delete(itob(id))
	})
}

func (o *order) Add(order Orders ) error {
	return o.dataBase.Batch(func(tx *bbolt.Tx) error {
		if order.SizeName == "" {
			return ErrOrdersSize
		}
		if order.Color == "" {
			return ErrOrdersColor
		}
		if order.IdModel <= 0 {
			return ErrOrdersId
		}
		if order.Username == "" {
			return ErrOrdersUsername
		}

		if order.Container == "" {
			return ErrOrdersContainer
		}

		if order.Kind == "" {
			return ErrOrdersType
		}
		b := tx.Bucket([]byte("orders"))

		id, _ := b.NextSequence()
		idnext := int(id)
		order.Id = uint(idnext)
		data, err := json.Marshal(order)
		if err != nil {
			return err
		}
		b.Put(itob(idnext), data)
		return nil
	})

}

func (o *order) Get() []string {
	var allorder []string
	o.dataBase.View(func(tx *bbolt.Tx) error {
		next := tx.Bucket([]byte("orders")).Cursor()
		for _, v := next.First(); v != nil; _, v = next.Next() {
			allorder = append(allorder, string(v))
		}
		return nil
	})

	return allorder
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
