package db

import (
	"encoding/json"
	"log"

	"go.etcd.io/bbolt"
)

type modelsOutStock struct {
	IdModel   int    `json:"idModel"`
	Container string `json:"container"`
	Kind      string `json:"kind"`
	Size      string `json:"size"`
	Color     string `json:"color"`
}

type OutStock struct {
	dataBase *bbolt.DB
}

func (OutStock *OutStock) add(model *modelsOutStock) error {
	return OutStock.dataBase.Batch(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		if bucket = tx.Bucket([]byte("outstock")); bucket == nil {
			log.Panicln("bucket is nil")
		}
		value, _ := json.Marshal(model)
		bucket.Put(itob(int(bucket.Sequence())), value)
		bucket.NextSequence()
		return nil
	})
}

func (OutStock *OutStock) Get() (model []modelsOutStock) {
	OutStock.dataBase.Batch(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		if bucket = tx.Bucket([]byte("outstock")); bucket == nil {
			log.Panicln("bucket is nil")
		}
		bucket.ForEach(func(_, v []byte) error {
			modelTemp := modelsOutStock{}
			json.Unmarshal(v, &modelTemp)
			model = append(model, modelTemp)
			return nil
		})
		return nil
	})
	return model
}
func (OutStock *OutStock) Delete(id int) error {
	return OutStock.dataBase.Batch(func(tx *bbolt.Tx) error {
		var bucket *bbolt.Bucket
		if bucket = tx.Bucket([]byte("outstock")); bucket == nil {
			log.Panicln("bucket is nil")
		}
		bucket.Delete(itob(id))
		return nil
	})

}
