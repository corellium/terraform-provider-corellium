// Package api provides a simple way to get and set the access token.
package api

// Token is the access token to the Corellium API.
type Token struct {
	value string
}

// token is the instance of token that is used to store the access token.
var token *Token

// GetAccessToken returns the access token.
func GetAccessToken() string {
	return token.value
}

// SetAccessToken sets the access token.
func SetAccessToken(value string) {
	if token == nil {
		token = &Token{
			value,
		}
	}
}
