package gochan

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/garyburd/redigo/redis"
	"github.com/youtube/vitess/go/pools"
	"net/http"
	"time"
)

const (
	REDIS_POOL_SIZE = 20
	REDIS_POOL_INIT = 5
	REDIS_POOL_TIME = time.Minute
)

var (
	redisAddr = ":6379"
)

func randBytes(n int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	return bytes
}

func GenerateUid(r *http.Request) string {
	addr := r.RemoteAddr
	hash := sha256.New()
	hash.Sum([]byte(addr))
	hash.Sum(randBytes(100))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

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
