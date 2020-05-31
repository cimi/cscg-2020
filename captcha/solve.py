from keras.models import load_model
from keras.preprocessing import image
from sklearn.preprocessing import LabelEncoder
from processor import split_image
from validate import get_encodings, decode_results
import cv2
import numpy as np
import os
from bs4 import BeautifulSoup
import requests
import base64

encodings = get_encodings()
model = load_model('model-homogenous-500.hdf5')
s = requests.Session()

def guess_captcha(filename):
  letters = []
  tmpfile = "letter.png"
  for letter in split_image(filename):
    cv2.imwrite(tmpfile, letter)
    tmp_img = image.load_img(tmpfile, target_size=[30, 30, 1], color_mode='grayscale')
    tmp_img = image.img_to_array(tmp_img)
    tmp_img = tmp_img/255
    letters.append(tmp_img)

  results = model.predict(np.array(letters), verbose=0)
  return decode_results(results, encodings)

def get_images_b64(text):
  html = BeautifulSoup(text, features="html.parser")
  img_tags = html.body.find_all('img')
  print("Need to guess ", len(img_tags))
  results = []
  for img_tag in img_tags:
    img_b64 = img_tag['src'][len("data:image/png;base64,"):]
    results.append(img_b64)
  return results

def solve():
  round = 0
  while round <= 4:
    target = "http://hax1.allesctf.net:9200/captcha/" + str(round)
    resp = s.get(target)
    if round == 4:
      print("🎉 WIN!")
      print(resp.text)
      break
    imgs = get_images_b64(resp.text)
    tmpfile = "captcha.png"
    formdata = {}
    for idx, img_b64 in enumerate(imgs):
      with open(tmpfile, "wb") as captcha_img:
        captcha_img.write(base64.b64decode(img_b64))
      formdata[str(idx)] = guess_captcha(tmpfile)
    
    resp = s.post(target, formdata)
    if "Human detected" in resp.text:
      print("🚫 Failed guess", formdata)
      break
    else:
      print("✅ Correct guess", formdata)
      round += 1
