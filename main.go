package main

import (
	"bytes"
	"encoding/gob"
	"github.com/garyburd/redigo/redis"
	"github.com/youtube/vitess/go/pools"
	"log"
	"time"
)

type Post struct {
	Date time.Time
	Text []byte
}

const (
	REDIS_POOL_SIZE = 20
	REDIS_POOL_INIT = 5
	REDIS_POOL_TIME = time.Minute
)

var (
	redisAddr = ":6379"
)

// ResourceConn adapts a Redigo connection to a Vitess Resource.
type ResourceConn struct {
	redis.Conn
}

func (r ResourceConn) Close() {
	r.Conn.Close()
}

func GetPool() *pools.ResourcePool {
	return pools.NewResourcePool(func() (pools.Resource, error) {
		c, err := redis.Dial("tcp", redisAddr)
		return ResourceConn{c}, err
	}, REDIS_POOL_INIT, REDIS_POOL_SIZE, REDIS_POOL_TIME)
}

func main() {
	p := GetPool()
	defer p.Close()
	r, err := p.Get()
	if err != nil {
		log.Fatal(err)
	}
	defer p.Put(r)
	c := r.(ResourceConn)
	n, err := c.Do("INFO")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("info=%s", n)

	post := &Post{time.Now(), []byte("test text")}
	var buff bytes.Buffer
	b := gob.NewEncoder(&buff)
	be := gob.NewDecoder(&buff)
	err = b.Encode(post)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(buff.Bytes()))
	log.Println(len(buff.Bytes()))
	newpost := &Post{}
	err = be.Decode(newpost)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(newpost.Date, string(newpost.Text))
	log.Println(post.Date, string(post.Text))

}
