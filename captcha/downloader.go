package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

const maxDownloads = 10 * 1000
const target = "http://hax1.allesctf.net:9200/captcha/0"

var hc = &http.Client{}

var reImg = regexp.MustCompile(`"data:image/png;base64,([^"]+)"`)
var reSol = regexp.MustCompile(`The solution would have been <b>([^<]+)</b>`)

func main() {
	for i := 0; i < maxDownloads; i++ {
		// get cookie and image
		base64Img, cookie := get()
		// submit some garbage
		// print(cookie)
		resp := post(cookie)
		bytes := dump(resp)

		match := reSol.FindAllStringSubmatch(string(bytes), -1)
		solution := match[0][1]
		imgBytes, err := base64.StdEncoding.DecodeString(base64Img)
		check(err)

		ioutil.WriteFile("training/captchas/"+solution+".png", imgBytes, 0644)
	}
}

func get() (string, string) {
	resp, err := http.Get(target)
	check(err)
	cookie := ""
	for _, c := range resp.Cookies() {
		if c.Name == "session" {
			cookie = c.Value
			break
		}
	}
	if cookie == "" {
		panic("Could not extract cookie")
	}
	bytes := dump(resp)
	match := reImg.FindAllStringSubmatch(string(bytes), -1)
	base64Img := match[0][1]
	return base64Img, cookie
}

func post(cookie string) *http.Response {
	form := url.Values{}
	form.Add("0", "g4rb4g3")
	req, err := http.NewRequest("POST", target, strings.NewReader(form.Encode()))
	check(err)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "session="+cookie)
	resp, err := hc.Do(req)
	check(err)
	return resp
}

func dump(resp *http.Response) []byte {
	dump, err := httputil.DumpResponse(resp, true)
	check(err)
	return dump
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
