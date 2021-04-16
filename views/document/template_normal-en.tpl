# MinDoc Introduction

[![Build Status](https://travis-ci.org/lifei6671/mindoc.svg?branch=master)](https://travis-ci.org/lifei6671/mindoc)

MinDoc is a simple and easy-to-use document management system developed for IT teams.

MinDoc The predecessor of Mindoc was the SmartWiki documentation system. SmartWiki is a document management system developed based on the PHP framework laravel. Because the deployment of PHP is too complicated for ordinary users, it is convenient for users to deploy and use Golang instead.

The origin of the development is that the company’s IT department needs a simple and practical system for project interface document management and sharing. Its function and interface are derived from kancloud.

It can be used to store API documents, database dictionaries, manual instructions and other documents. Built-in project management, user management, authority management and other functions can meet the document management needs of most small and medium-sized teams.

Demo site: [http://doc.iminho.me](http://doc.iminho.me)

# Installation and use

**If the golang is not installed on your server, please manually set an environment variable as follows: the key name is ZONEINFO, and the value is MinDoc and /lib/time/zoneinfo.zip in the directory.**

**Windows tutorial:** [https://github.com/mindoc-org/mindoc/blob/master/README_WIN.md](docs/README_WIN.md)

**Linux  tutorial:**  [https://github.com/mindoc-org/mindoc/blob/master/README_LINUX.md](docs/README_LINUX.md)

**PDF Export configuration tutorial**  [https://github.com/mindoc-org/mindoc/blob/master/docs/README_LINUX.md](docs/WKHTMLTOPDF.md)

For users without Golang experience, you can download the compiled program from here. [https://github.com/mindoc-org/mindoc/releases](https://github.com/mindoc-org/mindoc/releases) 

If you have Golang development experience, it is recommended to compile and install.

```bash
git clone https://github.com/mindoc-org/mindoc.git

go get -d ./...

go build -ldflags "-w"

```

MinDoc uses MySQL to store data, and the encoding must be `utf8mb4_general_ci`.

Please Change database config locate in `conf/app.conf` before install.

If `app.conf` does not exist in the conf directory, please rename `app.conf.example` to `app.conf`.

If the install.lock file exists in the root directory, it means that the database has been initialized. If you want to reinitialize the database,  delete the file and restart the program.

**The default program will automatically create the table and initialize a super administrator user: admin password: 123456. Please reset your password after logging in.**

## Linux Running in Background

If you want the program to run in the background, you can execute the following command:

```bash
nohup ./godoc &
```

This command will make the program execute in the background, but the service will not start automatically after the server restarts.

Using supervisor as a service can automatically restart MinDoc after the server restarts.

## Windows Running in Background

Running in the background under Windows requires the help of CMD command line commands：

```bash
# Create slave.vbs file in the MinDoc root directory:

Set ws = CreateObject("Wscript.Shell")
ws.run "cmd /c start.bat",vbhide

# Create start.bat file:
@echo off

godoc_windows_amd64.exe

```

Double-click slave.bat to start, After the program initializes the database, an install.lock file will be created in this directory, indicating that the installation has been successful.

If you compile it yourself, you can use the following command to compile a program that does not rely on the cmd command to run in the background:

```bash
go build -ldflags "-H=windowsgui"
```
Compiled by this command runs in the background by default on Windows.

Please add MinDoc to the boot list.

## Password retrieval

The password retrieval function depends on the mail service. Therefore, you need to configure the mail service to use this function. The configuration is located in `conf/app.conf`


```bash

#mail service configuration
enable_mail=true
smtp_user_name=admin@iminho.me
smtp_host=smtp.ym.163.com
smtp_password=1q2w3e__ABC
smtp_port=25
form_user_name=admin@iminho.me
mail_expired=30
```


# Use Docker deployment

Refer to the built-in Dockerfile project files to compile the mirror.

The following environment variables need to be provided when starting the mirror:

```ini
MYSQL_PORT_3306_TCP_ADDR    MySQL Address
MYSQL_PORT_3306_TCP_PORT    MySQL Port
MYSQL_INSTANCE_NAME         MySQL Database name
MYSQL_USERNAME              MySQL Username
MYSQL_PASSWORD              MySQL Password
HTTP_PORT                   Listen Port
```

For Example

```bash
docker run -p 8181:8181 -e MYSQL_PORT_3306_TCP_ADDR=127.0.0.1 -e MYSQL_PORT_3306_TCP_PORT=3306 -e MYSQL_INSTANCE_NAME=mindoc_db -e MYSQL_USERNAME=root -e MYSQL_PASSWORD=123456 -e httpport=8181 -d daocloud.io/lifei6671/mindoc:latest
```

# Technology used

- beego 1.8.1
- mysql 5.6
- editor.md
- bootstrap 3.2
- jquery 
- layer 
- webuploader 
- Nprogress 
- jstree 
- font awesome 
- cropper 
- highlight 
- to-markdown 
- wangEditor


# Main function

- Project management, you can edit the project, add members, etc.
- Document management, adding and deleting documents, etc.
- Comment management, you can manage document comments and comments posted by yourself.
- User management, adding and disabling users, changing personal information, etc.
- User authority management, change user roles.
- Project encryption, you can set the public status of the project, and private projects need to be accessed through Token.
- Site configuration, anonymous access, verification code, etc.

# Contributing

We welcome you to report issue or pull request on the GitHub.

If you are not familiar with GitHub's Fork and Pull development model, you can read the GitHub documentation (https://help.github.com/articles/using-pull-requests) for more information.

