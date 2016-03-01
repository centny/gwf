package util

import (
	"encoding/base64"
	"encoding/binary"
	"hash/crc32"
)

func ShortLink(v string) string {
	uv := crc32.ChecksumIEEE([]byte(v))
	bv := make([]byte, 4)
	binary.BigEndian.PutUint32(bv, uv)
	return base64.NewEncoding("1234567890-=qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM").EncodeToString(bv)
}
