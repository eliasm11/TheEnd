package db

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"slices"
	"strings"
	"structs"

	"go.etcd.io/bbolt"
)

var (
	ErrDataBase = errors.New("data base err")

	ErrNotFoundInStock            = errors.New("not found in Stock")
	ErrContainerIsExistInStock    = errors.New("Container  is exist")
	ErrContainerIsNotExistInStock = errors.New("Container is not exist")
	ErrContainerIsEmpty           = errors.New("Name of Container is Empty")
	ErrKindNotFound               = errors.New("cant found the kind")
	ErrKindIsExited               = errors.New(" kind is exited ")
	ErrKindIsEmpty                = errors.New("Name of Kind is Empty")

	ErrModelSize         = errors.New("model size is empty or Colors is empty")
	ErrModelPrice        = errors.New("Price is Not acceptable")
	ErrModelDescription  = errors.New("Description is Not acceptable")
	ErrModelLinkesImage  = errors.New("linkes Image is empty")
	ErrModelSizeNotFound = errors.New("Size not found")
	ErrModelIdNotFound   = errors.New("Id of Product not Found")

	ErrModelCommintEmpty = errors.New("Commint is not accptable")
	ErrModelCommintStars = errors.New("Stars is not accptable")
)

type stock struct {
	pathStock string
	database  map[string]map[string]*bbolt.DB
}

func (s *stock) DeleteKind(container, kind string) error {

	if len(container) == 0 {
		return ErrContainerIsEmpty
	}
	if len(kind) == 0 {
		return ErrKindIsEmpty
	}
	v, ok := s.database[container]
	if !ok {
		return ErrContainerIsNotExistInStock
	}
	kind1, ok1 := v[kind]
	if !ok1 {
		return ErrKindNotFound
	}
	deleteKindPath := kind1.Path()
	kind1.Close()
	os.Remove(deleteKindPath)
	delete(s.database[container], kind)

	return nil
}

func (s *stock) DeleteCommint(id int, Container, kind, username, commint string) error {

	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return errors.New("no model has that id")
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}
		model.Commint = slices.DeleteFunc(model.Commint, func(s structs.UserCommint) bool {
			return (s.Username == username && s.Commint == commint)
		})

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) addCommint(id int, Container, kind string, commint structs.UserCommint) error {
	if len(commint.Commint) == 0 {
		return ErrModelCommintEmpty
	}
	if commint.Stars <= 0 {
		return ErrModelCommintStars
	}

	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return errors.New("no model has that id")
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.Commint = append(model.Commint, commint)

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpateModel(id int, Container, kind string, model *structs.Model) error {

	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	if len(model.Sizes) == 0 {
		return ErrModelSize
	}
	for _, v := range model.Sizes {
		if len(v) == 0 {
			return ErrModelSize
		} else {
			for _, s := range v {
				if s == 0 {
					return ErrModelSize
				}
			}
		}
	}

	if model.Price <= 0 {
		return ErrModelPrice
	}
	if len(model.Description) == 0 {
		return ErrModelDescription
	}
	if len(model.LinkesImage) == 0 {
		return ErrModelLinkesImage
	}
	return v.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		iserror := b.Get(itob(id))
		if iserror == nil {
			return ErrModelIdNotFound
		}
		model.Id = id
		data, err := json.Marshal(model)
		if err != nil {
			return err
		}
		return b.Put(itob(id), data)
	})
}

