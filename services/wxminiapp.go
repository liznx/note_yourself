package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"noteyouself/db"
	"noteyouself/util"
	"strings"
	"time"
)

type WxtokenResp struct {
	Openid     string `json:"openid"`
	Sessionkey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

const wxminiapp_url_prefix string = "https://api.weixin.qq.com"
const wxminiapp_appid string = "wxc5788af60e6d3cda"
const wxminiapp_secret string = "b1ddc1da8445297754d1cb9b4e6698bb"

func RequestWxMiniAppToken(code string) (int, string, string) {
	//调用微信鉴权接口 param:appid、secret、js_code、grant_type
	path := wxminiapp_url_prefix + "/sns/jscode2session?appid=" + wxminiapp_appid + "&secret=" + wxminiapp_secret +
		"&js_code=" + code + "&grant_type=authorization_code"

	rs, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer rs.Body.Close()
	resultbody, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resp WxtokenResp
	json.Unmarshal(resultbody, &resp)
	if resp.Errcode != 0 {
		return resp.Errcode, resp.Errmsg, ""
	}
	//生成token
	token := "noteyouself_" + util.Uuid()
	db.RedisClient.Set(token, resp.Openid+"|"+resp.Sessionkey, 30*time.Minute)
	return resp.Errcode, resp.Errmsg, token
}

func WxDecryptData(token string, encryptData string, iv string, signature string) (int, string, map[string]interface{}) {
	value, err := db.RedisClient.Get(token).Result()
	if err != nil {
		return 999, "系统异常！", nil
	}
	session_key := strings.Split(value, "|")[1]
	decryptData, err := util.DecryptWXOpenData(wxminiapp_appid, session_key, encryptData, iv)
	if err != nil {
		return 999, "解密失败！", nil
	}
	return 0, "success", decryptData
}
