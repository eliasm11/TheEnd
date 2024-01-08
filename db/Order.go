package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

var (
	ErrOrdersUsername  = errors.New("Username is empty")
	ErrOrdersColor     = errors.New("Color is empty")
	ErrOrdersSize      = errors.New("Size is empty")
	ErrOrdersId        = errors.New("Id is zero or less then zore")
	ErrOrdersContainer = errors.New("Container is empty")
	ErrOrdersType      = errors.New("Type is empty")
	ErrOrderNotFound   = errors.New("Order id not found")
)

/*
Id        uint   `json:"id"`
IdModel   int    `json:"idModel"`
Color     string `json:"color"`
SizeName  string `json:"size"`
*/
type UserOrder struct {
	Username  string `json:"username"`
	UserAddr  string `json:"addr"`
	UserPhone string `json:"phone"`
}
type ItemOrder struct {
	Container string                    `json:"Container"`
	Kind      string                    `json:"Kind"`
	Id        int                       `json:"idModel"`
	Size      map[string]map[string]int `json:"sizes"`
}
type Qty int
type Container string
type Kind string
type Size string
type Color string
type Id uint
type Orders struct {
	Id   uint                                                 `json:"id"`
	User UserOrder                                            `json:"user"`
	Item map[Container]map[Kind]map[Id]map[Size]map[Color]Qty `json:"item"`
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

func (o *order) Add(order *Orders) error {
	return o.dataBase.Batch(func(tx *bbolt.Tx) error {

		if len(order.Item) == 0 {
			return ErrOrdersContainer
		}

		if order.User.Username == "" {
			return ErrOrdersUsername
		}
		for _, kind := range order.Item {
			if len(kind) == 0 {
				return ErrOrdersType
			} else {
				for _, id := range kind {
					if len(id) == 0 {
						return ErrOrdersType
					} else {
						for _, size := range id {
							if len(size) == 0 {
								return ErrOrdersId
							} else {
								for _, color := range size {
									if len(color) == 0 {
										return ErrOrdersColor
									} else {
										for _, qty := range color {
											if qty <= 0 {
												return errors.New("Qty can't be zero or less")
											}
										}
									}

								}
							}
						}
					}
				}
			}
		}

		tim := time.Now()
		b, _ := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%d-%d-%d", tim.Year(), tim.Month(), tim.Day())))

		data := b.Get([]byte(order.User.Username))

		if data == nil {
			id, _ := b.NextSequence()
			idnext := int(id)
			order.Id = uint(idnext)
			data, err := json.Marshal(order)
			if err != nil {
				return err
			}
			b.Put([]byte(order.User.Username), data)
		} else {
			oldOlder := Orders{}
			json.Unmarshal(data, &oldOlder)
			fmt.Println("befor\n", oldOlder)

			for keyContainer, kind := range order.Item {
				if _, ok := oldOlder.Item[keyContainer]; ok == false {
					oldOlder.Item[keyContainer] = order.Item[keyContainer]
				} else {
					for keyKind, id := range kind {
						if _, ok := oldOlder.Item[keyContainer][keyKind]; ok == false {
							oldOlder.Item[keyContainer][keyKind] = order.Item[keyContainer][keyKind]
						} else {
							for keyid, size := range id {
								if _, ok := oldOlder.Item[keyContainer][keyKind][keyid]; ok == false {
									oldOlder.Item[keyContainer][keyKind][keyid] = order.Item[keyContainer][keyKind][keyid]
								} else {
									for keySize, color := range size {
										if _, ok := oldOlder.Item[keyContainer][keyKind][keyid][keySize]; ok == false {
											oldOlder.Item[keyContainer][keyKind][keyid][keySize] = order.Item[keyContainer][keyKind][keyid][keySize]
										} else {
											for keyColor := range color {
												if _, ok := oldOlder.Item[keyContainer][keyKind][keyid][keySize][keyColor]; ok == false {
													oldOlder.Item[keyContainer][keyKind][keyid][keySize][keyColor] = order.Item[keyContainer][keyKind][keyid][keySize][keyColor]
												}else {
													oldOlder.Item[keyContainer][keyKind][keyid][keySize][keyColor] += order.Item[keyContainer][keyKind][keyid][keySize][keyColor]
												}
											}

										}
									}
								}

							}
						}
					}
				}
			}

			data, err := json.Marshal(&oldOlder)
			fmt.Println("after\n", oldOlder)
			if err != nil {
				return err
			}
			b.Put([]byte(order.User.Username), data)
		}

		return nil
	})

}

func (o *order) Get(Time string) []string {
	var allorder []string
	o.dataBase.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(Time))
		if bucket == nil {
			return nil
		}
		next := bucket.Cursor()
		for _, v := next.First(); v != nil; _, v = next.Next() {
			allorder = append(allorder, string(v))
		}
		return nil
	})

	return allorder
}

func (o *order) GetUser(Time, username string, order *Orders) (err error) {
	err = o.dataBase.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(Time))
		if bucket == nil {
			return errors.New("No Data ")
		}
		data := bucket.Get([]byte(username))
		if data == nil {
			return ErrUsereNotFound
		}
		if err := json.Unmarshal(data, order); err != nil {
			return err
		}
		return nil
	})
	return err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

/*

	for keyContainer, kind := range o.Item {
		if _, ok := o.Item[keyContainer]; ok == false {
			o.Item[keyContainer] = order.Item[keyContainer]
		} else {
			for keyKind, id := range kind {
				if _, ok := o.Item[keyContainer][keyKind]; ok == false {
					o.Item[keyContainer][keyKind] = order.Item[keyContainer][keyKind]
				} else {
					for keyId, size := range id {
						if _, ok := o.Item[keyContainer][keyKind][keyId]; ok == false {
							o.Item[keyContainer][keyKind][keyId] = order.Item[keyContainer][keyKind][keyId]
						} else {
							for keySize, color := range size {
								if _, ok := o.Item[keyContainer][keyKind][keyId][keySize]; ok == false {
									o.Item[keyContainer][keyKind][keyId][keySize] = order.Item[keyContainer][keyKind][keyId][keySize]
								} else {
									for KeyColor := range color {
										if _, ok := o.Item[keyContainer][keyKind][keyId][keySize][KeyColor]; ok == false {
											o.Item[keyContainer][keyKind][keyId][keySize][KeyColor] = order.Item[keyContainer][keyKind][keyId][keySize][KeyColor]
										} else {
											o.Item[keyContainer][keyKind][keyId][keySize][KeyColor]++
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
*/