func (s *stock) UpateLinkesImage(id int, Container, kind string, linkesImage []string) error {
	if len(linkesImage) == 0 {
		return ErrModelLinkesImage
	}
	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrModelIdNotFound
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.LinkesImage = linkesImage

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpateDescription(id int, Description, Container, kind string) error {
	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.Description = Description

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpatePrice(id int, Price float32, Container, kind string) error {
	if Price <= 0 {
		return ErrModelPrice
	}
	v, err := s.ContainerAndKindIsExited(Container, kind)
	if err != nil {
		return err
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}

		model.Price = Price

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

/*
	func (s *stock) GetALLSize(id int, Container, kind, sizeName string) (*map[string][]map[string]int, error) {
		Container = strings.TrimSpace(Container)
		kind = strings.TrimSpace(kind)
		v, ok := s.database[Container][kind]
		if !ok {
			return nil, errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
		}

		Size := make(map[string][]structs.Size)
		err := v.View(func(tx *bbolt.Tx) error {

			b := tx.Bucket([]byte(kind))
			if b == nil {
				return ErrKindNotFound
			}
			return b.ForEach(func(k, v []byte) error {
				model := structs.Model{}
				if err := json.Unmarshal(v, &model); err != nil {
					return err
				}
				for k, v := range model.Sizes {
					Size[k] = append(Size[k], v)
				}

				return nil
			})
		})

		if err != nil {
			return nil, err
		}

		return &Size, nil
	}

	func (s *stock) GetModleSize(id int, Container, kind, sizeName string) (*structs.Size, error) {
		Container = strings.TrimSpace(Container)
		kind = strings.TrimSpace(kind)
		v, ok := s.database[Container][kind]
		if !ok {
			return nil, errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
		}
		var size structs.Size
		err := v.View(func(tx *bbolt.Tx) error {

			b := tx.Bucket([]byte(kind))
			if b == nil {
				return ErrKindNotFound
			}

			data := b.Get(itob(id))
			if data == nil {
				return ErrDataBase
			}

			model := structs.Model{}

			err := json.Unmarshal(data, &model)
			if err != nil {
				return err
			}
			v, ok := model.Sizes[sizeName]
			if !ok {
				return ErrModelSizeNotFound
			}
			size = v
			return nil
		})
		if err != nil {
			return nil, err
		}

		return &size, nil
	}
*/
func (s *stock) DeleteSizeFromModel(id int, Container, kind, sizeName string) error {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}
		delete(model.Sizes, sizeName)

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) UpdataSizeFromModel(id int, Container, kind, sizeName string, size map[string]int) error {

	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}

	return v.Batch(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(kind))

		if b == nil {
			return ErrKindNotFound
		}

		data := b.Get(itob(id))
		if data == nil {
			return ErrDataBase
		}

		model := structs.Model{}

		err := json.Unmarshal(data, &model)
		if err != nil {
			return err
		}
		if model.Sizes[sizeName] == nil {
			model.Sizes[sizeName] = make(map[string]int)
		}
		model.Sizes[sizeName] = size

		dataTodatabase, err := json.Marshal(model)
		if err != nil {
			return err
		}

		return b.Put(itob(model.Id), dataTodatabase)
	})

}

func (s *stock) AddModelToKind(Container, kind string, model *structs.Model) error {

	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	if len(model.Sizes) == 0 {
		return ErrModelSize
	}
	for _, v := range model.Sizes {
		if len(v) == 0 {
			return ErrModelSize
		} else {
			for _, s := range v {
				if s == 0 {
					return ErrModelSize
				}
			}
		}
	}

	if model.Price <= 0 {
		return ErrModelPrice
	}
	if len(model.Description) == 0 {
		return ErrModelDescription
	}
	if len(model.LinkesImage) == 0 {
		return ErrModelLinkesImage
	}
	return v.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}

		c, err := b.NextSequence()
		if err != nil {
			// this not to clinc
			return err
		}
		model.Id = int(c)

		data, err := json.Marshal(&model)
		if err != nil {
			// this not to clinc
			return err
		}
		return b.Put(itob(model.Id), data)
	})

}
func (s *stock) NewKindtoContainer(Container, kind string) error {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	if len(Container) == 0 {
		return ErrContainerIsEmpty
	}
	if len(kind) == 0 {
		return ErrKindIsEmpty
	}
	_, ok := s.database[Container][kind]
	if ok {
		return ErrKindIsExited
	}
	if _, ok := s.database[Container]; !ok {
		return ErrContainerIsNotExistInStock
	}

	b, err := bbolt.Open(s.pathStock+"/"+Container+"/"+kind+".db", 0600, nil)
	if err != nil {
		return err
	}

	if err = b.Batch(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucket([]byte(kind))
		return err
	}); err != nil {
		return err
	}
	if s.database[Container] == nil {
		s.database[Container] = make(map[string]*bbolt.DB)
	}
	s.database[Container][kind] = b

	return nil
}

func (s *stock) DeleteModelFromKind(Container, kind string, id int) error {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	return v.Batch(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}
		if err := b.Get(itob(id)); err == nil {
			return errors.New("ID Not found")
		}

		b.Delete(itob(id))
		return nil
	})

}

