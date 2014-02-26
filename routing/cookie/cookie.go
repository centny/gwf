package cookie

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/Centny/Cny4go/log"
	"github.com/Centny/Cny4go/routing"
	"net/http"
)

type CookieSession struct {
	W       http.ResponseWriter
	R       *http.Request
	Sb      *CookieSessionBuilder
	kvs     map[string]interface{}
	updated bool
}

func (c *CookieSession) Val(key string) interface{} {
	if v, ok := c.kvs[key]; ok {
		return v
	} else {
		return nil
	}
}
func (c *CookieSession) Set(key string, val interface{}) {
	if val == nil {
		delete(c.kvs, key)
	} else {
		c.kvs[key] = val
	}
	c.updated = true
}
func (c *CookieSession) Flush() error {
	if !c.updated {
		return nil
	}
	val, err := c.Crypto()
	if err != nil {
		return err
	}
	cookie := &http.Cookie{}
	cookie.Name = "C"
	cookie.Domain = c.Sb.Domain
	cookie.Path = c.Sb.Path
	cookie.Value = val
	cookie.MaxAge = 0
	http.SetCookie(c.W, cookie)
	c.updated = false
	return nil
}
func (c *CookieSession) Crypto() (string, error) {
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	// Encoding the map
	err := e.Encode(c.kvs)
	if err != nil {
		return "", err
	}
	bys, err := c.Sb.Crypto(b.Bytes())
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bys), nil
}
func (c *CookieSession) UnCrypto(v string) {
	if len(v) < 1 {
		return
	}
	bys, err := hex.DecodeString(v)
	if err != nil {
		log.D("UnCrypto erro:%s", err.Error())
		return
	}
	ubys, err := c.Sb.UnCrypto(bys)
	if err != nil {
		log.D("UnCrypto erro:%s", err.Error())
		return
	}
	d := gob.NewDecoder(bytes.NewBuffer(ubys))
	// Decoding the serialized data
	err = d.Decode(&c.kvs)
	if err != nil {
		log.D("UnCrypto erro:%s", err.Error())
		return
	}
}

//
type CookieCryptoFunc func(bys []byte) ([]byte, error)

//
type CookieSessionBuilder struct {
	//
	Domain   string
	Path     string
	Crypto   CookieCryptoFunc
	UnCrypto CookieCryptoFunc
}

func NewCookieSessionBuilder(domain string, path string) *CookieSessionBuilder {
	sb := CookieSessionBuilder{}
	sb.Domain = domain
	sb.Path = path
	sb.Crypto = func(bys []byte) ([]byte, error) {
		return bys, nil
	}
	sb.UnCrypto = func(bys []byte) ([]byte, error) {
		return bys, nil
	}
	return &sb
}
func (s *CookieSessionBuilder) FindSession(w http.ResponseWriter, r *http.Request) routing.Session {
	c, err := r.Cookie("C")
	cs := &CookieSession{
		W:       w,
		R:       r,
		Sb:      s,
		kvs:     map[string]interface{}{},
		updated: false,
	}
	if err == nil {
		cs.UnCrypto(c.Value)
	}
	return cs
}
