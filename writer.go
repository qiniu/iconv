package iconv

import (
	"io"
	"syscall"
)

type Writer struct {
	inbuf []byte
	outbuf []byte
	cd Iconv
	output io.Writer
	n int // inbuf[0:n] is valid
	autoSync bool
}

func NewWriter(cd Iconv, output io.Writer, bufSize int, autoSync bool) *Writer {
	if bufSize < 16 { bufSize = DefaultBufSize }
	outbuf := make([]byte, bufSize)
	var inbuf []byte
	if !autoSync {
		inbuf = make([]byte, bufSize)
	}
	return &Writer{inbuf, outbuf, cd, output, 0, autoSync}
}

func (w *Writer) Output(w1 io.Writer) {
	w.Sync()
	w.output = w1
	w.n = 0
}

func (w *Writer) AutoSync(b bool) {
	w.autoSync = b
	if !b && w.inbuf == nil {
		w.inbuf = make([]byte, len(w.outbuf))
	}
}

func (w *Writer) Sync() error {

	if w.n == 0 { return nil }
	
	inleft, err := w.cd.DoWrite(w.output, w.inbuf, w.n, w.outbuf)
	if inleft > 0 {
		copy(w.inbuf, w.inbuf[w.n-inleft:w.n])
	}
	w.n = inleft
	return err
}

func (w *Writer) Write(b []byte) (n int, err error) {

	if w.autoSync {
		var inleft int
		inleft, err = w.cd.DoWrite(w.output, b, len(b), w.outbuf)
		n = len(b) - inleft
		return
	}
	for {
		n1 := copy(w.inbuf[w.n:], b)
		if n1 == 0 {
			if len(b) > 0 { return n, EILSEQ }
			break
		}
		w.n += n1
		n += n1
		if w.n == len(w.inbuf) {
			err = w.Sync()
			if err != nil && err != syscall.EINVAL { return }
		}
		if len(b) == n1 { break }
		b = b[n1:]
	}
	return n, nil
}

func (w *Writer) WriteString(b string) (n int, err error) {

	if w.autoSync {
		var inleft int
		inleft, err = w.cd.DoWrite(w.output, []byte(b), len(b), w.outbuf)
		n = len(b) - inleft
		return
	}
	for {
		n1 := copy(w.inbuf[w.n:], b)
		if n1 == 0 {
			if len(b) > 0 { return n, EILSEQ }
			break
		}
		w.n += n1
		n += n1
		if w.n == len(w.inbuf) {
			err = w.Sync()
			if err != nil && err != syscall.EINVAL { return }
		}
		if len(b) == n1 { break }
		b = b[n1:]
	}
	return n, nil
}

