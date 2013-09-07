package iconv

import (
	"io"
)

type Reader struct {
	rdbuf []byte
	cnvbuf []byte
	cd Iconv
	input io.Reader
	from, to int // rdbuf[from:to] is valid
	m int // cnvbuf[:m] is valid
	err error
}

func NewReader(cd Iconv, input io.Reader, bufSize int) *Reader {
	if bufSize < 16 { bufSize = DefaultBufSize }
	rdbuf := make([]byte, bufSize)
	cnvbuf := make([]byte, bufSize)
	return &Reader{rdbuf, cnvbuf, cd, input, 0, 0, 0, nil}
}

func (r *Reader) Input(r1 io.Reader) {
	r.input = r1
	r.from, r.to, r.m = 0, 0, 0
	r.err = nil
}

func (r *Reader) fetch() error {

	var m int

	if r.err != nil { return r.err }

	m, r.err = r.input.Read(r.cnvbuf[r.m:])
	m += r.m
	if m == 0 { return io.EOF }

	r.from = 0
	r.to, r.m, r.err = r.cd.Do(r.cnvbuf, m, r.rdbuf)
	if r.err != EILSEQ {
		r.err = nil
	}
	if r.m > 0 {
		copy(r.cnvbuf[:r.m], r.cnvbuf[m-r.m:m])
	}
	if r.to == 0 {
		if r.err == nil { return io.EOF }
		return r.err
	}
	return nil
}

func (r *Reader) Read(b []byte) (n int, err error) {

	for {
		if r.from < r.to {
			n1 := copy(b, r.rdbuf[r.from:r.to])
			n += n1
			r.from += n1
			if n1 == len(b) {
				break
			}
			b = b[n1:]
		}
		err = r.fetch()
		if err != nil { break }
	}
	return
}

