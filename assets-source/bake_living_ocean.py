"""Reproducible v0.3 Lighthouse Keeper art bake (offline only).

Input: original approved-ocean-study.png already in Garage GitHub.
Output: clean animated sea, fixed rock/shore, water domain, spray atlas.
Requires Pillow, NumPy, OpenCV on authoring/CI only; no runtime dependency.
Never classify rocks by pixel brightness at runtime.
"""
from pathlib import Path
import cv2
import numpy as np
from PIL import Image, ImageDraw, ImageFilter

ROOT=Path(__file__).resolve().parent.parent
OUT=ROOT/'cmd'/'lighthousekeeper'/'assets'
OUT.mkdir(parents=True,exist_ok=True)
ORIGINAL=ROOT/'assets-source'/'approved-ocean-study.png'
W,H=960,720
original=np.asarray(Image.open(ORIGINAL).convert('RGB').crop((395,0,1640,928)).resize((W,H),Image.Resampling.LANCZOS))
# Remove the source painting's decorative boat BEFORE sampling sky/water donors.
# The gameplay boat is a separate moving object; source donors must be boat-free.
boat_x,boat_y,boat_rx,boat_ry=635,407,82,73
boat_donor=original[boat_y-boat_ry:boat_y+boat_ry+1,405-boat_rx:405+boat_rx+1].copy()
boat_mask=np.zeros((2*boat_ry+1,2*boat_rx+1),np.uint8)
cv2.ellipse(boat_mask,(boat_rx,boat_ry),(boat_rx-5,boat_ry-5),0,0,360,255,-1)
original=cv2.cvtColor(cv2.seamlessClone(
    cv2.cvtColor(boat_donor,cv2.COLOR_RGB2BGR),
    cv2.cvtColor(original,cv2.COLOR_RGB2BGR),boat_mask,
    (boat_x,boat_y),cv2.NORMAL_CLONE),cv2.COLOR_BGR2RGB)
Y,X=np.mgrid[:H,:W].astype('float32')

def poly_mask(points,feather=0):
    m=Image.new('L',(W,H),0)
    ImageDraw.Draw(m).polygon(points,fill=255)
    if feather:m=m.filter(ImageFilter.GaussianBlur(feather))
    return np.asarray(m).copy()

def export_2x(img,name,quality=None):
    image=Image.fromarray(img) if isinstance(img,np.ndarray) else img
    large=image.resize((image.width*2,image.height*2),Image.Resampling.LANCZOS)
    if quality:
        large.save(OUT/name,quality=quality,optimize=True,subsampling=0)
    else:large.save(OUT/name,optimize=True)

# 1. The baked lighthouse beam must not exist in the game's static plate.
# Author-selected sky/water donors replace its entire original cone offline;
# the game separately draws the only real, moving lighthouse beam.
beam_poly=[(-10,123),(90,153),(250,199),(425,269),(620,345),
           (735,425),(735,493),(580,463),(405,419),(200,341),(-10,268)]
a=cv2.GaussianBlur(poly_mask(beam_poly),(0,0),28).astype('float32')/255
sx=np.clip(570-.37*X,0,W-1).astype('float32')
sy=np.clip(Y-112,0,H-1).astype('float32')
wx=np.clip(540+.28*X,0,W-1).astype('float32')
wy=np.clip(Y+125,0,H-1).astype('float32')
sky=cv2.remap(original,sx,sy,cv2.INTER_LINEAR)
sea=cv2.remap(original,wx,wy,cv2.INTER_LINEAR)
f=np.clip((Y-255)/125,0,1)**2
donor=sky*(1-f[:,:,None])+sea*f[:,:,None]
clean=np.uint8(np.clip(original.astype('float32')*(1-a[:,:,None])+donor*a[:,:,None],0,255))

