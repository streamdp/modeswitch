package device

type Sender interface {
	Send(cmd ...string) ([]byte, error)
	Close() error
}
