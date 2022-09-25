package webtoy_base

import "golang.org/x/crypto/bcrypt"

func BcryptHashPasswd(passwd string) (string, error) {
	var passwdBytes = []byte(passwd)
	hashBytes, err := bcrypt.GenerateFromPassword(passwdBytes, bcrypt.DefaultCost)
	return string(hashBytes), err
}

func BCryptMatchPasswd(hashPasswd, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPasswd), []byte(passwd))
	return err == nil
}
