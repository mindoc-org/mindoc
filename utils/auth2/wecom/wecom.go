package wecom

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

// doc
// - 全局错误码: https://work.weixin.qq.com/api/doc/90000/90139/90313

const (
	AppName = "workwx"

	auth2Url      = "https://open.weixin.qq.com/connect/oauth2/authorize"
	ssoUrl        = "https://login.work.weixin.qq.com/wwlogin/sso/login"
	callbackState = "mindoc"
)

type BasicResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (r *BasicResponse) Error() string {
	return fmt.Sprintf("errcode=%d,errmsg=%s", r.ErrCode, r.ErrMsg)
}

func (r *BasicResponse) AsError() error {
	if r.ErrCode != 0 {
		return r
	}
	return nil
}

// 获取用户Id-请求响应结构
type UserIdResponse struct {
	// 接口文档: https://developer.work.weixin.qq.com/document/path/91023
	*BasicResponse

	UserId         string `json:"userid"`          // 企业成员UserID
	UserTicket     string `json:"user_ticket"`     // 用于获取敏感信息
	OpenId         string `json:"openid"`          // 非企业成员的标识，对当前企业唯一
	ExternalUserId string `json:"external_userid"` // 外部联系人ID
}

// 获取用户信息-请求响应结构
type UserInfoResponse struct {
	// 接口文档: https://developer.work.weixin.qq.com/document/path/90196
	*BasicResponse

	UserId         string `json:"userid"`            // 企业成员UserID
	Name           string `json:"name"`              // 成员名称
	Department     []int  `json:"department"`        // 成员所属部门id列表
	IsLeaderInDept []int  `json:"is_leader_in_dept"` // 表示在所在的部门内是否为上级
	IsLeader       int    `json:"isleader"`          // 是否是部门上级(领导)
	Alias          string `json:"alias"`             // 别名
	Status         int    `json:"status"`            // 激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业
	MainDepartment int    `json:"main_department"`   // 主部门
}

type UserPrivateInfoResponse struct {
	// 文档地址: https://developer.work.weixin.qq.com/document/path/95833
	*BasicResponse

	UserId  string `json:"userid"`   // 企业成员userid
	Gender  string `json:"gender"`   // 成员性别
	Avatar  string `json:"avatar"`   // 头像
	QrCode  string `json:"qr_code"`  // 二维码
	Mobile  string `json:"mobile"`   // 手机号
	Mail    string `json:"mail"`     // 邮箱
	BizMail string `json:"biz_mail"` // 企业邮箱
	Address string `json:"address"`  // 地址
}

