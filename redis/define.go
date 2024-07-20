package redis

type UnlockFunc func() error

type Option struct {
	Addr     string
	Password string
	DB       int
	timeout  int64
}
