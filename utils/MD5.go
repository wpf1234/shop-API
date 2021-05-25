package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func StrMd5(str string ) string{
	w := md5.New()
	_,err := io.WriteString(w, str)
	if err!=nil{
		fmt.Errorf("%s",err)
		return err.Error()
	}
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}
