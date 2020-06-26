// Package iconv is golang bindings to libiconv that converts string to
// requested character encoding.
package iconv

// #cgo darwin  LDFLAGS: -liconv
// #cgo freebsd LDFLAGS: -liconv
// #cgo windows LDFLAGS: -liconv
// #include <iconv.h>
// #include <stdlib.h>
// #include <errno.h>
//
// size_t bridge_iconv(iconv_t cd,
//		       char *inbuf, size_t *inbytesleft,
//                     char *outbuf, size_t *outbytesleft) {
//   return iconv(cd, &inbuf, inbytesleft, &outbuf, outbytesleft);
// }
import "C"

import (
	"bytes"
	"io"
	"syscall"
	"unsafe"
)

var (
	// EILSEQ error
	EILSEQ = syscall.Errno(C.EILSEQ)
	// E2BIG error
	E2BIG = syscall.Errno(C.E2BIG)
)

// DefaultBufSize const
const DefaultBufSize = 4096

// Iconv represents an iconv handle.
type Iconv struct {
	Handle C.iconv_t
}

// Open returns a conversion descriptor cd, cd contains a conversion state and can not be used in multiple threads simultaneously.
func Open(tocode string, fromcode string) (cd Iconv, err error) {
	tocode1 := C.CString(tocode)
	defer C.free(unsafe.Pointer(tocode1))

	fromcode1 := C.CString(fromcode)
	defer C.free(unsafe.Pointer(fromcode1))

	ret, err := C.iconv_open(tocode1, fromcode1)
	if err != nil {
		return
	}
	cd = Iconv{ret}
	return
}

// Close closes the iconv handle.
func (cd Iconv) Close() error {
	_, err := C.iconv_close(cd.Handle)
	return err
}

// Conv converts text to requested character encoding.
func (cd Iconv) Conv(b []byte, outbuf []byte) (out []byte, inleft int, err error) {
	outn, inleft, err := cd.Do(b, len(b), outbuf)
	if err == nil || err != E2BIG {
		out = outbuf[:outn]
		return
	}

	w := bytes.NewBuffer(nil)
	w.Write(outbuf[:outn])

	inleft, err = cd.DoWrite(w, b[len(b)-inleft:], inleft, outbuf)
	if err != nil {
		return
	}
	out = w.Bytes()
	return
}

// ConvString converts string to requested character encoding.
func (cd Iconv) ConvString(s string) string {
	var outbuf [512]byte
	s1, _, err := cd.Conv([]byte(s), outbuf[:])
	if err != nil {
		return ""
	}
	return string(s1)
}

// Do converts text to requested character encoding.
func (cd Iconv) Do(inbuf []byte, in int, outbuf []byte) (out, inleft int, err error) {
	if in == 0 {
		return
	}
	inbytes := C.size_t(in)
	inptr := &inbuf[0]

	outbytes := C.size_t(len(outbuf))
	outptr := &outbuf[0]
	_, err = C.bridge_iconv(cd.Handle,
		(*C.char)(unsafe.Pointer(inptr)), &inbytes,
		(*C.char)(unsafe.Pointer(outptr)), &outbytes)

	out = len(outbuf) - int(outbytes)
	inleft = int(inbytes)
	return
}

// DoWrite converts text to requested character encoding and writes into a Writer.
func (cd Iconv) DoWrite(w io.Writer, inbuf []byte, in int, outbuf []byte) (inleft int, err error) {
	if in == 0 {
		return
	}
	inbytes := C.size_t(in)
	for inbytes > 0 {
		in = int(inbytes)
		inptr := &inbuf[len(inbuf)-in]
		outbytes := C.size_t(len(outbuf))
		outptr := &outbuf[0]
		_, err = C.bridge_iconv(cd.Handle,
			(*C.char)(unsafe.Pointer(inptr)), &inbytes,
			(*C.char)(unsafe.Pointer(outptr)), &outbytes)
		w.Write(outbuf[:len(outbuf)-int(outbytes)])
		if err != nil && err != E2BIG {
			return int(inbytes), err
		}
	}
	return 0, nil
}
