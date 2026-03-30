package constants

const (
	MinPasswordLength = 8
	MaxPasswordLength = 16
	MinNameLength = 2
	MaxNameLength = 64
	MinAge = 18
	MaxAge = 120
)

var (
	PhoneNumberPattern = "^1[3-9][0-9]{9}$" 
	EmailPattern = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
	DIDPattern = "^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$"
)
