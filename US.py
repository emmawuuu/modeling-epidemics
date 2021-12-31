import numpy as np
from PIL import Image
from scipy import signal

image_name = "US.png"
image = Image.open(image_name)
npimg = np.asarray(image)
npimg = np.where(npimg == 255, 0, 1)[:, :, 0]
npimg = np.pad(npimg, ((300, 300), (0, 0)))
convarr = np.ones((12, 12))
def strideConv(arr,arr2,s):
    cc= signal.convolve(arr,arr2[::-1,::-1],mode='valid')
    idx=(np.arange(0,cc.shape[1],s), np.arange(0,cc.shape[0],s))
    xidx,yidx=np.meshgrid(*idx)
    return cc[yidx,xidx]
npimg = strideConv(npimg, convarr, 12)
npimg = np.where(npimg < 12 * 6, 0, 1)

# 1200 
# 1 -> ALL Squares 
# 2 -> 0 for water 1 for land 
# 3 -> 1 for all 
f = open("US.txt","w")
f.write("100\n")
f.write(" ".join(["{} {}".format(x, y) for x in range(100) for y in range(100)]) + "\n")
f.write(" ".join(npimg.astype(str).flatten()) + "\n")
f.write(" ".join(["1" for x in range(100 * 100)]))
f.close()
