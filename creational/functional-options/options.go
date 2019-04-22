package functional_options

type Client interface {
	Do()
}
type client struct {
	address     string
	timeout     int
	retries     int
	isCheatMode bool
}

func (c *client) Do() {}

type Option func(*client)

func WithTimeout(t int) Option {
	return func(c *client) {
		c.timeout = t
	}
}
func WithRetries(r int) Option {
	return func(c *client) {
		c.retries = r
	}
}
func SetCheatMode() Option {
	return func(c *client) {
		c.isCheatMode = true
	}
}

func NewClient(addr string, opts ...Option) Client {
	c := client{address: addr}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}
