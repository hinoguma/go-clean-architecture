package adapter

type StrIDGenerator interface {
	Issue() string
}
