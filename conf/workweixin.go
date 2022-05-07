package conf

import (
    "github.com/beego/beego/v2/server/web"
)

type WorkWeixinConf struct {
    CorpId        string // 企业ID
    AgentId       string // 应用ID
    Secret        string // 应用密钥
    ContactSecret string // 通讯录密钥
}

func GetWorkWeixinConfig() *WorkWeixinConf {
    corpid, _ := web.AppConfig.String("workweixin_corpid")
    agentid, _ := web.AppConfig.String("workweixin_agentid")
    secret, _ := web.AppConfig.String("workweixin_secret")
    contact_secret, _ := web.AppConfig.String("workweixin_contact_secret")

    c := &WorkWeixinConf{
        CorpId:        corpid,
        AgentId:       agentid,
        Secret:        secret,
        ContactSecret: contact_secret,
    }
    return c
}
