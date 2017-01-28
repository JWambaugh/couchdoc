package couchdoc

import gocb "gopkg.in/couchbaselabs/gocb.v1"

type Document struct {
	Id      string   `json:"-"`
	Expire_ uint32   `json:"-"`
	Cas_    gocb.Cas `json:"-"`
}
