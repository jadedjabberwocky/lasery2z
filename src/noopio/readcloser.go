package noopio

import "io"

type ReadCloser struct{}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (r *ReadCloser) Close() error {
	return nil
}
