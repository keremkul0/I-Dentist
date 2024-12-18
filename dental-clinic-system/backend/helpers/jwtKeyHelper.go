package helpers

var jWTKey []byte

func SetJWTKey(jwtKey string) {
	jWTKey = []byte(jwtKey)
	return
}

func GetJWTKey() []byte {
	return jWTKey
}
