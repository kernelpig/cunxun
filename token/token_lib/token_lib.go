package token_lib

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"sync"
)

type counter struct {
	sync.Mutex
	num uint32
}

var count counter

type Payload struct {
	IssueTime   uint32 // 时间戳，秒
	TTL         uint16 // 单位为分钟
	UserId      uint32 // 4B
	LoginSource string // 3B
}

var privateKey *ecdsa.PrivateKey
var publicKeys []*ecdsa.PublicKey

func initPrivateKey(key []byte) {
	block, _ := pem.Decode(key)
	if block == nil {
		panic("private key invalid")
	}
	prvk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	privateKey = prvk
}

func initPublicKeys(keys ...[]byte) {
	if len(keys) == 0 {
		panic("no public key")
	}
	publicKeys = make([]*ecdsa.PublicKey, 0)
	for _, key := range keys {
		block, _ := pem.Decode(key)
		if block == nil {
			panic("public key invalid")
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			panic(err)
		}
		pubKey := pubInterface.(*ecdsa.PublicKey)
		publicKeys = append(publicKeys, pubKey)
	}
}

func InitKeyPem(publicKeyFile, privateKeyFile string) {
	publickKeyPem, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		panic(err)
	}

	privateKeyPem, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}

	initPrivateKey(privateKeyPem)
	initPublicKeys(publickKeyPem)
}

func Encrypt(version int, tk *Payload) (string, error) {

	count.Lock()
	seq := count.num
	count.num += 1
	count.Unlock()

	switch version {
	case 1:
		data, err := tk.encryptV1(seq)
		if err != nil {
			return "", err
		}
		token := base64.URLEncoding.EncodeToString(data)
		return token, nil
	default:
		return "", errors.New("invalid version")
	}
}

func Decrypt(token string) (*Payload, error) {
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, errors.New("invalid base64 token")
	}
	switch int(data[0]) {
	case 1:
		token := &Payload{}
		err := token.decryptV1(data)
		if err != nil {
			return nil, err
		}
		return token, nil
	default:
		return nil, errors.New("invalid version")
	}
}

func (t *Payload) encryptV1(seq uint32) ([]byte, error) {

	var datas = make([]byte, 10)
	var sigs = make([]byte, 1)
	var seqs = make([]byte, 4)

	// SIG = VERSION(1B) + SEQ(3B)
	sigs[0] = 0x01
	binary.LittleEndian.PutUint32(seqs[0:4], seq)
	sigs = append(sigs, seqs[0:3]...)

	// PAYLOAD = ISSUETIME(4B) + TTL(2B) + ACCOUNTID(32B) + SOURCE(>0B)
	binary.LittleEndian.PutUint32(datas[0:4], t.IssueTime)
	binary.LittleEndian.PutUint16(datas[4:6], uint16(t.TTL))
	binary.LittleEndian.PutUint32(datas[6:10], t.UserId)
	datas = append(datas, []byte(t.LoginSource)...)

	// 对SIG+PAYLOAD做签名，SIGN = SIGN_R(32B)+SIGN_S(32B)
	h := md5.New()
	h.Write(append(sigs, datas...))
	hashed := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return nil, err
	}

	// TOKEN = SIG{VERSION(1B)+SEQ(3B)} + SIGN{SIGN_R(32B)+SIGN_S(32B)} + PAYLOAD{ISSUETIME(4B) + TTL(2B) + ACCOUNTID(32B) + SOURCE(>0B)}
	sigs = append(sigs, r.Bytes()...)
	sigs = append(sigs, s.Bytes()...)
	sigs = append(sigs, datas...)

	return sigs, nil
}

func (t *Payload) decryptV1(data []byte) error {

	// SIG = VERSION(1B) + SEQ(3B)
	sigs := data[0:4]

	// PAYLOAD = ISSUETIME(4B) + TTL(2B) + ACCOUNTID(32B) + SOURCE(>0B)
	payload := data[68:]

	// SIGN = SIGN_R(32B)+SIGN_S(32B)
	r := big.NewInt(0)
	s := big.NewInt(0)
	r = r.SetBytes(data[4:36])
	s = s.SetBytes(data[36:68])

	// 签名校验
	h := md5.New()
	h.Write(append(sigs, payload...))
	hashed := h.Sum(nil)

	// 签名校验
	valid := false
	for _, pubk := range publicKeys {
		ok := ecdsa.Verify(pubk, hashed, r, s)
		if ok {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("sign verify failed")
	}

	// 解析载荷
	t.IssueTime = uint32(binary.LittleEndian.Uint32(data[68:72]))
	t.TTL = uint16(binary.LittleEndian.Uint16(data[72:74]))
	t.UserId = uint32(binary.LittleEndian.Uint32(data[74:78]))
	t.LoginSource = string(data[78:])

	return nil
}
