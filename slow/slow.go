package slow

import (
	"io"
	"time"
)

type SlowReader struct {
	delay time.Duration
	r     io.Reader
}

func (sr SlowReader) Read(p []byte) (int, error) {
	time.Sleep(sr.delay)
	return sr.r.Read(p[:1])
}

func NewReader(r io.Reader, bps int) io.Reader {
	delay := time.Second / time.Duration(bps)
	return SlowReader{
		r:     r,
		delay: delay,
	}
}

/*
func example() {
	s := strings.NewReader("Not very long line...")
	r := NewReader(s, 4) //4byte par sec
	io.Copy(os.Stdout, r)
}
*/
