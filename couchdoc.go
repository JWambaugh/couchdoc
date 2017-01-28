package couchdoc

import (
	"reflect"

	"gopkg.in/couchbaselabs/gocb.v1"
)

var bucket *gocb.Bucket

func SetBucket(b *gocb.Bucket) {
	bucket = b
}

func Doc(doc interface{}) *DocManipulator {
	if reflect.ValueOf(doc).Kind() != reflect.Ptr {
		panic("You must pass a pointer to couchdoc.Doc! Did you forget the '&'? Example: couchdoc.Doc(&myDoc) ")
	}
	man := new(DocManipulator)
	man.doc = doc
	return man
}

func Find(key string, res interface{}) (gocb.Cas, error) {
	return bucket.Get(key, res)
}
