
I attempted to solve another captcha cracking challenge 2019 X-MAS CTF and I failed. The challenge was slightly different in that there were only a few symbols but they were custom squiggles, not digits or letters. Also, they were spread out on a larger 2d rectangle compared to the one line script we have here.

Then I tried to use the AWS machine learning service, train a model in the cloud and then call an API to solve the challenge captcha interactively. Unfortunately, after spending a couple of hours labelling images I realised that the AWS service doesn't offer an interactive API, it only analysed images in batch and the expected response time was a few minutes, so I gave up.

I followed up after the contest ended and I found a very nice write-up of a solution that trained a neural net using python, keras and tensor flow. I thought it was really cool so I implemented something similar here. https://ctftime.org/writeup/17585

![flag.jpg]

The algorithm doesn't work all the time for the final large batch of catchas, but after retraining it a few times I managed to get the accuracy around 98% which is enough to make it pass some runs.

TODO: cleanup the code and explain in more detail
