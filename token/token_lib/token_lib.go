package token_lib

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"sync"

	e "wangqingang/cunxun/error"
)

type counter struct {
	sync.Mutex
	num uint32
}

var count counter

type Payload struct {
	IssueTime   uint32 // 4B, 时间戳，秒
	TTL         uint16 // 2B, 单位为分钟
	Role        uint16 // 2B, 用户角色
	UserId      uint64 // 8B
	LoginSource string // 3B
}

var privateKey *ecdsa.PrivateKey
var publicKeys []*ecdsa.PublicKey

func initPrivateKey(key []byte) error {
	block, _ := pem.Decode(key)
	if block == nil {
		return e.S(e.MTokenErr, e.TokenDecryptErr)
	}
	prvk, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return e.S(e.MTokenErr, e.TokenParsePriKey)
	}
	privateKey = prvk
	return nil
}

func initPublicKeys(keys ...[]byte) error {
	if len(keys) == 0 {
		return e.S(e.MTokenErr, e.TokenIsEmpty)
	}
	publicKeys = make([]*ecdsa.PublicKey, 0)
	for _, key := range keys {
		block, _ := pem.Decode(key)
		if block == nil {
			return e.S(e.MTokenErr, e.TokenInvalidPubKey)
		}
		pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return e.S(e.MTokenErr, e.TokenInvalidPubKey)
		}
		pubKey := pubInterface.(*ecdsa.PublicKey)
		publicKeys = append(publicKeys, pubKey)
	}
	return nil
}

func InitKeyPem(publicKeyFile, privateKeyFile string) error {
	publickKeyPem, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return e.S(e.MTokenErr, e.TokenReadPubKeyFileErr)
	}
	privateKeyPem, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return e.S(e.MTokenErr, e.TokenReadPriKeyFileErr)
	}
	if err := initPrivateKey(privateKeyPem); err != nil {
		return err
	}
	if err := initPublicKeys(publickKeyPem); err != nil {
		return err
	}
	return nil
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
			return "", e.SP(e.MTokenErr, e.TokenEcryptErr, err)
		}
		token := base64.URLEncoding.EncodeToString(data)
		return token, nil
	default:
		return "", e.SP(e.MTokenErr, e.TokenInvalidVersion, nil)
	}
}

func Decrypt(token string) (*Payload, error) {
	data, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, e.SP(e.MTokenErr, e.TokenBase64DecodeErr, nil)
	}
	switch int(data[0]) {
	case 1:
		token := &Payload{}
		err := token.decryptV1(data)
		if err != nil {
			return nil, e.SP(e.MTokenErr, e.TokenDecryptErr, nil)
		}
		return token, nil
	default:
		return nil, e.SP(e.MTokenErr, e.TokenInvalidVersion, nil)
	}
}

func (t *Payload) encryptV1(seq uint32) ([]byte, error) {
	// 固定字节bytes缓冲区, 16B = ISSUETIME(4B) + TTL(2B) + ROLE(2B) + USERID(8B)
	var datas = make([]byte, 16)
	var sigs = make([]byte, 1)
	var seqs = make([]byte, 4)

	// SIG = VERSION(1B) + SEQ(3B)
	sigs[0] = 0x01
	binary.LittleEndian.PutUint32(seqs[0:4], seq)
	sigs = append(sigs, seqs[0:3]...)

	// PAYLOAD = ISSUETIME(4B) + TTL(2B) + ROLE(2B) + USERID(4B) + SOURCE(>0B)
	binary.LittleEndian.PutUint32(datas[0:4], t.IssueTime)
	binary.LittleEndian.PutUint16(datas[4:6], uint16(t.TTL))
	binary.LittleEndian.PutUint16(datas[6:8], uint16(t.Role))
	binary.LittleEndian.PutUint64(datas[8:16], t.UserId)
	datas = append(datas, []byte(t.LoginSource)...)

	// 对SIG+PAYLOAD做签名，SIGN = SIGN_R(32B)+SIGN_S(32B)
	h := md5.New()
	h.Write(append(sigs, datas...))
	hashed := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed)
	if err != nil {
		return nil, e.SP(e.MTokenErr, e.TokenSignErr, err)
	}

	// TOKEN = SIG{VERSION(1B)+SEQ(3B)} + SIGN{SIGN_R(32B)+SIGN_S(32B)} + PAYLOAD{ISSUETIME(4B) + TTL(2B) + USERID(32B) + SOURCE(>0B)}
	sigs = append(sigs, r.Bytes()...)
	sigs = append(sigs, s.Bytes()...)
	sigs = append(sigs, datas...)

	return sigs, nil
}

func (t *Payload) decryptV1(data []byte) error {

	// SIG = VERSION(1B) + SEQ(3B)
	sigs := data[0:4]

	// PAYLOAD = ISSUETIME(4B) + TTL(2B) + ROLE(2B) + USERID(8B) + SOURCE(>0B)
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
		return e.SP(e.MTokenErr, e.TokenSignVerifyErr, nil)
	}

	// 解析载荷
	t.IssueTime = uint32(binary.LittleEndian.Uint32(data[68:72]))
	t.TTL = uint16(binary.LittleEndian.Uint16(data[72:74]))
	t.Role = uint16(binary.LittleEndian.Uint16(data[74:76]))
	t.UserId = uint64(binary.LittleEndian.Uint64(data[76:84]))
	t.LoginSource = string(data[84:])

	return nil
}
