This can be solved without writing any code.

I noticed that the product ids start with two, so we change the url to see what's product/1 - it's the flag!

http://staywoke.hax1.allesctf.net/products/1

We need 1337€ but we don't have anything yet.

We notice a hidden form field to a `payment-api` host. Here I lost a few hours thinking there has to be some way to send the payment request to another API that would return valid and that would approve my purchase.

While I was playing with a script to automate this, I noticed that the payment API returns debug information when it errors. So we can replace the endpoint with the help endpoint `http://payment-api:9090/help?`, we get this

```
Error from Payment API: {"endpoints":[{"method":"GET","path":"/wallets/:id/balance","description":"check wallet balance"},{"method":"GET","path":"/wallets","description":"list all wallets"},{"method":"GET","path":"/help","description":"this help message"}]}
```

There's some automation in the go script, but this can be solved just by clicking the UI.

We list the wallets by changing the API URL in the form using the browser's devtools, we find the ID of an account with 1335€ in it. So we're short of two euros for getting the flag.

The website has a redeem code functionality and we can see the code `I<3CORONA` in the news ticker on the top. If we play with this a bit we see that it can't be applied to the flag, but it can be applied to other items. And if we delete the other items from the cart, we keep the discount!

So we make a two euro discount, we add the flag, remove all other items then win:

```
CSCG{c00l_k1ds_st4y_@_/home/}
```

![flag.png]