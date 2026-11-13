"""Offline moving sprite and habitat regression."""
from pathlib import Path
import numpy as np
from PIL import Image
A=Path(__file__).resolve().parent.parent/'cmd/lighthousekeeper/assets'
b=np.asarray(Image.open(A/'boat.png').convert('RGBA'))
assert b.shape==(256,256,4)
assert np.count_nonzero(b[225:,:,3]>20)==0,'fixed wake below hull'
assert np.count_nonzero(b[:,0:20,3]>20)==0
assert np.count_nonzero(b[:,238:,3]>20)==0
assert np.count_nonzero(b[:,:,3]>100)>2500
atlas=np.asarray(Image.open(A/'sightings-live.png').convert('RGBA'))
assert atlas.shape==(240,1440,4)
for i in range(6):
 m=atlas[:,i*240:(i+1)*240,3]
 assert np.count_nonzero(m>20)>350,(i,'subject lost')
 assert np.count_nonzero(m[:20,:]>20)==0
 assert np.count_nonzero(m[220:,:]>20)==0
 assert np.count_nonzero(m[:,:20]>20)==0
 assert np.count_nonzero(m[:,220:]>20)==0
water=np.asarray(Image.open(A/'water-domain.png').convert('L'))
for l,r,t,b in [(415,755,205,450),(540,780,215,425)]:
 for x in range(l,r+1,20):
  for y in range(t,b+1,20):
   assert water[int((338+.42*y)*2),x*2]>235,(x,y)
print('PASS: clean boat and six sighting sprites, sea-only habitat')
