"""Offline geometry and art invariants; NOT a hardware FPS test."""
from pathlib import Path
import numpy as np
import cv2
from PIL import Image

ROOT=Path(__file__).resolve().parent.parent
ART=ROOT/'cmd/lighthousekeeper/assets'
W,H=960,720
sea=np.asarray(Image.open(ART/'seascape.jpg').convert('RGB').resize((W,H)))
mask=np.asarray(Image.open(ART/'water-domain.png').convert('L').resize((W,H)),dtype='float32')/255
rock=Image.open(ART/'rocks-fixed.png').convert('RGBA').resize((208,121))
shore=Image.open(ART/'shore-fixed.png').convert('RGBA').resize((545,340))
bottom=Image.open(ART/'rocks-bottom.png').convert('RGBA').resize((325,75))
atlas=Image.open(ART/'impact-atlas.png').convert('RGBA')
assert atlas.size==(640,256)
assert rock.getchannel('A').getextrema()==(0,255)
alpha=np.asarray(rock.getchannel('A'))
core=(alpha==255)
assert core.sum()>2500
Y,X=np.mgrid[:H,:W].astype('float32')

def moving(t):
 depth=np.clip((Y-570/2)/(840/2),0,1)
 depth=depth*depth*(3-2*depth)
 a=np.sin(Y*.055*2-t*1.8+X*.0075*2)
 b=np.sin(Y*.023*2+X*.009*2-t*.72)
 c=np.sin(X*.048*2+Y*.012*2+t*1.4)
 amp=(.32+depth*1.28)*mask
 xx=np.clip(X+(3.2*a+2.1*b+1.3*c)*amp,1,W-2).astype('float32')
 yy=np.clip(Y+(1.1*a+1.5*b)*amp,1,H-2).astype('float32')
 water=cv2.remap(sea,xx,yy,cv2.INTER_LINEAR)
 valid=cv2.remap(mask,xx,yy,cv2.INTER_LINEAR)
 f=(mask*valid)[:,:,None]
 return Image.fromarray(np.uint8(np.clip(sea*(1-f)+water*f,0,255)),'RGB').convert('RGBA')

def fixed(t):
 im=moving(t)
 im.alpha_composite(shore,(0,380))
 im.alpha_composite(rock,(157,520))
 im.alpha_composite(bottom,(450,645))
 return np.asarray(im.convert('RGB'))

im0=fixed(.15);im1=fixed(1.25)
roi0=im0[520:641,157:365];roi1=im1[520:641,157:365]
assert np.array_equal(roi0[core],roi1[core]),'rock pixels changed despite fixed foreground layer'
water_roi=np.zeros((H,W),bool);water_roi[430:590,430:790]=True
movement=np.abs(im0.astype('int16')-im1.astype('int16')).max(axis=2)>5
assert np.count_nonzero(movement&water_roi)>2500,'sea not visibly alive'
# The 10 frames are genuinely distinct, not one immovable spray repeated.
frames=[]
for j in range(10):frames.append(np.asarray(atlas.crop(((j%5)*128,(j//5)*128,(j%5+1)*128,(j//5+1)*128))))
assert np.count_nonzero(frames[0]!=frames[2])>1000
assert np.count_nonzero(frames[5]!=frames[7])>1000
# The art bake must be small; no giant sprite sheet, no videos in runtime.
extra=['water-domain.png','shore-fixed.png','rocks-fixed.png','rocks-bottom.png','impact-atlas.png']
size=sum((ART/p).stat().st_size for p in extra)
assert size<1600000,size
print('PASS: fixed rock',int(core.sum()),'pixels; changing sea',int(np.count_nonzero(movement&water_roi)),'pixels; additional art',size,'bytes')
