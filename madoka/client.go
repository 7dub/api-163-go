package madoka

import (
	"os"
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
	"net/url"
	"strconv"
	"crypto/md5"
	"encoding/base64"
)

/** 
 * 执行搜索
 * params: 	关键词 类型 页码 数量
 * return:	字符串形式的请求结果
 */
func Search(words string, stype string, page int, limit int) string {
	// 创建客户端
	client := &http.Client{}
	// 格式化参数
	_o, _l := formatParams(page, limit)
	// 设置body
	form := url.Values{}
	form.Set("s", words)
	form.Set("type", stype)
	form.Set("limit", _l)
	form.Set("offset", _o)
	body := strings.NewReader(form.Encode())
	// 创建请求
	request, _ := http.NewRequest("POST", "http://music.163.com/api/search/get/", body)
	//设置头部
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Cookie", "appver=2.0.2")
	request.Header.Set("Referer", "http://music.163.com")
	request.Header.Set("Content-Length", (string)(body.Len()))
	// 发起请求
	response, reqErr := client.Do(request)
	// 错误处理
	if reqErr!= nil {
		fmt.Println("Fatal error ", reqErr.Error())
		os.Exit(0)
	}
	defer response.Body.Close()
	resBody, _ := ioutil.ReadAll(response.Body)
	return string(resBody)
}

/**
 * 根据传入id返回生成的mp3地址
 */
func Download(id string) string {
	return "http://music.163.com/api/song/enhance/download/url?br=320000&id=" + id
}

/**
* 传入 搜索类型 页码 数量
* 返回 搜索类型 偏移 数量
*/
func formatParams(page int, limit int) (string, string) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 0
	}
	return strconv.Itoa((page - 1) * limit), strconv.Itoa(limit)
}

func encryptSongId(_id string) string {
	// 5249677
	rs := [50]byte{}
	magic := []byte("3go8&$8*3*3h0k(2)2")
	magic_len := len(magic)
	id := []byte(_id)
	id_len := len(id)
	for i := 0; i < id_len; i++ {
		rs[i] = id[i] ^ magic[i % magic_len]
	}
	m := md5.New()
	m.Write(id)
	md5Sum := m.Sum(nil)
	fmt.Println(md5Sum)
	b := base64.StdEncoding.EncodeToString(md5Sum)
	b = strings.Replace(b, "/", "_", -1)
	b = strings.Replace(b, "+", "-", -1)
    return b
}
