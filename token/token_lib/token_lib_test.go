package token_lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
)

const (
	testTicketUserId      = 123
	testTicketUserRole    = 2
	testTicketLoginSource = common.WebSource
	testTicketTTL         = 0Xffff
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
	UserId:      testTicketUserId,
	IssueTime:   uint32(time.Now().Unix()),
	Role:        testTicketUserRole,
	TTL:         testTicketTTL,
	LoginSource: testTicketLoginSource,
}

func TestEncrypt(t *testing.T) {
	assert := assert.New(t)

	initPrivateKey([]byte(testPrivateKey))
	initPublicKeys([]byte(testPublicKey))

	token, err := Encrypt(1, &accessPayload)
	assert.Nil(err)
	assert.NotEmpty(token)
	t.Log("access token: ", token)

	tokenPayload, err := Decrypt(token)
	assert.Nil(err)
	assert.Equal(accessPayload, *tokenPayload)
}
