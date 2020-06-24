package builder

var (
	// DefaultDriver 数据库引擎, 默认sqlite3
	DefaultDriver = "sqlite3"
	// DefaultURI 数据连接, 默认:memory:
	// mysql uri: [username[:password]@][protocol[(address)]]/dbname
	DefaultURI = ":memory:"
	// DefaultCharset 字符集，默认utf8mb4，仅支持mysql
	DefaultCharset = "utf8mb4"
	// DefaultParseTime 是否转义时间格式，仅支持mysql
	DefaultParseTime = true
	// DefaultMaxClient 最大客户端连接数，默认为10
	DefaultMaxClient = 10
)

// Options 连接池参数集
type Options struct {
	Driver    string
	URI       string
	Charset   string
	ParseTime bool
	MaxClient int
}

// Option 参数函数
type Option func(o *Options)

// NewOptions 返回一个新的参数集
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

// Driver 设置数据库引擎
func Driver(v string) Option {
	return func(o *Options) {
		o.Driver = v
	}
}

// URI 设置数据库连接地址
func URI(v string) Option {
	return func(o *Options) {
		o.URI = v
	}
}

// Charset 设置字符集
func Charset(v string) Option {
	return func(o *Options) {
		o.Charset = v
	}
}

// ParseTime 设置是否转义时间字段
func ParseTime(v bool) Option {
	return func(o *Options) {
		o.ParseTime = v
	}
}

// MaxClient 设置最大客户端连接数
func MaxClient(v int) Option {
	return func(o *Options) {
		o.MaxClient = v
	}
}
