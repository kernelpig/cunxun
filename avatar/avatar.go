package avatar

import (
	"io/ioutil"

	e "wangqingang/cunxun/error"
)

var AvatarBytes []byte

// TODO: 头像数据初始化暂时放这
func InitAvatar(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(e.SP(e.MConfigErr, e.ConfigLoadAvatarErr, err))
	}
	AvatarBytes = bytes
}
