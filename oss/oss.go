package oss

import (
	"fmt"
	"io"
	"path"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"

	"wangqingang/cunxun/common"
	e "wangqingang/cunxun/error"
)

func PutImageByFile(fileName string, reader io.Reader) (string, error) {
	client, err := alioss.New(common.Config.Oss.Endpoint, common.Config.Oss.AliAccessId, common.Config.Oss.AliAccessSecret)
	if err != nil {
		return "", e.SP(e.MOssErr, e.OssClientInitErr, err)
	}
	bucket, err := client.Bucket(common.Config.Oss.Bucket)
	if err != nil {
		return "", e.SP(e.MOssErr, e.OssBucketGetErr, err)
	}
	fileExt := path.Ext(fileName)
	if fileExt == "" || fileExt == "." {
		fileExt = ".png" // 默认使用png类型
	}
	options := []alioss.Option{
		alioss.ContentType("image/" + fileExt[1:]),
	}
	err = bucket.PutObject(fileName, reader, options...)
	if err != nil {
		return "", e.SP(e.MOssErr, e.OssPutObjectByBytesErr, err)
	}
	imagePath := path.Join(common.Config.Oss.Domain, fileName)
	return fmt.Sprintf("http://%s", imagePath), nil
}
