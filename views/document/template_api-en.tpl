#### Introduction：

- User Login

#### Version：

|ver|Maker|Set date|Revision date|
|:----    |:---|:----- |-----   |
|2.1.0 |Ben  |2021-04-15 |  xxxx-xx-xx |

#### 请求URL:

- http://xx.com/api/login

#### Request Method：

- GET
- POST

#### Request Header：

|Paramter|Must|Type|Description|
|:----    |:---|:----- |-----   |
|Content-Type |Y  |string |request type： application/json   |
|Content-MD5 |Y  |string | request sign    |


#### request paramters:

|Paramter|Must|Type|description|
|:----    |:---|:----- |-----   |
|username |Y  |string |   |
|password |Y  |string |    |

#### Return Example:

**Correct Return:**

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

**Error Return:**


```
{
    "errcode": 500,
    "errmsg": "invalid appid"
}
```

#### Return Paramters:

|Paramter|Type|description|
|:-----  |:-----|-----                           |
|group_level |int   |user group id，1：Super administrator；2：normal  |

#### Remark:

- Check Error Code for more return error message