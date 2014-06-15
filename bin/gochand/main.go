package main

import (
	"bytes"
	"encoding/gob"
	"github.com/cydev/gochan"
	"log"
	"net/http"
	"time"
)

func main() {
	p := gochan.GetPool()
	defer p.Close()
	r, err := p.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer p.Put(r)
	c := r.(gochan.ResourceConn)
	n, err := c.Do("INFO")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info=%s", n)

	post := &gochan.Post{}
	post.Date = time.Now()
	post.Text = "test text"
	post.User = gochan.GenerateUid(&http.Request{})
	var buff bytes.Buffer
	b := gob.NewEncoder(&buff)
	be := gob.NewDecoder(&buff)
	err = b.Encode(post)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(buff.Bytes()))
	log.Println(len(buff.Bytes()))
	newpost := &gochan.Post{}
	err = be.Decode(newpost)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(post.Date, post.Text, post.User)
}