// 访问凭据缓存-结构
type AccessToken struct {
	*BasicResponse

	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`

	createTime time.Time `json:"create_time"`
}

func (a AccessToken) GetToken() string {
	return a.AccessToken
}

func (a AccessToken) GetExpireIn() time.Duration {
	return time.Duration(a.ExpiresIn) * time.Second
}

func (a AccessToken) GetExpireTime() time.Time {
	return a.createTime.Add(a.GetExpireIn())
}

// 企业微信用户敏感信息-结构
type WorkWeixinUserPrivateInfo struct {
	UserId  string `json:"userid"`   // 企业成员userid
	Name    string `json:"name"`     // 姓名
	Gender  string `json:"gender"`   // 成员性别
	Avatar  string `json:"avatar"`   // 头像
	QrCode  string `json:"qr_code"`  // 二维码
	Mobile  string `json:"mobile"`   // 手机号
	Mail    string `json:"mail"`     // 邮箱
	BizMail string `json:"biz_mail"` // 企业邮箱
	Address string `json:"address"`  // 地址
}

// 企业微信用户信息-结构
type WorkWeixinUserInfo struct {
	UserId         string `json:"UserId"`            // 企业成员UserID
	Name           string `json:"name"`              // 成员名称
	HideMobile     int    `json:"hide_mobile"`       // 是否隐藏了手机号码
	Mobile         string `json:"mobile"`            // 手机号码
	Department     []int  `json:"department"`        // 成员所属部门id列表
	Email          string `json:"email"`             // 邮箱
	IsLeaderInDept []int  `json:"is_leader_in_dept"` // 表示在所在的部门内是否为上级
	IsLeader       int    `json:"isleader"`          // 是否是部门上级(领导)
	Avatar         string `json:"avatar"`            // 头像url
	Alias          string `json:"alias"`             // 别名
	Status         int    `json:"status"`            // 激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业
	MainDepartment int    `json:"main_department"`   // 主部门
}

func NewClient(corpId, appId, appSecrete string) auth2.Client {
	return NewWorkWechatClient(corpId, appId, appSecrete)
}
func NewWorkWechatClient(corpId, appId, appSecrete string) *WorkWechatClient {
	return &WorkWechatClient{
		CorpId:    corpId,
		AppId:     appId,
		AppSecret: appSecrete,
	}
}

type WorkWechatClient struct {
	CorpId    string
	AppId     string
	AppSecret string

	token auth2.IAccessToken
}

func (c *WorkWechatClient) GetAccessToken(ctx context.Context) (auth2.IAccessToken, error) {
	if c.token != nil {
		return c.token, nil
	}

	endpoint := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", c.CorpId, c.AppSecret)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	var token AccessToken
	if err := auth2.Request(req, &token); err != nil {
		return token, err
	}
	token.createTime = time.Now()
	return token, nil
}

func (c *WorkWechatClient) SetAccessToken(token auth2.IAccessToken) {
	c.token = token
	return
}

func (c *WorkWechatClient) BuildURL(callback string, isAppBrowser bool) string {
	var endpoint string
	if isAppBrowser {
		// 企业微信内-网页授权登录
		urlFmt := "%s?appid=%s&agentid=%s&redirect_uri=%s&response_type=code&scope=snsapi_privateinfo&state=%s#wechat_redirect"
		endpoint = fmt.Sprintf(urlFmt, auth2Url, c.CorpId, c.AppId, url.PathEscape(callback), callbackState)
	} else {
		// 浏览器内-扫码授权登录
		urlFmt := "%s?login_type=CorpApp&appid=%s&agentid=%s&redirect_uri=%s&state=%s"
		endpoint = fmt.Sprintf(urlFmt, ssoUrl, c.CorpId, c.AppId, url.PathEscape(callback), callbackState)
	}
	return endpoint
}

func (c *WorkWechatClient) ValidateCallback(state string) error {
	if state != callbackState {
		return errors.New("auth2.state.wrong")
	}
	return nil
}

func (c *WorkWechatClient) getUserId(ctx context.Context, code string) (UserIdResponse, error) {
	var userId UserIdResponse

	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return userId, err
	}
	endpoint := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo?access_token=%s&code=%s", token.GetToken(), code)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err := auth2.Request(req, &userId); err != nil {
		return userId, err
	}

	if userId.UserId == "" {
		return userId, errors.New("auth2.userid.empty")
	}

	return userId, nil
}

func (c *WorkWechatClient) getUserInfo(ctx context.Context, userid string) (UserInfoResponse, error) {
	var userInfo UserInfoResponse
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return userInfo, err
	}

	endpoint := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s", token.GetToken(), userid)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)

	if err := auth2.Request(req, &userInfo); err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (c *WorkWechatClient) getUserPrivateInfo(ctx context.Context, ticket string) (UserPrivateInfoResponse, error) {
	var userInfo UserPrivateInfoResponse

	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return userInfo, err
	}
	endpoint := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/auth/getuserdetail?access_token=%s", token.GetToken())

	b, _ := json.Marshal(map[string]string{
		"user_ticket": ticket,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(b))

	if err := auth2.Request(req, &userInfo); err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func (c *WorkWechatClient) GetUserInfo(ctx context.Context, code string) (auth2.UserInfo, error) {
	var info auth2.UserInfo

	userid, err := c.getUserId(ctx, code)
	if err != nil {
		return info, err
	}

	userInfo, err := c.getUserInfo(ctx, userid.UserId)
	if err != nil {
		return info, err
	}

	info.UserId = userInfo.UserId
	info.Name = userInfo.Name

	if userid.UserTicket == "" {
		return info, nil
	}

	private, err := c.getUserPrivateInfo(ctx, userid.UserTicket)
	if err != nil {
		return info, err
	}

	info.Mail = private.BizMail
	info.Avatar = private.Avatar
	info.Mobile = private.Mobile
	return info, nil
}
