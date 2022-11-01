package api

import (
	"log"
	"net/http"
	_ "noteyouself/docs"
	"noteyouself/services"

	"github.com/gin-gonic/gin"
)

// @Summary 通过code获取用户唯一token
// @Tags 微信小程序交互
// @title 随心记接口文档
// @version 1.0
// @description 这是随心记微信小程序的后台接口文档
// @host 127.0.0.1:8001
// @Produce json
// @Param code body string true "小程序中前台wx.login获取的code"
// @Success 200 {string} json "{"code":200,"data":"","msg":"ok"}"
// @Router /api/v1/wxminiapp/getWxOpenidByCode [post]
func GetWxOpenidByCode(c *gin.Context) {

	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
    if(err != nil){
        log.Printf("参数解析失败了！")
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":  500,
            "msg":   "系统异常",
        })
    }
	log.Printf("%v", &json)
	code := json["code"].(string)
	errcode, errmsg, token := services.RequestWxMiniAppToken(code)
	c.JSON(http.StatusOK, gin.H{
		"code":  errcode,
		"msg":   errmsg,
		"token": token,
	})
}

// @Summary 解密小程序前台获取的加密信息
// @Tags 微信小程序交互
// @title 随心记接口文档
// @version 1.0
// @description 这是随心记微信小程序的后台接口文档
// @host 127.0.0.1:8001
// @Produce json
// @Param token body string true "小程序前后台交互唯一凭证"
// @Param encryptData body string true "加密数据"
// @Success 200 {string} json "{"code":200,"data":"","msg":"ok"}"
// @Router /api/v1/wxminiapp/decryptWxData [post]
func DecryptWxData(c *gin.Context) {
	json := make(map[string]interface{}) //注意该结构接受的内容
	err := c.BindJSON(&json)
    if(err != nil){
        log.Printf("参数解析失败了！")
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":  500,
            "msg":   "系统异常",
            })
    }
	log.Printf("%v", &json)
	token := json["token"].(string)
	encryptData := json["encryptData"].(string)
	iv := json["iv"].(string)
	signature := json["signature"].(string)
	errcode, errmsg, data := services.WxDecryptData(token, encryptData, iv, signature)
	c.JSON(http.StatusOK, gin.H{
		"code": errcode,
		"msg":  errmsg,
		"data": data,
	})
}
