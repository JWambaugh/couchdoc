package couchdoc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
)

type DocManipulator struct {
	doc interface{}
}

func (man *DocManipulator) Document() interface{} {
	return man.doc
}
func (man *DocManipulator) Save() error {
	key, err := man.GetKey()
	if err != nil {
		return err
	}
	fmt.Println(reflect.TypeOf(man.doc).String())
	b, _ := json.Marshal(man.doc)
	fmt.Printf("%s", b)
	bucket.Upsert(key, man.doc, man.Get("Expire_").(uint32))
	return nil
}

func (man *DocManipulator) Load(id string) error {
	if id != "" {
		if man.Set("Id", id) != nil {
			return errors.New("Cannot set Id! Does " + man.GetStructName() + " embed couchdoc.CouchbaseDocument?")
		}
	}

	key, err := man.GetKey()
	if err != nil {
		return err
	}

	cas, err := bucket.Get(key, man.doc)
	if err != nil {
		return err
	}

	err = man.Set("Cas_", cas)
	if err != nil {
		return err
	}
	return nil
}
func (man *DocManipulator) GetKey() (string, error) {
	err := man.generateId()
	if err != nil {
		return "", err
	}

	key := man.GetStructName()
	id := man.Get("Id").(string)
	if id != "" {
		key += "_" + id
	}

	return key, nil
}
func (man *DocManipulator) generateId() error {
	if man.Get("Override_Auto_Id_Gen") != nil {
		fmt.Println("Not generating key")
		return nil
	}
	if man.Get("Id") == "" {
		//generate new UUID for art piece image
		u, err := uuid.NewV4()
		if err != nil {
			return err
		}
		if man.Set("Id", u.String()) != nil {
			return errors.New("Cannot set Id! Does " + man.GetStructName() + " embed couchdoc.CouchbaseDocument?")
		}
	}
	return nil

}

func (man *DocManipulator) Set(key string, value interface{}) error {
	v := reflect.ValueOf(man.doc).Elem()
	field := v.FieldByName(key)
	if !field.IsValid() {
		return errors.New("Cannot set key '" + key + "' because it doesn't exist.")
	}
	field.Set(reflect.ValueOf(value))
	return nil
}

func (man *DocManipulator) Get(key string) (value interface{}) {
	v := reflect.ValueOf(man.doc).Elem()
	field := v.FieldByName(key)
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}

func (man *DocManipulator) GetStructName() string {
	s := reflect.TypeOf(man.doc).String()
	arr := strings.Split(s, ".")
	return arr[len(arr)-1]
}
