## wkhtmltopdf 安装

导出 pdf 格式文档使用的是 wkhtmltopdf 工具，工具下载地址为：[https://wkhtmltopdf.org/downloads.html](https://wkhtmltopdf.org/downloads.html)。

### Windows 下配置

下载 Windows 版本，安装即可。

### Linux 下配置

请下载和你服务器对应的版本，Linux版本依赖一下库：

`zlib fontconfig freetype X11 libs (libX11, libXext, libXrender)`

请自行搜索安装以上依赖包，下面的命令是安装 libXrender 和 libXext。

```bash
apt-get install -y libxrender-dev
apt-get install -y libxext-dev
```

## 配置

请将 wkhtmltopdf 可执行文件所在目录配置到 MinDoc 根目录下 conf/app.conf 的 wkhtmltopdf 节点。

配置完后请重启程序。
