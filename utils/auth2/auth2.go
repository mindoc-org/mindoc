package auth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserInfo struct {
	UserId string `json:"userid"` // 企业成员userid
	Name   string `json:"name"`   // 姓名
	Avatar string `json:"avatar"` // 头像
	Mobile string `json:"mobile"` // 手机号
	Mail   string `json:"mail"`   // 邮箱
}

func NewAccessToken(token IAccessToken) AccessTokenCache {
	return AccessTokenCache{
		AccessToken: token.GetToken(),
		ExpireIn:    token.GetExpireIn(),
		ExpireTime:  token.GetExpireTime(),
	}
}

type AccessTokenCache struct {
	ExpireIn    time.Duration
	ExpireTime  time.Time
	AccessToken string
}

func (a AccessTokenCache) GetToken() string {
	return a.AccessToken
}

func (a AccessTokenCache) GetExpireIn() time.Duration {
	return a.ExpireIn
}

func (a AccessTokenCache) GetExpireTime() time.Time {
	return a.ExpireTime
}

func (a AccessTokenCache) IsExpired() bool {
	return time.Now().After(a.ExpireTime)
}

type IAccessToken interface {
	GetToken() string
	GetExpireIn() time.Duration
	GetExpireTime() time.Time
}

type Client interface {
	GetAccessToken(ctx context.Context) (IAccessToken, error)
	SetAccessToken(token IAccessToken)
	BuildURL(callback string, isAppBrowser bool) string
	ValidateCallback(state string) error
	GetUserInfo(ctx context.Context, code string) (UserInfo, error)
}

type IResponse interface {
	AsError() error
}

func Request(req *http.Request, v IResponse) error {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("status = %d, msg = %s", response.StatusCode, string(b))
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return v.AsError()
}
