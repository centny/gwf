package routing

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Centny/gwf/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type MultipartFile struct {
	Name     string
	Filename string
	SavePath string
	Sha1     []byte
	Md5      []byte
	Length   int64
}
type MultipartValues struct {
	Files  []MultipartFile
	Values map[string][][]byte
}

func (h *HTTPSession) RecMultipart(sha bool, save_path_f func(*multipart.Part) string) (*MultipartValues, error) {
	mr, err := h.R.MultipartReader()
	if err != nil {
		return nil, util.Err("MultipartReader err(%v)", err.Error())
	}
	vals := &MultipartValues{
		Files:  []MultipartFile{},
		Values: map[string][][]byte{},
	}
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return vals, util.Err("NextPart err(%v)", err.Error())
		}
		if len(part.FileName()) < 1 {
			bys, err := ioutil.ReadAll(part)
			if err != nil {
				part.Close()
				return vals, err
			}
			part.Close()
			vals.Values[part.FormName()] = append(vals.Values[part.FormName()], bys)
			continue
		}
		tfile := save_path_f(part)
		if len(tfile) < 1 {
			part.Close()
			continue
		}
		_, fn := filepath.Split(part.FileName())
		if strings.HasSuffix(tfile, "/") {
			tfile = tfile + fn
		}
		var w int64
		var sha_, md5_ []byte
		if sha {
			w, sha_, md5_, err = util.Copyp2(tfile, part)
		} else {
			w, err = util.Copyp(tfile, part)
		}
		if err != nil {
			part.Close()
			return vals, err
		}
		part.Close()
		vals.Files = append(vals.Files, MultipartFile{
			Name:     part.FormName(),
			Filename: part.FileName(),
			SavePath: tfile,
			Sha1:     sha_,
			Md5:      md5_,
			Length:   w,
		})
	}
	return vals, nil
}

func (h *HTTPSession) RecF(name, tfile string) (int64, error) {
	_, w, _, _, err := h.RecFvV(false, name, func(part *multipart.Part) string {
		return tfile
	})
	return w, err
}
func (h *HTTPSession) RecBys(name string, max int) ([]byte, error) {
	if max < 100 {
		max = 100
	}
	src, _, err := h.R.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer src.Close()
	dst_b := bytes.NewBuffer(nil)
	dst := bufio.NewWriterSize(dst_b, max)
	_, err = io.Copy(dst, src)
	if err == nil {
		return dst_b.Bytes(), nil
	} else {
		return nil, err
	}
}

func (h *HTTPSession) RecF2(name, tfile string) (w int64, sha_ string, md5_ string, err error) {
	_, w, sha_, md5_, err = h.RecFv2(name, tfile)
	return
}
func (h *HTTPSession) RecFv(name, tfile string) (w int64, sha_ []byte, md5_ []byte, err error) {
	_, w, sha_, md5_, err = h.RecFvN(name, tfile)
	return
}
func (h *HTTPSession) RecFv2(name, tfile string) (fn string, w int64, sha_ string, md5_ string, err error) {
	fn, w, sh, md, err := h.RecFvN(name, tfile)
	return fn, w, fmt.Sprintf("%x", sh), fmt.Sprintf("%x", md), err
}
func (h *HTTPSession) RecFvN(name, tfile string) (fn string, w int64, sha_ []byte, md5_ []byte, err error) {
	return h.RecFvV(true, name, func(*multipart.Part) string {
		return tfile
	})
}
func (h *HTTPSession) RecFvV(sha bool, name string, tfile_f func(*multipart.Part) string) (fn string, w int64, sha_ []byte, md5_ []byte, err error) {
	res, err := h.RecMultipart(sha, func(part *multipart.Part) string {
		if name == part.FormName() {
			return tfile_f(part)
		} else {
			return ""
		}
	})
	if err != nil {
		return "", 0, nil, nil, err
	}
	if len(res.Files) < 1 {
		return "", 0, nil, nil, util.NOT_FOUND
	}
	tf := res.Files[0]
	return tf.Filename, tf.Length, tf.Sha1, tf.Md5, nil
}

func (h *HTTPSession) RecFvV2(name string, tfile_f func(*multipart.Part) string) (fn string, w int64, sha_ string, md5_ string, err error) {
	fn, w, sh, md, err := h.RecFvV(true, name, tfile_f)
	return fn, w, fmt.Sprintf("%x", sh), fmt.Sprintf("%x", md), err
}
