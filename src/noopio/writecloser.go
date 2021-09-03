package noopio

type WriteCloser struct{}

func (r *WriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (r *WriteCloser) Close() error {
	return nil
}
