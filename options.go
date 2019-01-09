package sqlxt

var (
	DefaultDriver    = "mysql"
	DefaultURI       = "root:123456@tcp(127.0.0.1:3306)/test"
	DefaultCharset   = "utf8mb4"
	DefaultParseTime = true
	DefaultMaxClient = 10
)

type Options struct {
	Driver    string
	URI       string
	Charset   string
	ParseTime bool
	MaxClient int
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	opt := Options{
		Driver:    DefaultDriver,
		URI:       DefaultURI,
		Charset:   DefaultCharset,
		ParseTime: DefaultParseTime,
		MaxClient: DefaultMaxClient,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Driver(v string) Option {
	return func(o *Options) {
		o.Driver = v
	}
}

func URI(v string) Option {
	return func(o *Options) {
		o.URI = v
	}
}

func Charset(v string) Option {
	return func(o *Options) {
		o.Charset = v
	}
}

func ParseTime(v bool) Option {
	return func(o *Options) {
		o.ParseTime = v
	}
}

func MaxClient(v int) Option {
	return func(o *Options) {
		o.MaxClient = v
	}
}
