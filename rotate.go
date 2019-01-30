package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

func Rotate(filename string, copies int) error {
	var merr error
	for i := copies; i >= 2; i-- {
		err := os.Rename(
			filename+"."+strconv.Itoa(i-1),
			filename+"."+strconv.Itoa(i),
		)
		if !os.IsNotExist(err) {
			merr = multierror.Append(merr, err)
		}
	}
	f1, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}
	defer f1.Close()

	f2, err := os.Create(filename + ".1")
	if err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}
	defer f2.Close()

	_, err = io.Copy(f2, f1)
	if err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}

	err = f1.Truncate(0)
	if err != nil {
		merr = multierror.Append(merr, err)
		return merr
	}

	return merr
}

func RotateDir(dirname string, copies int) error {
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	var merr error
	for _, info := range infos {
		if info.IsDir() {
			err := RotateDir(path.Join(dirname, info.Name()), copies)
			if err != nil {
				merr = multierror.Append(merr, err)
			}
		} else if regexp.MustCompile("\\.log$").MatchString(info.Name()) {
			err := Rotate(path.Join(dirname, info.Name()), copies)
			if err != nil {
				merr = multierror.Append(merr, err)
			}
		} else if regexp.MustCompile("\\.log\\.(\\d+)$").MatchString(info.Name()) {
			
		}
	}

	return merr
}

func main() {
	root := "/var/lib/docker/containers"
	copies := 5
	for {
		err := RotateDir(root, copies)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Hour)
	}
}
