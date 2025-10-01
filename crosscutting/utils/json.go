package utils

type JsonString string

func (v JsonString) String() string {
	return string(v)
}
