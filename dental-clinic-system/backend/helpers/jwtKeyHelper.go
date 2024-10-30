package helpers

var jWTKey []byte

func NewJWTKeyHelper(jwtKey []byte) {
	jWTKey = jwtKey
	return
}

func GetJWTKey() []byte {
	return jWTKey
}
