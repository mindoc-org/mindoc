# MinDoc 简介

[![Build Status](https://travis-ci.org/lifei6671/mindoc.svg?branch=master)](https://travis-ci.org/lifei6671/mindoc)
[![Build status](https://ci.appveyor.com/api/projects/status/ik70whjrioyvfy18/branch/master?svg=true)](https://ci.appveyor.com/project/lifei6671/godoc/branch/master)

MinDoc 是一款针对IT团队开发的简单好用的文档管理系统。

MinDoc 的前身是 SmartWiki 文档系统。SmartWiki 是基于 PHP 框架 laravel 开发的一款文档管理系统。因 PHP 的部署对普通用户来说太复杂，所以改用 Golang 开发。可以方便用户部署和实用。

开发缘起是公司IT部门需要一款简单实用的项目接口文档管理和分享的系统。其功能和界面源于 kancloud 。

可以用来储存日常接口文档，数据库字典，手册说明等文档。内置项目管理，用户管理，权限管理等功能，能够满足大部分中小团队的文档管理需求。

演示站点： [http://doc.iminho.me](http://doc.iminho.me)

# 安装与使用

**如果你的服务器上没有安装golang程序请手动设置一个环境变量如下：键名为 ZONEINFO，值为MinDoc跟目录下的/lib/time/zoneinfo.zip 。**

更多信息请查看手册： [MinDoc 使用手册](https://github.com/lifei6671/mindoc/wiki)

对于没有Golang使用经验的用户，可以从 [https://github.com/lifei6671/mindoc/releases](https://github.com/lifei6671/mindoc/releases) 这里下载编译完的程序。

如果有Golang开发经验，建议通过编译安装，在此之前，您需要先安装Golang官方包管理工具，详见[Install dep](https://golang.github.io/dep/docs/installation.html)。

```bash
git clone https://github.com/lifei6671/mindoc.git

dep ensure

go build -ldflags "-w"

./mindoc install

./mindoc

```

MinDoc 使用MySQL储存数据，且编码必须是`utf8mb4_general_ci`。请在安装前，把数据库配置填充到项目目录下的 conf/app.conf 中。

如果conf目录下不存在 app.conf 请重命名 app.conf.example 为 app.conf。

**默认程序会自动初始化一个超级管理员用户：admin 密码：123456 。请登录后重新设置密码。**


```bash

#邮件配置
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
如果是Docker用户，可参考项目内置的Dockerfile文件编译镜像。

在启动镜像时需要提供如下的环境变量：

```ini
DB_ADAPTER                  制定 DB
MYSQL_PORT_3306_TCP_ADDR    MySQL地址
MYSQL_PORT_3306_TCP_PORT    MySQL端口号
MYSQL_INSTANCE_NAME         MySQL数据库名称
MYSQL_USERNAME              MySQL账号
MYSQL_PASSWORD              MySQL密码
HTTP_PORT                   程序监听的端口号
```

举个栗子

```bash
docker run -p 8181:8181 --name mindoc -e DB_ADAPTER=mysql -e MYSQL_PORT_3306_TCP_ADDR=10.xxx.xxx.xxx -e MYSQL_PORT_3306_TCP_PORT=3306 -e MYSQL_INSTANCE_NAME=mindoc -e MYSQL_USERNAME=root -e MYSQL_PASSWORD=123456 -e httpport=8181 -d daocloud.io/lifei6671/mindoc:latest
```
### docker-compose 一键安装
1. 生成配置文件
    > cp env-example .env

    按照自己的环境填写配置信息到.env文件中
2. 一键完成所有环境搭建
    > docker-compose up -d mindoc
3. 浏览器访问
    > http://localhost:8181/

    整个部署完成了
4. 常用命令参考
   - 启动
        > docker-compose up -d mindoc
   - 停止
        > docker-compose stop mindoc
   - 重启
        > docker-compose restart mindoc
   - 停止删除容器，释放所有资源
        > docker-compose down

   更多 docker-compose 的使用相关的内容 请查看官网文档或百度
   
# 支持 MinDoc

![支付宝](https://raw.githubusercontent.com/lifei6671/mindoc/master/static/images/alipay.png) ![微信支付](https://raw.githubusercontent.com/lifei6671/mindoc/master/static/images/weixin.png)

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


# 使用的技术

- beego 1.8.1
- mysql 5.6
- editor.md
- bootstrap 3.2
- jquery 库
- layer 弹出层框架
- webuploader 文件上传框架
- Nprogress 库
- jstree 树状结构库
- font awesome 字体库
- cropper 图片剪裁库
- layer 弹出层框架
- highlight 代码高亮库
- to-markdown HTML转Markdown库
- wangEditor 富文本编辑器


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

# 关于作者

一个不纯粹的PHPer，一个不自由的 gopher 。
