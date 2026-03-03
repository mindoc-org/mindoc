package requests

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"fmt"
)

//下载远程文件并保存到指定位置
func DownloadAndSaveFile(remoteUrl, dstFile string) (error) {
	client := &http.Client{}
	uri, err := url.Parse(remoteUrl)
	if err != nil {
		return err
	}
	// Create the file
	out, err := os.Create(dstFile)
	if err != nil  {
		return err
	}
	defer out.Close()

	request, err := http.NewRequest("GET", uri.String(), nil)
	request.Header.Add("Connection", "close")
	request.Header.Add("Host", uri.Host)
	request.Header.Add("Referer", uri.String())
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()


	if resp.StatusCode == http.StatusOK {
		_, err = io.Copy(out, resp.Body)
	}else{
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	return nil
}
