package webtoy_base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptMatch(t *testing.T) {
	passwd := "$%^&*)nmf"

	hashPasswd, err := BcryptHashPasswd(passwd)
	assert.Nil(t, err, "bcrypt hash password")

	ret := BCryptMatchPasswd(hashPasswd, passwd)
	assert.True(t, ret, "bcrypt match password")

	errPasswd := "*23JN89F(S]*"
	ret = BCryptMatchPasswd(hashPasswd, errPasswd)
	assert.False(t, ret, "bcrypt don't match error password")
}
