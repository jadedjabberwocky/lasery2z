package imagemap

type Options struct {
	imageWidth  float64
	imageHeight float64
}

func DefaultOptions() *Options {
	return &Options{}
}

func (o *Options) WithImageWidth(v float64) *Options {
	o.imageWidth = v
	return o
}

func (o *Options) WithImageHeight(v float64) *Options {
	o.imageHeight = v
	return o
}

func (o *Options) Check() (*Options, error) {
	return o, nil
}