// just for admin
// return all models  in kind  at Container
func (s *stock) GetAllModelsInKind(Container, kind string) (models []structs.Model, _ error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return models, ErrContainerIsNotExistInStock
	}
	err := v.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		// check if bucket is exit
		if b == nil {
			return ErrKindNotFound
		}

		return b.ForEach(func(_, v []byte) error {
			getModel := structs.Model{}
			if err := json.Unmarshal(v, &getModel); err != nil {
				return err
			}
			models = append(models, getModel)
			return nil
		})

	})
	if err != nil {
		return models, nil
	}

	return models, nil
}

func (s *stock) GetNumberModelsInKind(Container, kind string) (int, error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return 0, ErrContainerIsNotExistInStock
	}
	number := 0

	err := v.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}
		number = b.Stats().KeyN

		return nil
	})
	return number, err

}

func (s *stock) GetModelsInKind(formId, count int, Container, kind string) (models []structs.Model, _ error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return models, ErrContainerIsNotExistInStock
	}

	err := v.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		if b == nil {
			return ErrKindNotFound
		}
		c := b.Cursor()

		k, nextVuale := c.Seek(itob(formId))
		if k == nil {
			return ErrKindNotFound
		}
		model := structs.Model{}
		err := json.Unmarshal(nextVuale, &model)
		if err != nil {
			return err
		}
		models = append(models, model)

		k, nextVuale = c.Prev()
		i := 0

		for k != nil && i < count {
			model := structs.Model{}
			err := json.Unmarshal(nextVuale, &model)
			if err != nil {
				return err
			}

			models = append(models, model)

			k, nextVuale = c.Prev()
			i++
		}

		return nil
	})
	if err != nil {
		return models, nil
	}

	return models, nil
}

func (S *stock) GetAllContainer() (Container []string) {
	for k := range S.database {
		Container = append(Container, k)
	}
	return Container
}

func (S *stock) GetAllContainerAndKind() (ContainerAndKind map[string][]string) {
	ContainerAndKind = make(map[string][]string)
	for k1, v := range S.database {
		ContainerAndKind[k1] = make([]string, 0)
		for k := range v {
			ContainerAndKind[k1] = append(ContainerAndKind[k1], k)
		}
	}
	return ContainerAndKind
}

func (db *stock) IsExist(name string) bool {
	name = strings.TrimSpace(name)
	_, ok := db.database[name]
	return ok
}
func (db *stock) DeleteContainer(name string) error {
	if len(name) == 0 {
		return ErrContainerIsEmpty
	}
	v, ok := db.database[name]
	if !ok {
		return ErrContainerIsNotExistInStock
	}
	for _, database := range v {
		database.Close()
	}
	delete(db.database, name)

	err := os.RemoveAll(db.pathStock + "/" + name)
	if err != nil {
		return err
	}
	return nil
}

func (db *stock) AddNewContainer(name string) error {

	if len(name) == 0 {
		return ErrContainerIsEmpty
	}
	name = strings.TrimSpace(name)
	if db.IsExist(name) {
		return ErrContainerIsExistInStock
	}
	db.database[name] = nil
	err := os.Mkdir(db.pathStock+"/"+name, 0770)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func (db *stock) init() error {

	f, err := os.Open(db.pathStock)
	defer f.Close()
	if err != nil {
		return err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, v := range files {
		if v.IsDir() {
			db.database[v.Name()] = make(map[string]*bbolt.DB)

			SubDir, err := os.Open(db.pathStock + "/" + v.Name())

			if err != nil {
				return err
			}
			defer SubDir.Close()
			filesInSubDir, _ := SubDir.Readdir(0)
			if filesInSubDir == nil {
				continue
			}
			for _, file := range filesInSubDir {

				if !file.IsDir() && path.Ext(file.Name()) == ".db" {
					lenCut := len(file.Name()) - 3
					ptrdb, err := bbolt.Open(db.pathStock+"/"+v.Name()+"/"+file.Name(), 0600, nil)
					if err != nil {
						return err
					}
					db.database[v.Name()][file.Name()[:lenCut]] = ptrdb

				}

			}

		}
	}

	return nil
}

func (s *stock) ContainerAndKindIsExited(Container, kind string) (*bbolt.DB, error) {
	Container = strings.TrimSpace(Container)
	kind = strings.TrimSpace(kind)
	v, ok := s.database[Container][kind]
	if !ok {
		return nil, errors.Join(ErrKindNotFound, ErrContainerIsNotExistInStock)
	}
	return v, nil
}
