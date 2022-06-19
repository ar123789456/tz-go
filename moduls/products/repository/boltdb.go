package repository

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"log"
	"tz/models"
	"tz/moduls/products"

	"github.com/boltdb/bolt"
)

type Repository struct {
	db *bolt.DB
}

func NewRepository(db *bolt.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r *Repository) GetAll(c context.Context) ([]models.Product, error) {
	var posts []models.Product
	err := r.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("products"))

		err := b.ForEach(func(k, v []byte) error {
			buff := bytes.NewBuffer(v)
			p, err := decode(*buff)
			if err != nil {
				log.Println(err)
				return err
			}
			posts = append(posts, p)
			return nil
		})
		if err != nil {
			log.Println(err)
			return err
		}
		return err
	})
	return posts, err
}

func (r *Repository) Get(c context.Context, id int) (models.Product, error) {
	var prod models.Product
	err := r.db.View(func(t *bolt.Tx) error {
		b2 := t.Bucket([]byte("products"))
		buf := b2.Get(itob(id))

		var err error
		prod, err = decode(*bytes.NewBuffer(buf))
		return err
	})
	return prod, err
}

func (r *Repository) Post(c context.Context, p models.Product) (int, error) {
	err := r.db.Update(func(t *bolt.Tx) error {
		b2 := t.Bucket([]byte("name"))
		id := b2.Get([]byte(p.Name))
		if len(id) != 0 {
			log.Println(products.ErrInvalidName)
			return products.ErrInvalidName
		}
		if p.Price < 0 {
			log.Println(products.ErrInvalidPrice)
			return products.ErrInvalidPrice
		}
		b := t.Bucket([]byte("products"))
		id64, err := b.NextSequence()
		if err != nil {
			log.Println(err)
			return err
		}
		p.Id = int(id64)
		buff, err := encode(p)
		if err != nil {
			log.Println(err)
			return err
		}
		err = b.Put(itob(p.Id), buff.Bytes())
		if err != nil {
			log.Println(err)
			return err
		}
		return b2.Put([]byte(p.Name), itob(p.Id))
	})
	log.Println(err)
	return p.Id, err
}

func (r *Repository) Delete(c context.Context, i int) error {
	return r.db.Update(func(t *bolt.Tx) error {
		b1 := t.Bucket([]byte("products"))
		val := b1.Get(itob(i))

		p, err := decode(*bytes.NewBuffer(val))
		if err != nil {
			return err
		}

		err = b1.Delete(itob(i))
		if err != nil {
			return err
		}
		b2 := t.Bucket([]byte("name"))
		return b2.Delete([]byte(p.Name))
	})
}

func (r *Repository) Update(c context.Context, p models.Product) error {
	return r.db.Update(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("products"))
		b2 := t.Bucket([]byte("name"))

		buf := b.Get(itob(p.Id))
		prod, err := decode(*bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
		buf = b2.Get([]byte(prod.Name))
		if !bytes.Equal(buf, itob(p.Id)) {
			return errors.New("invalid id")
		}
		err = b.Delete(itob(p.Id))
		if err != nil {
			return err
		}
		err = b2.Delete([]byte(prod.Name))
		if err != nil {
			return err
		}

		buff, err := encode(p)
		if err != nil {
			return err
		}
		err = b.Put(itob(p.Id), buff.Bytes())
		if err != nil {
			return err
		}
		err = b2.Put([]byte(p.Name), itob(p.Id))
		return err
	})
}

func (r *Repository) Find(c context.Context, name string) (models.Product, error) {
	var prod models.Product
	err := r.db.View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte("name"))
		id := b.Get([]byte(name))
		log.Println(id)

		b2 := t.Bucket([]byte("products"))
		buf := b2.Get(id)

		var err error
		prod, err = decode(*bytes.NewBuffer(buf))
		return err
	})
	return prod, err
}

func decode(buff bytes.Buffer) (models.Product, error) {
	out := models.Product{}
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&out)
	log.Println(out)
	return out, err
}

func encode(p models.Product) (bytes.Buffer, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	return buf, err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
