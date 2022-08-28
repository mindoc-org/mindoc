package workweixin

import (
	"context"
	"crypto/tls"
	// "encoding/json"
	"net/http"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/cache"
	"github.com/mindoc-org/mindoc/conf"
)

// doc
// - 全局错误码: https://work.weixin.qq.com/api/doc/90000/90139/90313

const (
	AccessTokenCacheKey        = "access-token-cache-key"
	ContactAccessTokenCacheKey = "contact-access-token-cache-key"
)

// 获取访问凭据-请求响应结构
type AccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"` // 获取到的凭证,最长为512字节
	ExpiresIn   int    `json:"expires_in"`   // 凭证的有效时间（秒）
}

// 获取用户Id-请求响应结构
type UserIdResponse struct {
	ErrCode  int    `json:"errcode"`
	ErrMsg   string `json:"errmsg"`
	UserId   string `json:"UserId"`   // 企业成员UserID
	OpenId   string `json:"OpenId"`   // 非企业成员的标识，对当前企业唯一
	DeviceId string `json:"DeviceId"` // 设备号
}

// 获取成员ID列表-请求响应结构
type UserListIdResponse struct {
	ErrCode    int                      `json:"errcode"`
	ErrMsg     string                   `json:"errmsg"`
	NextCursor string                   `json:"next_cursor"` // 分页游标,下次请求时填写以获取之后分页的记录
	DeptUser   []WorkWeixinDeptUserInfo `json:"dept_user"`   // 用户-部门关系列表

}

