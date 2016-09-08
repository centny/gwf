package io

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileSize int64

func (f FileSize) String() string {
	size := int64(f)
	keys := []string{"B", "KB", "MB", "GB", "TB"}
	for i := 0; i < 4; i++ {
		if size < 1024 {
			return fmt.Sprintf("%v%v", size, keys[i])
		} else {
			size = size / 1024
		}
	}
	return fmt.Sprintf("%v%v", size, keys[len(keys)-1])
}

//check if data match regex in list
func MatchRegex(data string, regs ...*regexp.Regexp) bool {
	for _, reg := range regs {
		if reg.MatchString(data) {
			return true
		}
	}
	return false
}

//list all file or folder in root folder by include regex and exclude regex
func Walk(root string, dir bool, inc []string, exc []string, call func(string) string) []string {
	pathes := []string{}
	reg_inc := []*regexp.Regexp{}
	reg_exc := []*regexp.Regexp{}
	for _, i := range inc {
		reg_inc = append(reg_inc, regexp.MustCompile(i))
	}
	for _, e := range exc {
		reg_exc = append(reg_exc, regexp.MustCompile(e))
	}
	add_path := func(path string) {
		if call == nil {
			pathes = append(pathes, path)
		} else {
			pathes = append(pathes, call(path))
		}
	}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil || (dir && !info.IsDir()) || (!dir && info.IsDir()) {
			return nil
		}
		if MatchRegex(path, reg_inc...) {
			add_path(path)
			return nil
		} else if MatchRegex(path, reg_exc...) {
			return nil
		} else {
			add_path(path)
			return nil
		}
	})
	return pathes
}

//check if file path exist and whether dir or not
func FileExists(path string) (dir, exists bool) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), true
	} else {
		return false, false
	}
}

//touch file by mode
func TouchFileMode(path string, mode os.FileMode) error {
	file, err := os.Open(path)
	if err != nil {
		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, mode)
		if err != nil {
			return err
		}
		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	info, err := file.Stat()
	if err == nil && info.IsDir() {
		return errors.New("can't touch path")
	} else {
		return err
	}
}

//touch file by default mode
func TouchFile(path string) error {
	return TouchFileMode(path, os.ModePerm)
}

//write []byte data to file
func WriteFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

//write string to file
func WriteFileString(path string, data string) error {
	return WriteFile(path, []byte(data))
}

//write a reader to file
func WriteFileReader(path string, buf io.Reader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, buf)
	return err
}

//append data to file
func AppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

//append string to file
func AppendFileString(path, data string) error {
	return AppendFile(path, []byte(data))
}

//append reader to file
func AppendFileReader(path string, buf io.Reader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, buf)
	return err
}

//copy file
func CopyFile(src string, dst string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()
	df, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer df.Close()
	return io.Copy(df, sf)
}

//check all file whether file exist or not order by list.
//it will read the first found file and return data.
//if all file not exists, return os.ErrNotExist
func CheckReadFile(paths ...string) ([]byte, error) {
	for _, path := range paths {
		if _, ok := FileExists(path); ok {
			return ioutil.ReadFile(path)
		}
	}
	return nil, os.ErrNotExist
}

func FileProtocolPath(path string) string {
	path = strings.Trim(path, " \t")
	if strings.HasPrefix(path, "file://") {
		return path
	}
	path, _ = filepath.Abs(path)
	path = strings.Replace(path, "\\", "/", -1)
	return "file://" + path
}