# 2. The accepted fixed hero rocks are manually authored silhouettes. First
# replace their pixels and immovable spray with open water of similar depth.
rock_polys=[
 [(219,570),(228,557),(240,551),(257,551),(275,552),(285,556),
  (295,563),(301,574),(306,585),(303,597),(296,605),(279,608),
  (263,604),(252,601),(239,596),(227,590),(219,581)],
 [(291,589),(297,579),(308,575),(321,578),(333,586),(344,597),
  (350,605),(350,612),(338,613),(323,615),(309,610),(299,604)],
]
hero=Image.new('L',(W,H),0)
for p in rock_polys:ImageDraw.Draw(hero).polygon(p,fill=255)
hero_m=np.asarray(hero).copy()
# Replace the whole static plume and its underlying rocks from one water donor.
cx,cy,rx,ry=265,565,148,116
sx,sy=525,557
donor=clean[sy-ry:sy+ry+1,sx-rx:sx+rx+1].copy()
mask=np.zeros((2*ry+1,2*rx+1),np.uint8)
cv2.ellipse(mask,(rx,ry),(rx-8,ry-8),0,0,360,255,-1)
clean=cv2.cvtColor(cv2.seamlessClone(
    cv2.cvtColor(donor,cv2.COLOR_RGB2BGR),
    cv2.cvtColor(clean,cv2.COLOR_RGB2BGR),mask,(cx,cy),cv2.NORMAL_CLONE),cv2.COLOR_BGR2RGB)

# 3. Also separate the broad shoreline/left cliff and bottom-centre rocks.
cliff=[(0,381),(36,382),(58,400),(78,433),(88,457),(117,454),
 (139,466),(149,490),(164,521),(160,544),(174,565),(165,586),
 (178,611),(222,639),(245,654),(272,664),(304,672),(332,674),
 (360,655),(409,641),(442,645),(464,655),(488,679),(524,720),(0,720)]
cliff_m=poly_mask(cliff,0)
plume=[(152,510),(178,516),(204,534),(220,569),(223,603),
       (261,638),(244,657),(193,634),(164,615),(145,586),(136,556),(147,531)]
cut=cv2.dilate(np.uint8(hero_m>0),np.ones((11,11),np.uint8))>0
cliff_m[cut]=0
psoft=poly_mask(plume,3)
cliff_m=np.uint8(np.clip(cliff_m.astype('int16')-psoft.astype('int16'),0,255))
cliff_m=np.asarray(Image.fromarray(cliff_m).filter(ImageFilter.GaussianBlur(.75)))
shore_rgba=np.dstack([original,cliff_m]).astype(np.uint8)
export_2x(Image.fromarray(shore_rgba,'RGBA').crop((0,380,545,720)),'shore-fixed.png')

# The rock art is never passed through the moving-water shader, even once.
# Only small isolated fixed white flecks are retouched in this OFFLINE bake.
rock_rgb=original.copy();mean=rock_rgb.mean(axis=2).astype(np.float32)
contrast=mean-cv2.GaussianBlur(mean,(0,0),2.5)
sat=rock_rgb.max(axis=2).astype('int16')-rock_rgb.min(axis=2).astype('int16')
foam=(mean>88)&(contrast>8)&(sat<75)&((rock_rgb[:,:,2].astype('int16')-rock_rgb[:,:,0])>8)&(Y>570)&(hero_m>0)
if foam.sum()>20:
    corrected=cv2.cvtColor(cv2.inpaint(cv2.cvtColor(original,cv2.COLOR_RGB2BGR),np.uint8(foam)*255,2.3,cv2.INPAINT_TELEA),cv2.COLOR_BGR2RGB)
    rock_rgb[foam]=corrected[foam]
alpha=np.asarray(Image.fromarray(hero_m).filter(ImageFilter.GaussianBlur(.75)))
rock_rgba=np.dstack([rock_rgb,alpha])
export_2x(Image.fromarray(rock_rgba,'RGBA').crop((157,520,365,641)),'rocks-fixed.png')

