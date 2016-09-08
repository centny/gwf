package util

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var ShowLog bool = false

func slog(f string, args ...interface{}) {
	if ShowLog {
		fmt.Println(fmt.Sprintf(f, args...))
	}
}

// var C_SH string = "/bin/bash"

func Exec(args ...string) (string, error) {
	return Exec2(strings.Join(args, " "))
}
func Exec2(cmds string) (string, error) {
	var bys []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		bys, err = exec.Command("cmd", "/C", cmds).Output()
	default:
		bys, err = exec.Command("bash", "-c", cmds).Output()
	}
	return string(bys), err
}

func Copy(dst io.Writer, src io.Reader) (written int64, sha_ []byte, md5_ []byte, err error) {
	md5_h, sha_h := md5.New(), sha1.New()
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			md5_h.Write(buf[0:nr])
			sha_h.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	sha_, md5_ = sha_h.Sum(nil), md5_h.Sum(nil)
	return
}
func Copy2(dst io.Writer, src io.Reader) (written int64, sha_ string, md5_ string, err error) {
	w, sh, md, err := Copy(dst, src)
	return w, fmt.Sprintf("%x", sh), fmt.Sprintf("%x", md), err
}
func Copyp(dst string, src io.Reader) (written int64, err error) {
	fp, _ := filepath.Split(dst)
	if !Fexists(fp) {
		os.MkdirAll(fp, os.ModePerm)
	}
	dst_, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer dst_.Close()
	return io.Copy(dst_, src)
}
func Copyp2(dst string, src io.Reader) (written int64, sha_ []byte, md5_ []byte, err error) {
	return Copyp2_(dst, src, os.ModePerm)
}
func Copyp2_(dst string, src io.Reader, mode os.FileMode) (written int64, sha_ []byte, md5_ []byte, err error) {
	fp, _ := filepath.Split(dst)
	if !Fexists(fp) {
		os.MkdirAll(fp, os.ModePerm)
	}
	dst_, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return 0, nil, nil, err
	}
	defer dst_.Close()
	return Copy(dst_, src)
}
func Copyp3(dst string, src io.Reader) (written int64, sha string, md5 string, err error) {
	written, sha_, md5_, err := Copyp2(dst, src)
	return written, fmt.Sprintf("%x", sha_), fmt.Sprintf("%x", md5_), err
}
func Copyp4(dst string, src io.Reader, mode os.FileMode) (written int64, sha string, md5 string, err error) {
	written, sha_, md5_, err := Copyp2_(dst, src, mode)
	return written, fmt.Sprintf("%x", sha_), fmt.Sprintf("%x", md5_), err
}

func ChkVer(n string, o string) (int, error) {
	if len(n) < 1 {
		return 0, Err("new version is empty")
	}
	if len(o) < 1 {
		return 1, nil
	}
	ns := strings.Split(n, ".")
	os := strings.Split(o, ".")
	ml := len(ns)
	if len(os) < ml {
		ml = len(os)
	}
	for i := 0; i < ml; i++ {
		ov, err := strconv.ParseInt(os[i], 10, 64)
		if err != nil {
			return 0, err
		}
		nv, err := strconv.ParseInt(ns[i], 10, 64)
		if err != nil {
			return 0, err
		}
		if nv == ov {
			continue
		} else {
			return int(nv - ov), nil
		}
	}
	return len(ns) - len(os), nil
}
