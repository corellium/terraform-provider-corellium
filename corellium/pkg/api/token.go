package api

type Token struct {
	value string
}

var token *Token

func GetAccessToken() string {
	return token.value
}

func SetAccessToken(value string) {
	if token == nil {
		token = &Token{
			value,
		}
	}
}
