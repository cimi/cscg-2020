from keras.models import load_model
from keras.preprocessing import image
from sklearn.preprocessing import LabelEncoder
from processor import split_image
import cv2
import numpy as np
import os

def get_encodings():
  letters = []
  for filename in os.listdir('training/letters/'):
    if '-' not in filename:
      continue
    letter = filename.split('-')[0]
    letters.append(letter)

  encoder = LabelEncoder()
  labels = encoder.fit_transform(letters)

  encodings = {}
  for idx, l in enumerate(letters):
    encodings[labels[idx]] = l
  return encodings

def decode_results(results, encodings):
  guess = ""
  for result in results:
    max_idx, max_val = -1, -1
    for idx, prob in enumerate(result):
      if prob > max_val:
        max_idx, max_val = idx, prob
    guess += encodings[max_idx]
  return guess

def validate(model_file):
  validation_dir = "captchas/"
  encodings = get_encodings()
  model = load_model(model_file)
  total, success = 0, 0
  for filename in os.listdir(validation_dir):
    total += 1
    letters = split_image(validation_dir + filename)
    captcha = filename.split(".")[0]
    tmpfile = "tmp-letter.png"
    images = []
    for l in letters:
      cv2.imwrite(tmpfile, l)
      img = image.load_img(tmpfile, target_size=[30, 30, 1], color_mode='grayscale')
      img = image.img_to_array(img)
      img = img/255
      images.append(img)

    results = model.predict(np.array(images), verbose=0)
    guess = decode_results(results, encodings)

    if guess == captcha:
      success += 1
  
  success_percentage = (success / total) * 100
  if success_percentage < 99:
    print(f"🚫 {model_file} success too low: {success_percentage:.2f}%")
  else:
    print(f"✅ {model_file} is a great success: {success_percentage:.2f}%")