// 获取用户信息-请求响应结构
type UserInfoResponse struct {
	ErrCode        int    `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
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

// 访问凭据缓存-结构
type AccessTokenCache struct {
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	UpdateTime  time.Time `json:"update_time"`
}

// 企业微信用户信息-结构
type WorkWeixinDeptUserInfo struct {
	UserId     string `json:"UserId"`     // 企业成员UserID
	Department int    `json:"department"` // 成员所属部门id列表
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

func httpFilter(next httplib.Filter) httplib.Filter {
	return func(ctx context.Context, req *httplib.BeegoHTTPRequest) (*http.Response, error) {
		r := req.GetRequest()
		logs.Info("filter-url: ", r.URL)
		// Never forget invoke this. Or the request will not be sent
		return next(ctx, req)
	}
}

// 获取访问凭据-请求
func RequestAccessToken(corpid string, secret string) (cache_token AccessTokenCache, ok bool) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	req := httplib.Get(url)
	req.Param("corpid", corpid)     // 企业ID
	req.Param("corpsecret", secret) // 应用的凭证密钥
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	req.AddFilters(httpFilter)
	resp, err := req.Response()
	_ = resp
	var token AccessTokenCache
	if err != nil {
		logs.Error(err)
		return token, false
	}
	var atr AccessTokenResponse
	err = req.ToJSON(&atr)
	if err != nil {
		logs.Error(err)
		return token, false
	}
	token = AccessTokenCache{
		AccessToken: atr.AccessToken,
		ExpiresIn:   atr.ExpiresIn,
		UpdateTime:  time.Now(),
	}
	return token, true
}

// 获取访问凭据
func GetAccessToken(is_contact bool) (access_token string, ok bool) {
	var cache_token AccessTokenCache
	cache_key := AccessTokenCacheKey
	if is_contact {
		cache_key = ContactAccessTokenCacheKey
	}
	err := cache.Get(cache_key, &cache_token)
	if err == nil {
		logs.Info("AccessToken从缓存读取成功")
		// TODO: access_token有效期判断, 刷新
		return cache_token.AccessToken, true
	} else {
		logs.Warning(err)
		workweixinConfig := conf.GetWorkWeixinConfig()
		logs.Debug("corp_id: ", workweixinConfig.CorpId)
		logs.Debug("agent_id: ", workweixinConfig.AgentId)
		logs.Debug("secret: ", workweixinConfig.Secret)
		logs.Debug("contact_secret: ", workweixinConfig.ContactSecret)
		secret := workweixinConfig.Secret
		if is_contact {
			secret = workweixinConfig.ContactSecret
		}
		new_token, ok := RequestAccessToken(workweixinConfig.CorpId, secret)
		if ok {
			logs.Debug(new_token)
			if err = cache.Put(cache_key, new_token, time.Second*time.Duration(new_token.ExpiresIn)); err == nil {
				logs.Info("AccessToken缓存写入成功")
				return new_token.AccessToken, true
			}
			logs.Warning("AccessToken缓存写入失败")
			return "", false
		}
		logs.Warning("AccessToken请求失败")
		return "", false
	}
}

// 获取用户id-请求
func RequestUserId(access_token string, code string) (user_id string, ok bool) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
	req := httplib.Get(url)
	req.Param("access_token", access_token) // 应用调用接口凭证
	req.Param("code", code)                 // 通过成员授权获取到的code
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	req.AddFilters(httpFilter)
	resp, err := req.Response()
	_ = resp
	if err != nil {
		logs.Error(err)
		return "", false
	}
	var uir UserIdResponse
	err = req.ToJSON(&uir)
	if err != nil {
		logs.Error(err)
		return "", false
	}
	return uir.UserId, true
}

/*
获取用户详细信息-请求
从2022年8月15日10点开始，“企业管理后台 - 管理工具 - 通讯录同步”的新增IP将不能再调用此接口
url:https://developer.work.weixin.qq.com/document/path/96079
*/
func RequestUserInfo(contact_access_token string, userid string) (user_info WorkWeixinUserInfo, error_msg string, ok bool) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	req := httplib.Get(url)
	req.Param("access_token", contact_access_token) // 通讯录应用调用接口凭证
	req.Param("userid", userid)                     // 成员UserID
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	req.AddFilters(httpFilter)
	resp_str, err := req.String()
	_ = resp_str
	var info WorkWeixinUserInfo
	if err != nil {
		logs.Error(err)
		return info, "请求失败", false
	} else {
		logs.Debug(resp_str)
	}
	var uir UserInfoResponse
	err = req.ToJSON(&uir)
	if err != nil {
		logs.Error(err)
		return info, "请求数据结果错误", false
	}
	if uir.ErrCode != 0 {
		return info, uir.ErrMsg, false
	}
	info = WorkWeixinUserInfo{
		UserId:         uir.UserId,
		Name:           uir.Name,
		HideMobile:     uir.HideMobile,
		Mobile:         uir.Mobile,
		Department:     uir.Department,
		Email:          uir.Email,
		IsLeaderInDept: uir.IsLeaderInDept,
		IsLeader:       uir.IsLeader,
		Avatar:         uir.Avatar,
		Alias:          uir.Alias,
		Status:         uir.Status,
		MainDepartment: uir.MainDepartment,
	}
	return info, "", true
}

/*
获取成员ID列表
*/
func GetUserListId(contact_access_token string, userid string) (user_info WorkWeixinDeptUserInfo, error_msg string, ok bool) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/user/list_id"
	req := httplib.Get(url)
	req.Param("access_token", contact_access_token) // 通讯录应用调用接口凭证
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: false})
	req.AddFilters(httpFilter)
	respStr, err := req.String()
	_ = respStr

	var info WorkWeixinDeptUserInfo
	if err != nil {
		logs.Error(err)
		return info, "请求失败", false
	} else {
		logs.Debug(respStr)
	}

	// 返回响应
	var uir UserListIdResponse
	//获取用户信息失败: 请求数据结果错误
	err = req.ToJSON(&uir)
	if err != nil {
		logs.Error(err)
		return info, "请求数据结果错误", false
	}

	if uir.ErrCode != 0 {
		return info, uir.ErrMsg, false
	}

	// 判断userid 中是否还有当前用户id
	for i := 0; i < len(uir.DeptUser); i++ {
		if uir.DeptUser[i].UserId == userid {
			info = WorkWeixinDeptUserInfo{
				UserId: uir.DeptUser[i].UserId,
			}
		} else {
			return info, uir.ErrMsg, false
		}
	}
	return info, "", true
}
