package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mindoc-org/mindoc/utils/auth2"
	"net/http"
	"net/url"
	"time"
)

const (
	AppName = "dingtalk"

	callbackState = "mindoc"
)

type BasicResponse struct {
	Message string `json:"errmsg"`
	Code    int    `json:"errcode"`
}

func (r *BasicResponse) Error() string {
	return fmt.Sprintf("errcode=%d, errmsg=%s", r.Code, r.Message)
}

func (r *BasicResponse) AsError() error {
	if r == nil {
		return nil
	}

	if r.Code != 0 || r.Message != "ok" {
		return r
	}
	return nil
}

type AccessToken struct {
	// 文档: https://open.dingtalk.com/document/orgapp/obtain-orgapp-token
	*BasicResponse

	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`

	createTime time.Time
}

func (a AccessToken) GetToken() string {
	return a.AccessToken
}

func (a AccessToken) GetExpireIn() time.Duration {
	return time.Duration(a.ExpireIn) * time.Second
}

func (a AccessToken) GetExpireTime() time.Time {
	return a.createTime.Add(a.GetExpireIn())
}

type UserAccessToken struct {
	// 文档: https://open.dingtalk.com/document/orgapp/obtain-user-token
	*BasicResponse // 此接口未返回错误代码信息，仅仅能检查HTTP状态码

	ExpireIn     int    `json:"expireIn"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	CorpId       string `json:"corpId"`
}

type UserInfo struct {
	// 文档: https://open.dingtalk.com/document/orgapp/dingtalk-retrieve-user-information
	*BasicResponse

	NickName  string `json:"nick"`
	Avatar    string `json:"avatarUrl"`
	Mobile    string `json:"mobile"`
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	Email     string `json:"email"`
	StateCode string `json:"stateCode"`
}

type UserIdByUnion struct {
	// 文档: https://open.dingtalk.com/document/isvapp/query-a-user-by-the-union-id
	*BasicResponse

	RequestId string `json:"request_id"`
	Result    struct {
		ContactType int    `json:"contact_type"`
		UserId      string `json:"userid"`
	} `json:"result"`
}

func NewClient(appSecret string, appKey string) auth2.Client {
	return NewDingtalkClient(appSecret, appKey)
}

func NewDingtalkClient(appSecret string, appKey string) *DingtalkClient {
	return &DingtalkClient{AppSecret: appSecret, AppKey: appKey}
}

type DingtalkClient struct {
	AppSecret string
	AppKey    string

	token auth2.IAccessToken
}

func (d *DingtalkClient) GetAccessToken(ctx context.Context) (auth2.IAccessToken, error) {
	if d.token != nil {
		return d.token, nil
	}

	endpoint := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", d.AppKey, d.AppSecret)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	var token AccessToken
	if err := auth2.Request(req, &token); err != nil {
		return nil, err
	}

	token.createTime = time.Now()
	return token, nil
}

func (d *DingtalkClient) SetAccessToken(token auth2.IAccessToken) {
	d.token = token
}

func (d *DingtalkClient) BuildURL(callback string, _ bool) string {
	v := url.Values{}
	v.Set("redirect_uri", callback)
	v.Set("response_type", "code")
	v.Set("client_id", d.AppKey)
	v.Set("scope", "openid")
	v.Set("state", callbackState)
	v.Set("prompt", "consent")
	return "https://login.dingtalk.com/oauth2/auth?" + v.Encode()
}

func (d *DingtalkClient) ValidateCallback(state string) error {
	if state != callbackState {
		return errors.New("auth2.state.wrong")
	}
	return nil
}

func (d *DingtalkClient) getUserAccessToken(ctx context.Context, code string) (UserAccessToken, error) {
	val := map[string]string{
		"clientId":     d.AppKey,
		"clientSecret": d.AppSecret,
		"code":         code,
		"grantType":    "authorization_code",
	}

	jv, _ := json.Marshal(val)

	endpoint := "https://api.dingtalk.com/v1.0/oauth2/userAccessToken"
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jv))
	req.Header.Set("Content-Type", "application/json")

	var token UserAccessToken
	if err := auth2.Request(req, &token); err != nil {
		return token, err
	}

	return token, nil
}

func (d *DingtalkClient) getUserInfo(ctx context.Context, userToken UserAccessToken, unionId string) (UserInfo, error) {
	var user UserInfo

	endpoint := fmt.Sprintf("https://api.dingtalk.com/v1.0/contact/users/%s", unionId)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	req.Header.Set("x-acs-dingtalk-access-token", userToken.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	if err := auth2.Request(req, &user); err != nil {
		return user, err
	}
	return user, nil
}

func (d *DingtalkClient) getUserIdByUnion(ctx context.Context, union string) (UserIdByUnion, error) {
	var userId UserIdByUnion
	token, err := d.GetAccessToken(ctx)
	if err != nil {
		return userId, err
	}
	endpoint := fmt.Sprintf("https://oapi.dingtalk.com/topapi/user/getbyunionid?access_token=%s", token.GetToken())
	b, _ := json.Marshal(map[string]string{
		"unionid": union,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	if err := auth2.Request(req, &userId); err != nil {
		return userId, err
	}

	return userId, nil
}

func (d *DingtalkClient) GetUserInfo(ctx context.Context, code string) (auth2.UserInfo, error) {
	var info auth2.UserInfo
	userToken, err := d.getUserAccessToken(ctx, code)
	if err != nil {
		return info, err
	}

	userInfo, err := d.getUserInfo(ctx, userToken, "me")
	if err != nil {
		return info, err
	}

	userId, err := d.getUserIdByUnion(ctx, userInfo.UnionId)
	if err != nil {
		return info, err
	}

	if userId.Result.ContactType > 0 {
		return info, errors.New("auth2.user.outer")
	}

	info.UserId = userId.Result.UserId
	info.Mail = userInfo.Email
	info.Mobile = userInfo.Mobile
	info.Name = userInfo.NickName
	info.Avatar = userInfo.Avatar
	return info, nil
}
