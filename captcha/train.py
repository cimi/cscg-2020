import keras
from keras.models import Sequential
from keras.layers import Dense, Dropout, Flatten
from keras.layers import Conv2D, MaxPooling2D
from keras.utils import to_categorical
from keras.preprocessing import image
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import LabelEncoder
from keras.utils import to_categorical
from tqdm import tqdm
import os

def get_train_files():
    train_files = []
    for filename in os.listdir('training/letters/'):
        if '-' not in filename:
            continue
        label = filename.split('-')[0]
        train_files.append((filename, label))

    return train_files


def train():
    train_files = get_train_files()
    train_image = []
    for i in range(len(train_files)):
        img = image.load_img('training/letters/' + train_files[i][0], target_size=[30, 30, 1], color_mode='grayscale')
        img = image.img_to_array(img)
        img = img/255
        train_image.append(img)

    X = np.array(train_image)
    encoder = LabelEncoder()
    transformed_labels = encoder.fit_transform([t[1] for t in train_files])
    y=pd.Series(transformed_labels)
    y = to_categorical(y)
    
    X_train, X_test, y_train, y_test = train_test_split(X, y, random_state=42, test_size=0.2)

    model = Sequential()
    model.add(Conv2D(32, kernel_size=(3, 3),activation='relu',input_shape=(30,30,1)))
    model.add(Conv2D(64, (3, 3), activation='relu'))
    model.add(MaxPooling2D(pool_size=(2, 2)))
    model.add(Dropout(0.25))
    model.add(Flatten())
    model.add(Dense(128, activation='relu'))
    model.add(Dropout(0.5))
    model.add(Dense(35, activation='softmax'))

    model.compile(loss='categorical_crossentropy',optimizer='Adam',metrics=['accuracy'])

    model.fit(X_train, y_train, epochs=10, validation_data=(X_test, y_test))
    model.save('model-homogenous-500.hdf5')

def clean_letters():
    max_count = 500
    counts = {}
    for filename in os.listdir('training/letters/'):
        if '-' not in filename:
            continue
        letter = filename.split('-')[0]
        try:
            counts[letter] += 1
        except:
            counts[letter] = 1
        if counts[letter] > max_count:
            os.remove('training/letters/' + filename)
            print("Deleted", filename)
    print(counts)
