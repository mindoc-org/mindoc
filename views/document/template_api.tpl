#### 简要描述：

- 用户登录接口

#### 接口版本：

|版本号|制定人|制定日期|修订日期|
|:----    |:---|:----- |-----   |
|2.1.0 |秦亮  |2017-03-20 |  xxxx-xx-xx |

#### 请求URL:

- http://xx.com/api/login

#### 请求方式：

- GET
- POST

#### 请求头：

|参数名|是否必须|类型|说明|
|:----    |:---|:----- |-----   |
|Content-Type |是  |string |请求类型： application/json   |
|Content-MD5 |是  |string | 请求内容签名    |


#### 请求参数:

|参数名|是否必须|类型|说明|
|:----    |:---|:----- |-----   |
|username |是  |string |用户名   |
|password |是  |string | 密码    |

#### 返回示例:

**正确时返回:**

```
{
    "errcode": 0,
    "data": {
        "uid": "1",
        "account": "admin",
        "nickname": "Minho",
        "group_level": 0 ,
        "create_time": "1436864169",
        "last_login_time": "0",
    }
}
```

**错误时返回:**


```
{
    "errcode": 500,
    "errmsg": "invalid appid"
}
```

#### 返回参数说明:

|参数名|类型|说明|
|:-----  |:-----|-----                           |
|group_level |int   |用户组id，1：超级管理员；2：普通用户；3：只读用户  |

#### 备注:

- 更多返回错误代码请看首页的错误代码描述