import os
import re
import cv2
import secrets
import hashlib

# we don't need bounding rectangles because the images are all the same height
# we may want to do this on characters later though if they don't fill vertical space?
def load_image_opencv(filename):
  im = cv2.imread(filename, cv2.IMREAD_GRAYSCALE)
  print(len(im))
  # scan image to get 
  ret, thresh = cv2.threshold(im, 127, 255, cv2.THRESH_BINARY_INV)
  contours, hierarchy = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

  contours_poly = [None]*len(contours)
  boundRect = [None]*len(contours)
  
  for i, c in enumerate(contours):
    epsilon = 0.1*cv2.arcLength(c,True)
    contours_poly[i] = cv2.approxPolyDP(c, epsilon, True)
    boundRect[i] = cv2.boundingRect(contours_poly[i])

  points = []
  for i in range(len(contours)):
    x, y, width, height = boundRect[i]
    points.append((x + width // 2, y + height // 2, width // height))

  return im, points


# simple image on all-blank vertical lines - all characters are contiguous (I hope)
def split_image(filename):
  im = cv2.imread(filename, cv2.IMREAD_GRAYSCALE)
  height = len(im)
  width = len(im[0])
  result = []
  curr = None

  # avgcol = 0
  # for x in range(0, width):
  #   col = 0
  #   for y in range(0, height):
  #     col += im[y][x]
  #   avgcol += col
  # avgcol = avgcol / width
  sums = []
  color_threshold = 3*255
  width_threshold = 21
  for x in range(0, width):
    # see if the column is blank in two passes
    # first we get the average in the image
    
    sum = 0
    for y in range(0, height):
      sum += im[y][x]
    sums.append(abs(255 * height - sum) < color_threshold)
    blank = abs(255 * height - sum) < color_threshold

    if blank and curr:
      # if the length is more than we expect break in half
      if (x+1 - curr) >= width_threshold:
        half = curr + ((x+1 - curr) // 2)
        crop_img = im[0:height,curr:half+1]
        result.append(crop_img)
        curr = half+1

      crop_img = im[0:height,curr:x+1]
      result.append(crop_img)
      curr = None
    elif not blank and not curr:
      curr = x
    else:
      continue
  if curr != None:
    crop_img = im[0:height,curr:width]
    result.append(crop_img)
  return result


def process_captchas():
  errors, total = 0, 0
  debug = None
  for filename in os.listdir("training/captchas/"):
    total += 1
    if debug != None and filename != debug+".png":
      continue
    letters = split_image("training/captchas/"+filename)
    captcha = filename.split(".")[0]
    
    if len(letters) != len(captcha):
      print("Number of letters different from expected!", captcha, len(letters), len(captcha))
      errors += 1
      continue

    # write to a temp file which we rename in step two
    # if debug == None:
    #   continue
    tmpfile = "tmp-" + secrets.token_hex(nbytes=6) + ".png"
    for idx, img in enumerate(letters):
      cv2.imwrite(tmpfile, img)

      # append the md5 to the file name to avoid duplicates
      hash_md5 = hashlib.md5()
      with open(tmpfile, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
          hash_md5.update(chunk)
      try:
        dstfile = "training/letters/" + "-".join([captcha[idx], hash_md5.hexdigest()[:5]]) + ".png"
        os.rename(tmpfile, dstfile)
      except:
        print("Failed to write letter file", captcha, idx, hash_md5.hexdigest())
    
  print("Errors: ", errors, total)
