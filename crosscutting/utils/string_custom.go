package utils

type Email string

func (e Email) String() string {
	return string(e)
}

type IPAddress string

func (ip IPAddress) String() string {
	return string(ip)
}
