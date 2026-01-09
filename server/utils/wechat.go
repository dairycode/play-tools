package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"play-tools/config"
)

// WechatSession 微信登录响应结构
type WechatSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// GetWechatOpenID 通过 code 获取微信用户的 openid
func GetWechatOpenID(code string) (*WechatSession, error) {
	appID := config.AppConfig.Wechat.AppID
	appSecret := config.AppConfig.Wechat.Secret

	if appID == "" || appSecret == "" {
		return nil, fmt.Errorf("微信配置未设置，请在 config.yaml 中配置 wechat.appid 和 wechat.secret")
	}

	// 微信 code2Session 接口
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appID, appSecret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求微信接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var session WechatSession
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查是否有错误
	if session.ErrCode != 0 {
		return nil, fmt.Errorf("微信登录失败: %s (错误码: %d)", session.ErrMsg, session.ErrCode)
	}

	if session.OpenID == "" {
		return nil, fmt.Errorf("获取 OpenID 失败")
	}

	return &session, nil
}
