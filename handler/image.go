package handler

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/oss"
)

// 处理froala_edit图片上传请求, 请求报文格式如下:
//------WebKitFormBoundaryAVJqmxYivAFBzgeg
//Content-Disposition: form-data; name="xToken"
//AQAAAK6llAfhtyHgF6UQyqSlxhgTB-RLGYp_nKGU7YZ7cSrGRC5EkIKI4doWVjtGeP1m9bqBWBUAMG0pJPoJXEa2TcOVecxZLAEDAAMAAAB3ZWI=
//------WebKitFormBoundaryAVJqmxYivAFBzgeg
//Content-Disposition: form-data; name="image_key"; filename="avatar.jpg"
//Content-Type: image/jpeg
//------WebKitFormBoundaryAVJqmxYivAFBzgeg--
func ImageCreateHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IImageCreate, e.MParamsErr, e.ParamsInvalidMultiForm, err))
		return
	}
	xToken, ok := form.Value["xToken"]
	if !ok || xToken == nil || len(xToken) < 1 {
		c.JSON(http.StatusBadRequest, e.I(e.IImageCreate, e.MTokenErr, e.TokenIsEmpty))
		return
	}
	if _, err := middleware.CheckAccessToken(xToken[0]); err != nil {
		c.JSON(http.StatusBadRequest, e.I(e.IImageCreate, e.MTokenErr, e.TokenSignVerifyErr))
		return
	}
	files, ok := form.File["image_key"]
	if !ok || files == nil || len(files) < 1 {
		c.JSON(http.StatusBadRequest, e.I(e.IImageCreate, e.MImageErr, e.ImageNotFound))
		return
	}

	// 暂时只处理一个文件上传
	file := files[0]
	newName := uuid.NewV4().String() + path.Ext(file.Filename)
	fd, err := file.Open()
	if err != nil || fd == nil {
		c.JSON(http.StatusInternalServerError, e.I(e.IImageCreate, e.MImageErr, e.ImageReadErr))
		return
	}
	defer fd.Close()

	link, err := oss.PutImageByFile(newName, fd)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IImageCreate, e.MImageErr, e.ImageSaveErr, err))
		return
	}

	// 必须为此格式, 为forala_edit处理特定格式
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"link": link,
	})
}
