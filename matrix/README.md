We get a wav file with a recording from the Matrix repeated 60 times.

Every fifth repetition is of a different length. I spent a few hours splitting the file and comparing the segments trying to find patterns, but I couldn't validate anything that looked interesting.

If we visualise the spectrogram of the file we see some text with the password!

Using this password we can extract data from the wav file and we get a file called `redpill.jpg`.

This looks similar to the Intro to stegano 2 challenge - we use Google image search to find the original photo. We then overlay them in GIMP and notice that some of the small lights are coloured red blue. If we translate this to binary we get this:

```python
bs = ["01101110","00100001","01000011","00110011","01011111","01010000","01010111","00111111"]

```

Then if we do a binwalk of the file we see there's an embedded archive. We extract it, then unzip it with the password we got from the image. Inside the file is the flag.

```
```
