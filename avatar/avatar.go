package avatar

import (
	"io/ioutil"
	"path"

	e "wangqingang/cunxun/error"
)

var AvatarBytes []byte

func InitAvatar(dir, file string) {
	bytes, err := ioutil.ReadFile(path.Join(dir, file))
	if err != nil {
		panic(e.SP(e.MConfigErr, e.ConfigLoadAvatarErr, err))
	}
	AvatarBytes = bytes
}

// 获取配置的默认头像, 有的话则更新bytes, 用于更换avatar免重启
func GetDefaultAvatar(dir, file string) []byte {
	bytes, err := ioutil.ReadFile(path.Join(dir, file))
	if err != nil {
		return AvatarBytes
	}
	AvatarBytes = bytes
	return bytes
}