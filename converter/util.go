//Author:TruthHun
//Email:TruthHun@QQ.COM
//Date:2018-01-21
package converter

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

//media-type
var MediaType = map[string]string{
	".jpeg":  "image/jpeg",
	".png":   "image/png",
	".jpg":   "image/jpeg",
	".gif":   "image/gif",
	".ico":   "image/x-icon",
	".bmp":   "image/bmp",
	".html":  "application/xhtml+xml",
	".xhtml": "application/xhtml+xml",
	".htm":   "application/xhtml+xml",
	".otf":   "application/x-font-opentype",
	".ttf":   "application/x-font-ttf",
	".js":    "application/x-javascript",
	".ncx":   "x-dtbncx+xml",
	".txt":   "text/plain",
	".xml":   "text/xml",
	".css":   "text/css",
}

//根据文件扩展名，获取media-type
func GetMediaType(ext string) string {
	if mt, ok := MediaType[strings.ToLower(ext)]; ok {
		return mt
	}
	return ""
}

//解析配置文件
func parseConfig(configFile string) (cfg Config, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(configFile); err == nil {
		err = json.Unmarshal(b, &cfg)
	}
	return
}

