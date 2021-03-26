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
	"strconv"
	"strings"
	"time"
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

// GetUserIDByUnionID 根据UnionID获取用户Userid
func (d *DingTalkAgent) GetUserIDByUnionID(unionid string) (string, error) {
	urlEndpoint, err := url.Parse("https://oapi.dingtalk.com/topapi/user/getbyunionid")
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Set("access_token", d.AccessToken)
	urlEndpoint.RawQuery = query.Encode()
	urlPath := urlEndpoint.String()

	resp, err := http.PostForm(urlPath, url.Values{"unionid": {unionid}})
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

	result := rdata["result"].(map[string]interface{})
	if result["contact_type"].(float64) != 0 {
		return "", errors.New("该用户不属于企业内部员工，无法登录。")
	}
	userid := result["userid"].(string)
	return userid, nil
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

// DingtalkQRLogin 用于钉钉扫码登录
type DingtalkQRLogin struct {
	AppSecret string
	AppKey    string
}

// NewDingtalkQRLogin 构造钉钉扫码登录实例
func NewDingtalkQRLogin(appSecret, appKey string) DingtalkQRLogin {
	return DingtalkQRLogin{
		AppSecret: appSecret,
		AppKey:    appKey,
	}
}

// GetUnionIDByCode 获取扫码用户UnionID
func (d *DingtalkQRLogin) GetUnionIDByCode(code string) (userid string, err error) {
	var resp *http.Response
	//服务端通过临时授权码获取授权用户的个人信息
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1000000, 10) // 毫秒时间戳
	signature := d.encodeSHA256(timestamp)                            // 加密签名
	urlPath := fmt.Sprintf(
		"https://oapi.dingtalk.com/sns/getuserinfo_bycode?accessKey=%s&timestamp=%s&signature=%s",
		d.AppKey, timestamp, signature)

	// 构造请求数据
	param := struct {
		Tmp_auth_code string `json:"tmp_auth_code"`
	}{code}
	paraByte, _ := json.Marshal(param)
	paraString := string(paraByte)

	resp, err = http.Post(urlPath, "application/json;charset=UTF-8", strings.NewReader(paraString))
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
	unionid := rdata["user_info"].(map[string]interface{})["unionid"].(string)
	return unionid, nil
}

func (d *DingtalkQRLogin) encodeSHA256(timestamp string) string {
	// 钉钉签名算法实现
	h := hmac.New(sha256.New, []byte(d.AppSecret))
	h.Write([]byte(timestamp))
	sum := h.Sum(nil) // 二进制流
	tmpMsg := base64.StdEncoding.EncodeToString(sum)

	uv := url.Values{}
	uv.Add("0", tmpMsg)
	message := uv.Encode()[2:]

	return message
}
