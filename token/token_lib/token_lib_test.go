package token_lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
)

const (
	TestTicketAccoutId    = "aid01234567890123456789012345678"
	TestTicketLoginSource = common.WebSource
	TestTicketTTL         = 0Xffff
)

const (
	testPrivateKey = `
-----BEGIN PRIVATE KEY-----
MHcCAQEEIBj5t3xVfkNcboOoMz7t9PFB3TGV0AAwKkkjFCgEtg7NoAoGCCqGSM49
AwEHoUQDQgAE/pqnVRDtspk0sSQu9ihL9plZgrxbgBU28nrAKjjwQUmSu1ZPyJw/
YrOxYrr92h6oOVNAK0pF6qAka6nryXZjQQ==
-----END PRIVATE KEY-----
`
	testPublicKey = `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/pqnVRDtspk0sSQu9ihL9plZgrxb
gBU28nrAKjjwQUmSu1ZPyJw/YrOxYrr92h6oOVNAK0pF6qAka6nryXZjQQ==
-----END PUBLIC KEY-----
`
)

var accessPayload = Payload{
	AccountId:   TestTicketAccoutId,
	IssueTime:   uint32(time.Now().Unix()),
	TTL:         TestTicketTTL,
	LoginSource: TestTicketLoginSource,
}

func TestEncrypt(t *testing.T) {
	assert := assert.New(t)

	InitPrivateKey([]byte(testPrivateKey))
	InitPublicKeys([]byte(testPublicKey))

	token, err := Encrypt(1, &accessPayload)
	assert.Nil(err)
	assert.NotEmpty(token)
	t.Log("access token: ", token)

	tokenPayload, err := Decrypt(token)
	assert.Nil(err)
	assert.Equal(accessPayload, *tokenPayload)
}
