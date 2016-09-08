package util

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
)

//Crc32 do hash for byte data by crc32 and return the bas64 string
func Crc32(v []byte) string {
	uv := crc32.ChecksumIEEE(v)
	bv := make([]byte, 4)
	binary.BigEndian.PutUint32(bv, uv)
	return base64.StdEncoding.EncodeToString(bv)
}

//Sha1 do hash for file by sha1 and return the base64 string
func Sha1(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var hash = sha1.New()
	_, err = bufio.NewReader(f).WriteTo(hash)
	return fmt.Sprintf("%x", hash.Sum(nil)), err
}

//Md5 do hash for file by md5 and return the base64 string
func Md5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	var hash = md5.New()
	_, err = bufio.NewReader(f).WriteTo(hash)
	return fmt.Sprintf("%x", hash.Sum(nil)), err
}

//Md5Byte do hash for bytes by md5 and return the base64 string
func Md5Byte(bys []byte) string {
	var hash = md5.New()
	hash.Write(bys)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//Sha1Byte do has for bytes by sha1 and return the base64 string
func Sha1Byte(bys []byte) string {
	var hash = sha1.New()
	hash.Write(bys)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//ShortLink do short link for string
func ShortLink(v string) string {
	uv := crc32.ChecksumIEEE([]byte(v))
	bv := make([]byte, 4)
	binary.BigEndian.PutUint32(bv, uv)
	return base64.NewEncoding("1234567890-=qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM").EncodeToString(bv)
}
