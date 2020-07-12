package api

const (
	DefaultLimit = 10
	DefaultHost  = "api.huobi.pro"
)

// CallOptions
const (
	DefaultSize = 20
)

// Options is an exchange option
type Options struct {
	TraderID  int64
	Host      string
	Type      string
	Name      string
	AccessKey string
	SecretKey string

	Limit  int64
	Source string
}

type Option func(opt *Options)

func newOption(opts ...Option) Options {
	opt := Options{
		Limit: DefaultLimit,
		Host:  DefaultHost,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func LimitOption(times int64) Option {
	return func(opts *Options) {
		opts.Limit = times
	}
}

func HostOption(host string) Option {
	return func(opts *Options) {
		opts.Host = host
	}
}

func SourceOption(source string) Option {
	return func(opts *Options) {
		opts.Source = source
	}
}

type CallOptions struct {
	Size int
}

type CallOption func(opts *CallOptions)

func Size(size int) CallOption {
	return func(opts *CallOptions) {
		opts.Size = size
	}
}

func newCallOption() CallOptions {
	opts := CallOptions{
		Size: DefaultSize,
	}
	return opts
}
