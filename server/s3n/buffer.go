package s3n

import (
	"bytes"
	"fmt"
	"strings"
	//"html/template"
)

type ResponseBuffer struct {
	b string
}

func (buffer *ResponseBuffer) ExecuteTemplate(name string, data interface{}) (err error) {
	var s bytes.Buffer
	err = T().ExecuteTemplate(&s, name, data)
	if err == nil {
		buffer.Print(s.String())
	}
	return
}

func (buffer *ResponseBuffer) Print(format string) {
	//buffer.b = fmt.Sprintf(buffer.b + format)
	buffer.b = buffer.b + format
}

func (buffer *ResponseBuffer) Println(format string) {
	buffer.Print(format + "\n")
}

func (buffer *ResponseBuffer) Printf(format string, a ...interface{}) {
	if len(a) == 0 {
		buffer.b = buffer.b + format
	} else {
		tmp := strings.Split(format, "%")
		for index, element := range tmp {
			if index == 0 {
				buffer.b += element
			} else {
				if index <= len(a) {
					buffer.b += fmt.Sprintf("%"+element, a[index-1])
				} else {
					buffer.b += "(!!MISING){s3n}[" + element + "]"
				}
			}
		}
	}
}

func (buffer *ResponseBuffer) Replace(what string, write string) {
	r := strings.NewReplacer(what, write)
	buffer.b = r.Replace(buffer.b)
}

func (buffer *ResponseBuffer) Addbefore(what string, write string) {
	buffer.Replace(what, write+what)
}

func (buffer *ResponseBuffer) Addafter(what string, write string) {
	buffer.Replace(what, what+write)
}

func (buffer *ResponseBuffer) Clear() {
	buffer.Empty()
}

func (buffer *ResponseBuffer) Empty() {
	buffer.b = ""
}

func (buffer *ResponseBuffer) Import(str string) {
	buffer.b = str
}

func (buffer *ResponseBuffer) Export() string {
	return buffer.b
}

func (buffer *ResponseBuffer) Debug() {
	fmt.Print(buffer.b + "\n============\n")
}
