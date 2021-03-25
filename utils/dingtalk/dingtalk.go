package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DingTalkAgent 用于钉钉交互
type DingTalkAgent struct {
	AppSecret   string
	AppKey      string
	AccessToken string
}

// NewDingTalkAgent 钉钉交互构造函数
func NewDingTalkAgent(appSecret, appKey string) *DingTalkAgent {
	return &DingTalkAgent{
		AppSecret: appSecret,
		AppKey:    appKey,
	}
}

// GetUserIDByCode 通过临时code获取当前用户ID
func (d *DingTalkAgent) GetUserIDByCode(code string) (string, error) {
	urlEndpoint, err := url.Parse("https://oapi.dingtalk.com/user/getuserinfo")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("access_token", d.AccessToken)
	query.Set("code", code)

	urlEndpoint.RawQuery = query.Encode()
	urlPath := urlEndpoint.String()

	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析钉钉返回数据
	var rdata map[string]interface{}
	err = json.Unmarshal(body, &rdata)
	if err != nil {
		return "", err
	}

	errcode := rdata["errcode"].(float64)
	if errcode != 0 {
		return "", errors.New(fmt.Sprintf("登录错误: %.0f, %s", errcode, rdata["errmsg"].(string)))
	}

	userid := rdata["userid"].(string)
	return userid, nil
}

// GetUserNameAndAvatarByUserID 通过userid获取当前用户姓名和头像
func (d *DingTalkAgent) GetUserNameAndAvatarByUserID(userid string) (string, string, error) {
	urlEndpoint, err := url.Parse("https://oapi.dingtalk.com/topapi/v2/user/get")
	if err != nil {
		return "", "", err
	}

	query := url.Values{}
	query.Set("access_token", d.AccessToken)

	urlEndpoint.RawQuery = query.Encode()
	urlPath := urlEndpoint.String()

	resp, err := http.PostForm(urlPath, url.Values{"userid": {userid}})
	if err != nil {
		return "", "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// 解析钉钉返回数据
	var rdata map[string]interface{}
	err = json.Unmarshal(body, &rdata)
	if err != nil {
		return "", "", err
	}

	errcode := rdata["errcode"].(float64)
	if errcode != 0 {
		return "", "", errors.New(fmt.Sprintf("登录错误: %.0f, %s", errcode, rdata["errmsg"].(string)))
	}

	userinfo := rdata["result"].(map[string]interface{})
	username := userinfo["name"].(string)
	avatar := userinfo["avatar"].(string)
	return username, avatar, nil
}

// GetAccesstoken 获取钉钉请求Token
func (d *DingTalkAgent) GetAccesstoken() (err error) {

	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", d.AppKey, d.AppSecret)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var i map[string]interface{}
	err = json.Unmarshal(body, &i)
	if err != nil {
		return err
	}

	if i["errcode"].(float64) == 0 {
		d.AccessToken = i["access_token"].(string)
		return nil
	}
	return errors.New("accesstoken获取错误：" + i["errmsg"].(string))
}

func (d *DingTalkAgent) encodeSHA256(message string) string {
	// 钉钉签名算法实现
	h := hmac.New(sha256.New, []byte(d.AppSecret))
	h.Write([]byte(message))
	sum := h.Sum(nil) // 二进制流
	tmpMsg := base64.StdEncoding.EncodeToString(sum)

	uv := url.Values{}
	uv.Add("0", tmpMsg)
	message = uv.Encode()[2:]

	return message
}
