# MinDoc 简介

[![Build Status](https://travis-ci.com/mindoc-org/mindoc.svg?branch=master)](https://travis-ci.com/mindoc-org/mindoc)
[![Build status](https://ci.appveyor.com/api/projects/status/7680ia6mu29m12wx?svg=true)](https://ci.appveyor.com/project/mindoc-org/mindoc)

MinDoc 是一款针对IT团队开发的简单好用的文档管理系统。

MinDoc 的前身是 [SmartWiki](https://github.com/lifei6671/SmartWiki) 文档系统。SmartWiki 是基于 PHP 框架 laravel 开发的一款文档管理系统。因 PHP 的部署对普通用户来说太复杂，所以改用 Golang 开发。可以方便用户部署和实用。

开发缘起是公司IT部门需要一款简单实用的项目接口文档管理和分享的系统。其功能和界面源于 kancloud 。

可以用来储存日常接口文档，数据库字典，手册说明等文档。内置项目管理，用户管理，权限管理等功能，能够满足大部分中小团队的文档管理需求。

##### 演示站点&文档:
- https://www.iminho.me/wiki/docs/mindoc/
- https://doc.gsw945.com/docs/mindoc-docs/

---

### 开发&维护&使用 悉知

- 感谢作者 [lifei6671](https://github.com/lifei6671) 创造了MinDoc，并持续维护了很久。
- 作者因工作等原因，精力有限，无法花费足够的时间来持续维护mindoc，已于北京时间2021年3月23日将mindoc交给社区(github组织[mindoc-org](https://github.com/mindoc-org))维护，期待热心开发者加入[mindoc-org](https://github.com/mindoc-org)一起来维护MinDoc。
- 遇到问题请提 [Issues](https://github.com/mindoc-org/mindoc/issues )，欢迎使用者和贡献者加入QQ群 `1051164153`
<a target="_blank" href="https://qm.qq.com/cgi-bin/qm/qr?k=bHFR7P3Qp1nsSPbsTw4KN_ZpFLUAblIU&jump_from=webapi"><img border="0" src="https://pub.idqqimg.com/wpa/images/group.png" alt="MinDoc使用&amp;开发交流群" title="MinDoc使用&amp;开发交流群"></a>
- 对开发感兴趣请关注 [Development](https://github.com/mindoc-org/mindoc/projects/1):
  - [Todo List](https://github.com/mindoc-org/mindoc/projects/1#column-13554511)
  - [Work in progress](https://github.com/mindoc-org/mindoc/projects/1#column-13554512)
  - [Review in progress](https://github.com/mindoc-org/mindoc/projects/1#column-13554513)
- Mindoc基于 [beeego](https://github.com/beego/beego) 开发，beego文档地址: https://github.com/beego/beego-doc/tree/main/docs/zh
- :warning: **特别声明**:
  - 原作者 [lifei6671](https://github.com/lifei6671) 已于 2021-08-06 删除了个人捐赠信息，参见: [1a179179c1fe4d0d4db95e0b757d863aee5bf395](https://github.com/mindoc-org/mindoc/commit/1a179179c1fe4d0d4db95e0b757d863aee5bf395)
  - 截止目前(2023-03-27)，[mindoc-org](https://github.com/mindoc-org) 暂未发布任何捐赠信息，请勿轻信

---

# 安装与使用

~~如果你的服务器上没有安装golang程序请手动设置一个环境变量如下：键名为 ZONEINFO，值为MinDoc跟目录下的/lib/time/zoneinfo.zip 。~~

更多信息请查看手册： [MinDoc 使用手册](https://www.iminho.me/wiki/docs/mindoc/mindoc-summary.md)

对于没有Golang使用经验的用户，可以从 [https://github.com/mindoc-org/mindoc/releases](https://github.com/mindoc-org/mindoc/releases) 这里下载编译完的程序。

如果有Golang开发经验，建议通过编译安装，要求golang版本不小于1.15.1(需支持`CGO`、`go mod`和`import _ "time/tzdata"`)(推荐Go版本为1.18.1)。
> 注意: CentOS7上GLibC版本低，常规编译版本不能使用。需要自行源码编译,或使用使用musl编译版本。

## 常规编译
```bash
# 克隆源码
git clone https://github.com/mindoc-org/mindoc.git
# go包安装
go mod tidy -v
# 编译(sqlite需要CGO支持)
go build -ldflags "-w" -o mindoc main.go
# 数据库初始化(此步骤执行之前，需配置`conf/app.conf`)
./mindoc install
# 执行
./mindoc
# 开发阶段运行
bee run
```

MinDoc 如果使用MySQL储存数据，则编码必须是`utf8mb4_general_ci`。请在安装前，把数据库配置填充到项目目录下的 `conf/app.conf` 中。

如果使用 `SQLite` 数据库，则直接在配置文件中配置数据库路径即可.

如果conf目录下不存在 `app.conf` 请重命名 `app.conf.example` 为 `app.conf`。

**默认程序会自动初始化一个超级管理员用户：admin 密码：123456 。请登录后重新设置密码。**

## Linux系统中不依赖gLibC的编译方式

### 安装 musl-gcc
```bash
wget -c http://musl.libc.org/releases/musl-1.2.2.tar.gz
tar -xvf musl-1.2.2.tar.gz
cd musl-1.2.2
./configure
make
sudo make install
```
### 使用 musl-gcc 编译 mindoc
```bash
go mod tidy -v
export GOARCH=amd64
export GOOS=linux
# 设置使用musl-gcc
export CC=/usr/local/musl/bin/musl-gcc
# 设置版本
export TRAVIS_TAG=temp-musl-v`date +%y%m%d`
go build -v -o mindoc_linux_musl_amd64 -ldflags="-linkmode external -extldflags '-static' -w -X 'github.com/mindoc-org/mindoc/conf.VERSION=$TRAVIS_TAG' -X 'github.com/mindoc-org/mindoc/conf.BUILD_TIME=`date`' -X 'github.com/mindoc-org/mindoc/conf.GO_VERSION=`go version`'"
# 验证
./mindoc_linux_musl_amd64 version
```

## Windows 上后台运行
 使用 [mindoc-daemon](https://github.com/mindoc-org/mindoc-daemon)


```ini
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

#### 举个栗子-当前(公开)镜像(信息页面: https://cr.console.aliyun.com/images/cn-hangzhou/mindoc-org/mindoc/detail , 需要登录阿里云账号才可访问列表)
##### Windows
```bash
set MINDOC=//d/mindoc
docker run -it --name=mindoc --restart=always -v "%MINDOC%/conf":"/mindoc/conf" -p 8181:8181 -e MINDOC_ENABLE_EXPORT=true -d registry.cn-hangzhou.aliyuncs.com/mindoc-org/mindoc:v2.1
```

##### Linux、Mac
```bash
export MINDOC=/home/ubuntu/mindoc-docker
docker run -it --name=mindoc --restart=always -v "${MINDOC}/conf":"/mindoc/conf" -p 8181:8181 -e MINDOC_ENABLE_EXPORT=true -d registry.cn-hangzhou.aliyuncs.com/mindoc-org/mindoc:v2.1
```

##### 举个栗子-更多环境变量示例(镜像已过期，仅供参考，请以当前镜像为准)
```bash
docker run -p 8181:8181 --name mindoc -e DB_ADAPTER=mysql -e MYSQL_PORT_3306_TCP_ADDR=10.xxx.xxx.xxx -e MYSQL_PORT_3306_TCP_PORT=3306 -e MYSQL_INSTANCE_NAME=mindoc -e MYSQL_USERNAME=root -e MYSQL_PASSWORD=123456 -e httpport=8181 -d daocloud.io/lifei6671/mindoc:latest
```

#### dockerfile内容参考
- [无需代理直接加速各种 GitHub 资源拉取 | 国内镜像赋能 | 助力开发](https://blog.frytea.com/archives/504/)
- [阿里云 - Ubuntu 镜像](https://developer.aliyun.com/mirror/ubuntu)

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
  - 参考
    - [wangEditor v4.7 富文本编辑器教程](https://www.bookstack.cn/books/wangeditor-4.7-zh)
    - [扩展菜单注册太过繁琐 #2493](https://github.com/wangeditor-team/wangEditor/issues/2493)
  - 工具： `https://babeljs.io/repl` + `@babel/plugin-transform-classes`
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


# 部署补充
- 若内网部署，draw.io无法使用外网，则需要用tomcat运行war包，见（https://github.com/jgraph/drawio） 从release下载，之后修改markdown.js的TODO行对应的链接即可
- 为了护眼，简单增加了编辑界面的主题切换，见editormd.js和markdown_edit_template.tpl
- (需重新编译项)为了对已删除文档/文档引用图片删除文字后，对悬空无引用的图片/附件进行清理，增加了清理接口，需重新编译
     - 编译后除二进制文件外还需更新三个文件: conf/lang/en-us.ini,zh-cn.ini; attach_list.tpl
     - 若不想重新编译，也可通过database/clean.py，手动执行对无引用图片/附件的文件清理和数据库记录双向清理。
- 若采用nginx二级部署，以yourpath/为例，需修改
     - conf/app.conf修改：`baseurl="/yourpath"`
     - static/js/kancloud.js文件中`url: "/comment/xxxxx` => `url: "/yourpath" + "/comment/xxxxx`, 共两处

     - nginx端口代理示例:
     ```
     增加
     location  /yourpath/ {
          rewrite ^/yourpath/(.*) /$1  break;
          proxy_pass http://127.0.0.1:8181;
     }
     ```
     注意使用的是127.0.0.1，根据自身选择替换，如果nginx是docker部署，则还需要在docker中托管运行mindoc，具体参考如下配置:
     - docker-compose代理示例(docker-nginx代理运行mindoc)
     ```
     version: '3'
     services:
       mynginx:
       image: nginx:latest
       ports:
         - "8880:80"
       command: 
         - bash
         - -c
         - |
             service nginx start
             cd /src/mindoc/ && ./mindoc
       volumes:
         - ..:/src
         - ./nginx:/etc/nginx/conf.d
     ```

     目录结构
     ```
     onefolder
     |
       - docker
       |
         - docker-compose.yml
         - nginx
         |
           - mynginx.conf
       
       - mindoc
       |
         - database/
         - conf/
         - ...
     ```

---

# 文档访问密码（本次开发）

为满足「单个目录（文章）单独设置访问密码」的需求，在原有项目级加密之外，新增了文档级访问密码功能。

## 功能说明

- 后台编辑某篇文章时，可勾选「访问密码保护」并设置自定义「访问密码」（密码为空时不允许开启保护）。
- 普通用户打开该文章时需输入密码才能阅读，同一会话（Session）内再次访问无需重复输入。
- 支持 Markdown（Editor.md）与 Cherry Markdown 两种编辑器的文档新建/编辑弹窗。
- 管理员（`Role == 0 || Role == 1`）可直接访问任何设置了密码的文档，无需输入密码。

## 使用步骤

1. 进入项目编辑页 → 新建或编辑某篇文章。
2. 勾选「访问密码保护」→ 在「访问密码」输入框填写密码（编辑时留空表示保持原密码）。
3. 保存后，普通用户点击该文章会跳转到密码输入页，输入正确密码即可阅读。

## 涉及改动

- 数据模型：`models/DocumentModel.go` 新增 `IsPassword`（是否启用）与 `Password`（访问密码）字段。
- 数据库：`md_documents` 表新增 `is_password`、`password` 两列。
- 控制器：`controllers/DocumentController.go` 新增 `CheckDocumentPassword`、`hasDocumentAccess`、`renderDocumentPasswordPage` 等逻辑，并在 `Read`/`Index` 中拦截密码保护文档。
- 路由：`routers/router.go` 新增 `/docs/:key/:id/check-password`。
- 前端：
  - `views/document/document_doc_password.tpl` 新增文档密码输入页。
  - `views/document/markdown_edit_template.tpl`、`cherry_markdown_edit_template.tpl` 新增密码复选框与密码输入框。
  - `static/js/editor.js` 处理密码输入框的创建/编辑回填与重置。
  - `static/js/kancloud.js` 处理 AJAX 加载加密文档时跳转到整页密码输入页（`errcode === 6010`）。
- 国际化：`conf/lang/zh-cn.ini`、`en-us.ini` 新增相关文案。

## 备注

- 密码当前以明文保存（与项目级 `BookPassword` 一致）。
- 修改了模型字段，部署时需重新编译并在数据库执行对应列的新增（或使用 `./mindoc install` 重新初始化）。