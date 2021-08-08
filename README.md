# MinDoc 简介

[![Build Status](https://travis-ci.com/mindoc-org/mindoc.svg?branch=master)](https://travis-ci.com/mindoc-org/mindoc)
[![Build status](https://ci.appveyor.com/api/projects/status/o3lcfmf5iy2cp9m6?svg=true)](https://ci.appveyor.com/project/gsw945/mindoc)

MinDoc 是一款针对IT团队开发的简单好用的文档管理系统。

MinDoc 的前身是 [SmartWiki](https://github.com/lifei6671/SmartWiki) 文档系统。SmartWiki 是基于 PHP 框架 laravel 开发的一款文档管理系统。因 PHP 的部署对普通用户来说太复杂，所以改用 Golang 开发。可以方便用户部署和实用。

开发缘起是公司IT部门需要一款简单实用的项目接口文档管理和分享的系统。其功能和界面源于 kancloud 。

可以用来储存日常接口文档，数据库字典，手册说明等文档。内置项目管理，用户管理，权限管理等功能，能够满足大部分中小团队的文档管理需求。

##### 演示站点:
- [https://www.iminho.me/wiki/](https://www.iminho.me/wiki/)
- https://doc.gsw945.com/

---

### 开发&维护&使用 悉知

感谢作者 [lifei6671](https://github.com/lifei6671) 创造了MinDoc，并持续维护了很久。

作者因工作等原因，精力有限，无法花费足够的时间来持续维护mindoc，已于北京时间2021年3月23日将mindoc交给社区(github组织[mindoc-org](https://github.com/mindoc-org))维护，期待热心开发者加入[mindoc-org](https://github.com/mindoc-org)一起来维护MinDoc。

遇到问题请提 [Issues](https://github.com/mindoc-org/mindoc/issues )，欢迎使用者和贡献者加入QQ群 `1051164153`
<a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=bHFR7P3Qp1nsSPbsTw4KN_ZpFLUAblIU&jump_from=webapi"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="MinDoc使用&amp;开发交流群" title="MinDoc使用&amp;开发交流群"></a>

对开发感兴趣请关注 [Development](https://github.com/mindoc-org/mindoc/projects/1):
- [Todo List](https://github.com/mindoc-org/mindoc/projects/1#column-13554511)
- [Work in progress](https://github.com/mindoc-org/mindoc/projects/1#column-13554512)
- [Review in progress](https://github.com/mindoc-org/mindoc/projects/1#column-13554513)

---

# 安装与使用

**如果你的服务器上没有安装golang程序请手动设置一个环境变量如下：键名为 ZONEINFO，值为MinDoc跟目录下的/lib/time/zoneinfo.zip 。**

更多信息请查看手册： [MinDoc 使用手册](https://www.iminho.me/wiki/docs/mindoc/mindoc-summary.md)

对于没有Golang使用经验的用户，可以从 [https://github.com/mindoc-org/mindoc/releases](https://github.com/mindoc-org/mindoc/releases) 这里下载编译完的程序。

如果有Golang开发经验，建议通过编译安装，要求golang版本不小于1.13(需支持`CGO`和`go mod`)。

```bash
# 克隆源码
git clone https://github.com/mindoc-org/mindoc.git
# go包安装
go mod tidy
# 编译(sqlite需要CGO支持)
go build -ldflags "-w"
# 数据库初始化(此步骤执行之前，需配置`conf/app.conf`)
./mindoc install
# 执行
./mindoc
```

MinDoc 如果使用MySQL储存数据，则编码必须是`utf8mb4_general_ci`。请在安装前，把数据库配置填充到项目目录下的 `conf/app.conf` 中。

如果使用 `SQLite` 数据库，则直接在配置文件中配置数据库路径即可.

如果conf目录下不存在 `app.conf` 请重命名 `app.conf.example` 为 `app.conf`。

**默认程序会自动初始化一个超级管理员用户：admin 密码：123456 。请登录后重新设置密码。**


```bash

#邮件配置-示例
#是否启用邮件
enable_mail=true
#smtp服务器的账号
smtp_user_name=admin@iminho.me
#smtp服务器的地址
smtp_host=smtp.ym.163.com
#密码
smtp_password=1q2w3e__ABC
#端口号
smtp_port=25
#邮件发送人的地址
form_user_name=admin@iminho.me
#邮件有效期30分钟
mail_expired=30
```


# 使用Docker部署
如果是Docker用户，可参考项目内置的Dockerfile文件自行编译镜像(编译命令见Dockerfile文件底部注释，仅供参考)。

在启动镜像时需要提供如下的常用环境变量(全部支持的环境变量请参考: [`conf/app.conf.example`](https://github.com/mindoc-org/mindoc/blob/master/conf/app.conf.example))：
```ini
DB_ADAPTER                  指定DB类型(默认为sqlite)
MYSQL_PORT_3306_TCP_ADDR    MySQL地址
MYSQL_PORT_3306_TCP_PORT    MySQL端口号
MYSQL_INSTANCE_NAME         MySQL数据库名称
MYSQL_USERNAME              MySQL账号
MYSQL_PASSWORD              MySQL密码
HTTP_PORT                   程序监听的端口号
MINDOC_ENABLE_EXPORT        开启导出(默认为false)
```

### 举个栗子-当前(公开)镜像(信息页面: https://cr.console.aliyun.com/images/cn-hangzhou/mindoc-org/mindoc/detail)
##### Windows
```bash
set MINDOC=//d/mindoc
docker run -it --name=mindoc --restart=always -v "%MINDOC%":"/mindoc-sync-host" -p 8181:8181 -e MINDOC_ENABLE_EXPORT=true -d registry.cn-hangzhou.aliyuncs.com/mindoc-org/mindoc:v2.1-beta.4
```

##### Linux、Mac
```bash
export MINDOC=/home/ubuntu/mindoc-docker
docker run -it --name=mindoc --restart=always -v "${MINDOC}":"/mindoc-sync-host" -p 8181:8181 -e MINDOC_ENABLE_EXPORT=true -d registry.cn-hangzhou.aliyuncs.com/mindoc-org/mindoc:v2.1-beta.4
```

##### 举个栗子-更多环境变量示例(镜像已过期，仅供参考，请以当前镜像为准)
```bash
docker run -p 8181:8181 --name mindoc -e DB_ADAPTER=mysql -e MYSQL_PORT_3306_TCP_ADDR=10.xxx.xxx.xxx -e MYSQL_PORT_3306_TCP_PORT=3306 -e MYSQL_INSTANCE_NAME=mindoc -e MYSQL_USERNAME=root -e MYSQL_PASSWORD=123456 -e httpport=8181 -d daocloud.io/lifei6671/mindoc:latest
```

### docker-compose 一键安装

1. 修改配置文件
    修改`docker-compose.yml`中的配置信息，主要修改`volumes`节点，将宿主机的两个目录映射到容器内。
    `environment`节点，配置自己的环境变量。
    
2. 一键完成所有环境搭建
    
    > docker-compose up -d
3. 浏览器访问
    > http://localhost:8181/

    整个部署完成了
4. 常用命令参考
   - 启动
        
        > docker-compose up -d
   - 停止
        
        > docker-compose stop
   - 重启
        
        > docker-compose restart
   - 停止删除容器，释放所有资源
        
        > docker-compose down
   - 删除并重新创建
        > docker-compose -f docker-compose.yml down && docker-compose up -d
        > 
        > 更多 docker-compose 的使用相关的内容 请查看官网文档或百度
   
# 项目截图

**创建项目**

![创建项目](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501204438.png)

**项目列表**

![项目列表](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501203542.png)

**项目概述**

![项目概述](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501203619.png)

**项目成员**

![项目成员](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501203637.png)

**项目设置**

![项目设置](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501203656.png)

**基于Editor.md开发的Markdown编辑器**

![基于Editor.md开发的Markdown编辑器](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501203854.png)

**基于wangEditor开发的富文本编辑器**

![基于wangEditor开发的富文本编辑器](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501204651.png)

**项目预览**

![项目预览](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501204609.png)

**超级管理员后台**

![超级管理员后台](https://raw.githubusercontent.com/lifei6671/mindoc/master/uploads/20170501204710.png)


# 使用的技术(TODO: 最新技术栈整理中，使用的第三方库升级中)

- [Beego](https://github.com/beego/beego) ~~1.10.0~~
- MySQL 5.6
- [editor.md](https://github.com/pandao/editor.md) Markdown 编辑器
- [Bootstrap](https://github.com/twbs/bootstrap) 3.2
- [jQuery](https://github.com/jquery/jquery) 库
- [WebUploader](https://github.com/fex-team/webuploader) 文件上传框架
- [NProgress](https://github.com/rstacruz/nprogress) 库
- [jsTree](https://github.com/vakata/jstree) 树状结构库
- [Font Awesome](https://github.com/FortAwesome/Font-Awesome) 字体库
- [Cropper](https://github.com/fengyuanchen/cropper) 图片剪裁库
- [layer](https://github.com/sentsin/layer) 弹出层框架
- [highlight.js](https://github.com/highlightjs/highlight.js) 代码高亮库
- ~~to-markdown~~[Turndown](https://github.com/domchristie/turndown) HTML转Markdown库
- ~~quill 富文本编辑器~~
- [wangEditor](https://github.com/wangeditor-team/wangEditor) 富文本编辑器
- [Vue.js](https://github.com/vuejs/vue) 框架


# 主要功能

- 项目管理，可以对项目进行编辑更改，成员添加等。
- 文档管理，添加和删除文档等。
- 评论管理，可以管理文档评论和自己发布的评论。
- 用户管理，添加和禁用用户，个人资料更改等。
- 用户权限管理 ， 实现用户角色的变更。
- 项目加密，可以设置项目公开状态，私有项目需要通过Token访问。
- 站点配置，可开启匿名访问、验证码等。

# 参与开发

我们欢迎您在 MinDoc 项目的 GitHub 上报告 issue 或者 pull request。

如果您还不熟悉GitHub的Fork and Pull开发模式，您可以阅读GitHub的文档（https://help.github.com/articles/using-pull-requests） 获得更多的信息。

# 关于作者[lifei6671](https://github.com/lifei6671)

一个不纯粹的PHPer，一个不自由的 gopher 。