bottom=[(464,720),(495,693),(523,671),(555,663),(592,669),
        (624,661),(670,665),(707,680),(748,720)]
bottom_m=poly_mask(bottom,1.0)
bottom_rgba=np.dstack([original,bottom_m])
export_2x(Image.fromarray(bottom_rgba,'RGBA').crop((450,645,775,720)),'rocks-bottom.png')
# Replace the static bottom stone inside the moving ocean source.
m=np.asarray(Image.fromarray(poly_mask(bottom)).filter(ImageFilter.GaussianBlur(3))).astype('float32')/255
water_donor=cv2.remap(clean,np.clip(X-240,0,W-1).astype('float32'),np.clip(Y-63,0,H-1).astype('float32'),cv2.INTER_LINEAR)
clean=np.uint8(np.clip(clean*(1-m[:,:,None])+water_donor*m[:,:,None],0,255))
export_2x(Image.fromarray(clean,'RGB'),'seascape.jpg',quality=86)

# One explicit geometric water-domain matte; no colour/brightness rock mask.
geo=Image.new('L',(W,H),0);draw=ImageDraw.Draw(geo)
draw.polygon([(0,346),(90,354),(135,381),(165,415),(180,440),(209,446),
 (240,420),(330,399),(410,382),(560,350),(705,332),(820,334),
 (955,321),(960,720),(0,720)],fill=255)
draw.polygon(cliff,fill=0)
export_2x(geo.filter(ImageFilter.GaussianBlur(2)),'water-domain.png')

# 4. Artistically placed whitewater poses from the *same original painting*.
# Offline sprite extraction, not runtime pixel drawing. The ten poses are a
# deliberately small atlas, two distinct five-pose impact animations.
# Whitewater regions are manually selected, excluding hard rock silhouettes.
raw=Image.fromarray(original,'RGB')
regions=[(132,491,220,570),(313,526,424,616)]
atlas=Image.new('RGBA',(640,256),(0,0,0,0))
for variation,box in enumerate(regions):
    crop=raw.crop(box).resize((112,112),Image.Resampling.LANCZOS).convert('RGB')
    pix=np.asarray(crop,dtype=np.float32)
    bright=(pix[:,:,0]+pix[:,:,1]+pix[:,:,2])/3
    # Purpose-made *offline spray* cutout, never used to classify game rocks.
    r,g,b=pix[:,:,0],pix[:,:,1],pix[:,:,2]
    whiteness=np.clip((bright-111)/85,0,1)**1.6 * np.clip((b-r+24)/55,0,1)
    yy,xx=np.mgrid[:112,:112];shape=np.clip(1-((xx-56)/65)**2-((yy-56)/70)**2,0,1)**.7
    raw_a=np.uint8(np.clip(whiteness*shape*255,0,255))
    # A 5-step growth-and-collapse cycle with different artist-directed poses.
    for frame in range(5):
        s=[.55,.84,1.10,.93,.62][frame]
        rise=[-1,-5,-15,-9,1][frame]
        angle=[-6,-2,2,5,8][frame]*(-1 if variation else 1)
        material=Image.fromarray(np.dstack([pix.astype(np.uint8),raw_a]),'RGBA')
        material=material.rotate(angle,Image.Resampling.BICUBIC,expand=False)
        material=material.resize((int(108*s),int(110*s)),Image.Resampling.BICUBIC)
        tile=Image.new('RGBA',(128,128),(0,0,0,0))
        tile.alpha_composite(material,((128-material.width)//2,(128-material.height)//2+rise))
        atlas.alpha_composite(tile,((variation*5+frame)%5*128,(variation*5+frame)//5*128))
atlas.save(OUT/'impact-atlas.png',optimize=True)

print('Living-sea build complete:')
for name in ['seascape.jpg','water-domain.png','shore-fixed.png','rocks-fixed.png','rocks-bottom.png','impact-atlas.png']:
    p=OUT/name;print(name,p.stat().st_size)
