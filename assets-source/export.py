"""Offline source-art export. No runtime dependency; reproducible original layers.
Uses user-approved original ChatGPT-generated concept paintings as source textures.
Run from the repository root with: python generate_assets.py
"""
from PIL import Image,ImageDraw,ImageFilter,ImageEnhance
from pathlib import Path
import math, random

ROOT=Path(__file__).resolve().parent.parent
A=ROOT/'cmd/lighthousekeeper/assets'; A.mkdir(parents=True,exist_ok=True)
R=3; random.seed(5731)
CONCEPTS=ROOT/'assets-source'
W,H=1920,1440

# Fully original repainting from approved concept colour palette and water texture.
# Sky gets its clouds and moon from an original concept, but no hero/gameplay objects.
sky=Image.new('RGBA',(W,800)); sd=ImageDraw.Draw(sky)
for y in range(800):
 t=y/800
 sd.line((0,y,W,y),fill=(int(8+12*t),int(19+23*t),int(44+40*t),255))
# Impressionistic clusters of softened blue night clouds, seeded, separately behind gameplay objects.
cloud=Image.new('RGBA',(W,800)); cd=ImageDraw.Draw(cloud,'RGBA')
for i in range(280):
 x=random.randint(-120,2040);y=random.randint(40,680)
 rx=random.randint(32,160); ry=random.randint(10,45)
 tint=random.choice([(95,118,159,22),(86,100,154,26),(120,108,165,20),(60,102,139,26)])
 cd.ellipse((x-rx,y-ry,x+rx,y+ry),fill=tint)
cloud=cloud.filter(ImageFilter.GaussianBlur(26));sky=Image.alpha_composite(sky,cloud)
# Moon sampled from generated concept as unaltered piece of the approved art set.
moon_src=Image.open(CONCEPTS/'concept-approved.png').convert('RGBA')
# moon in source approx x(250..480), y(30..300). Soft mask cutout, then circular recrop.
moon=moon_src.crop((270,35,490,255)).resize((285,285),Image.Resampling.LANCZOS)
mask=Image.new('L',(285,285));md=ImageDraw.Draw(mask);md.ellipse((20,20,265,265),fill=245)
mask=mask.filter(ImageFilter.GaussianBlur(15));moon.putalpha(mask)
glow=Image.new('RGBA',(W,800));gd=ImageDraw.Draw(glow,'RGBA')
for i in range(190,0,-6):
 a=max(0,int(30*(1-i/200)**1.5))
 gd.ellipse((330-i,176-i,330+i,176+i),fill=(215,205,152,a))
glow=glow.filter(ImageFilter.GaussianBlur(24));sky=Image.alpha_composite(sky,glow)
sky.alpha_composite(moon,(185,32))
# star field: size/brightness vary; keep calm and small.
stars=ImageDraw.Draw(sky,'RGBA')
for i in range(420):
 x=random.randrange(W);y=random.randrange(660)
 if (x-330)**2+(y-176)**2 < 160**2:continue
 radius=random.choice([1,1,1,2,2,3]); opacity=random.randint(95,240)
 stars.ellipse((x-radius,y-radius,x+radius,y+radius),fill=(221,239,253,opacity))
 if radius>=3:
  stars.line((x-7,y,x+7,y),fill=(210,224,255,opacity//2),width=1)
  stars.line((x,y-7,x,y+7),fill=(210,224,255,opacity//2),width=1)

# Sea: original hand-made textured ocean, using thousands of economical strokes pre-rendered offline.
water=Image.new('RGBA',(W,650));wd=ImageDraw.Draw(water,'RGBA')
for y in range(650):
 t=y/650
 wd.line((0,y,W,y),fill=(int(12+5*t),int(55+18*t),int(81+18*t),255))
for i in range(6500):
 x=random.randrange(-10,W);y=random.randrange(650)
 width=random.randint(8,100); thickness=random.choice([1,1,2,3])
 c=random.choice([(121,178,199,39),(149,194,195,36),(9,27,54,51),(47,100,127,35),(239,217,163,12)])
 wd.arc((x,y,x+width,y+random.randint(3,14)),180,345,fill=c,width=thickness)
# moonlight shimmer left of boat route, designed to be peripheral rather than obscure the task.
for i in range(1350):
 y=random.randint(0,650);x=int(random.gauss(325+0.07*y,55+y*.08))
 if x<0 or x>=W:continue
 wd.line((x,y,x+random.randint(2,24),y),fill=(236,212,148,random.randint(10,65)),width=random.randint(1,3))
# layered headlands and a warm distant harbour on right
scene=Image.new('RGBA',(W,H));scene.alpha_composite(sky,(0,0));scene.alpha_composite(water,(0,790))
d=ImageDraw.Draw(scene,'RGBA')
d.polygon([(0,740),(0,650),(105,647),(190,700),(260,742),(325,795)],fill=(12,36,53,245))
d.polygon([(W,735),(W,590),(1840,620),(1780,650),(1705,670),(1630,715),(1590,765)],fill=(10,31,49,250))
# harbour shoreline and 14 warm lighted cottage silhouettes, kept outside the boat lanes
for i in range(15):
 x=1630+i*21;y=650-int((x-1630)*0.08)+random.randrange(-15,16)
 ww=random.randrange(15,32);hh=random.randrange(28,50)
 d.rectangle((x,y-hh,x+ww,y),fill=(14,33,45,255))
 d.polygon([(x-2,y-hh),(x+ww/2,y-hh-15),(x+ww+2,y-hh)],fill=(30,46,60,255))
 for u in range(2):
  xx=x+6+u*10
  d.rectangle((xx,y-hh+12,xx+4,y-hh+19),fill=(255,177,78,220))
# foreground rocky edges, painterly highlights
for side in [0,1]:
 for i in range(32):
  x=random.randint(-100,140) if side==0 else random.randint(1790,2020)
  y=random.randint(970,1560)
  rx=random.randint(30,140);ry=random.randint(22,115)
  col=random.choice([(23,43,54,255),(30,53,62,255),(39,59,66,255),(15,39,51,255)])
  d.ellipse((x-rx,y-ry,x+rx,y+ry),fill=col)
  d.arc((x-rx,y-ry,x+rx,y+ry),215,325,fill=(113,130,121,80),width=3)
# lighthouse foundation and harbour mouth are built as dynamically separate gameplay layers.
scene.convert('RGB').save(A/'seascape.jpg',quality=90,optimize=True,subsampling=0)

# Reusable transparent art sprites, designed at higher-than-logical resolution.
def new_asset(w,h):return Image.new('RGBA',(w,h)),ImageDraw.Draw(Image.new('RGBA',(1,1)))
def save_asset(im,name):im.save(A/name,optimize=True)
# lighthouse at foot of screen; lamp top aligned with coordinate anchor.
L=Image.new('RGBA',(410,630)); d=ImageDraw.Draw(L,'RGBA')
d.ellipse((50,540,370,625),fill=(10,24,35,200))
d.polygon([(143,530),(258,530),(238,130),(167,130)],fill=(205,199,173,255))
for j in range(38):
 y=140+j*10;x1=167-int(j*.63);x2=238+int(j*.53)
 d.line((x1,y,x2,y),fill=(random.choice([(121,123,119,43),(225,223,194,45),(74,100,107,35)])),width=random.randint(1,3))
d.polygon([(165,355),(242,355),(239,390),(164,390)],fill=(138,68,49,255))
d.polygon([(164,466),(253,466),(256,491),(158,491)],fill=(148,64,49,255))
# Door and warmly lit windows
d.rounded_rectangle((185,465,222,531),radius=15,fill=(38,45,55,255),outline=(247,203,117,220),width=3)
for y in (250,350,432):
 d.rounded_rectangle((190,y,213,y+30),radius=7,fill=(255,183,66,250),outline=(76,48,37,255),width=4)
 d.line((200,y,200,y+30),fill=(113,72,48,220),width=2)
# Lantern and broad copper dome
d.rectangle((156,118,252,151),fill=(55,50,48,255))
d.rounded_rectangle((167,56,242,126),radius=12,fill=(255,202,102,240),outline=(91,65,47,255),width=8)
d.polygon([(160,66),(174,34),(199,16),(228,35),(250,66)],fill=(105,68,54,255))
d.line((200,16,200,0),fill=(179,159,101,255),width=5)
for x in (184,205,227): d.line((x,57,x,124),fill=(87,65,44,200),width=5)
d.ellipse((188,78,217,109),fill=(255,246,191,255))
# lamp anchor at image position (205,95) -> screen coordinate (480,640)
save_asset(L,'lighthouse.png')

# Boat 120x85, translucent wake is drawn by renderer, vessel reads distinctly at small scale.
B=Image.new('RGBA',(176,120));d=ImageDraw.Draw(B,'RGBA')
d.ellipse((13,82,162,108),fill=(150,213,227,43))
d.polygon([(12,54),(151,53),(135,90),(35,92)],fill=(23,45,65,255))
d.polygon([(14,53),(149,52),(135,75),(35,79)],fill=(140,57,48,255))
d.polygon([(25,44),(130,44),(155,57),(12,58)],fill=(214,197,139,255))
d.rounded_rectangle((61,12,114,52),radius=6,fill=(243,215,162,255),outline=(54,68,77,255),width=4)
d.rectangle((67,19,85,39),fill=(255,202,83,255));d.rectangle((91,19,107,39),fill=(255,202,83,255))
d.line((118,52,124,5),fill=(211,202,160,255),width=4);d.polygon([(124,7),(124,40),(142,36)],fill=(213,80,62,255))
for i in range(13):
 x=random.randrange(30,152);y=random.randint(59,70)
 d.line((x,y,x+random.randint(2,6),y),fill=(253,221,157,100),width=2)
save_asset(B,'boat.png')

# Parchment panel visual used by book, harbour log, and overlays; local illustration, no UI framework.
P=Image.new('RGBA',(1500,1030));d=ImageDraw.Draw(P,'RGBA')
d.rounded_rectangle((28,25,1470,1004),radius=90,fill=(222,195,143,245),outline=(115,77,51,255),width=20)
for i in range(16):
 d.rounded_rectangle((36+i,33+i,1463-i,996-i),radius=75,outline=(157,124,72,30),width=2)
for x in range(100,1430,16):
 y=random.randrange(80,950)
 d.line((x,y,x+random.randrange(4,30),y+random.randrange(-5,6)),fill=(115,94,57,random.randrange(5,35)),width=2)
d.line((749,55,749,966),fill=(139,111,79,115),width=6)
save_asset(P,'parchment.png')

# Light cones are encoded once in texture, resized/rotated by Ebitengine at runtime.
BW,BH=490,850
beam=Image.new('RGBA',(BW,BH));px=beam.load()
for y in range(BH):
 travel=1-y/(BH-1)
 half=max(6,travel*BW*.45)
 for x in range(BW):
  dist=abs(x-BW/2)/half
  # unlit outside and feathered edge; centre attenuates with distance
  if dist>=1:continue
  alpha=int(52*max(0,1-dist**3)**1.3*(.36+.64*(1-travel*.75)))
  px[x,y]=(255,207,121,alpha)
save_asset(beam,'beam.png')

fog=Image.new('RGBA',(360,230));pix=fog.load()
for y in range(230):
 for x in range(360):
  r=((x-180)/180)**2+((y-115)/115)**2
  if r<1:pix[x,y]=(204,222,237,int(125*(1-r)**2))
save_asset(fog,'fog.png')

# Six originals, consistent art style / recognisable silhouettes; one atlas and named boxes.
S=Image.new('RGBA',(900,360));d=ImageDraw.Draw(S,'RGBA')
# Whale
ox=0
for i in range(10):
 d.arc((ox+18+i,125+i,ox+150-i,213-i),170,352,fill=(95+i*4,175+i*2,193,145),width=3)
d.ellipse((ox+38,137,ox+154,200),fill=(51,114,145,255));d.polygon([(ox+45,170),(ox+12,141),(ox+30,200)],fill=(52,113,143,255))
d.polygon([(ox+110,167),(ox+152,134),(ox+143,187)],fill=(45,100,137,255))
d.ellipse((ox+118,155,ox+127,162),fill=(239,220,169,255))
for i in range(3): d.arc((ox+75+i*8,91-i*5,ox+111+i*10,136),205,305,fill=(169,219,233,220),width=4)
# shooting star
ox=150;d.line((ox+15,202,ox+114,121),fill=(245,209,130,180),width=16)
d.line((ox+30,204,ox+125,123),fill=(255,246,192,230),width=5)
star=[(ox+126+35*math.cos(i*math.pi/5-math.pi/2)*(1 if i%2==0 else .44),125+35*math.sin(i*math.pi/5-math.pi/2)*(1 if i%2==0 else .44))for i in range(10)]
d.polygon(star,fill=(255,230,147,255))
# seal
ox=300;d.ellipse((ox+35,156,ox+144,211),fill=(122,151,163,255));d.ellipse((ox+83,116,ox+150,178),fill=(143,162,171,255))
d.ellipse((ox+126,143,ox+136,152),fill=(20,39,52,255));d.ellipse((ox+143,158,ox+149,164),fill=(40,55,62,255))
d.polygon([(ox+35,175),(ox+9,160),(ox+13,197),(ox+41,193)],fill=(98,135,150,255))
for j in range(3):d.line((ox+135,161+j*4,ox+160,153+j*8),fill=(227,227,208,220),width=2)
# jellyfish
ox=450;d.pieslice((ox+24,101,ox+150,209),180,360,fill=(153,155,225,245),outline=(223,204,243,255),width=4)
for x in range(43,148,20):d.arc((ox+x,179,ox+x+17,259),15,180,fill=(182,180,241,235),width=6)
# owl
ox=600;d.ellipse((ox+39,114,ox+136,214),fill=(143,112,80,255))
d.polygon([(ox+42,132),(ox+34,98),(ox+71,114)],fill=(143,112,80,255));d.polygon([(ox+101,115),(ox+135,96),(ox+133,142)],fill=(143,112,80,255))
for x in (67,112):d.ellipse((ox+x-20,139,ox+x+14,173),fill=(237,216,163,255));d.ellipse((ox+x-5,147,ox+x+5,162),fill=(36,53,59,255))
d.polygon([(ox+83,171),(ox+102,173),(ox+92,191)],fill=(245,184,95,255))
# ferry
ox=750;d.polygon([(ox+10,177),(ox+149,176),(ox+126,209),(ox+37,210)],fill=(231,207,173,255))
d.rounded_rectangle((ox+54,130,ox+125,178),radius=6,fill=(239,229,184,255),outline=(100,117,131,255),width=4)
for x in range(65,124,19):d.rectangle((ox+x,146,ox+x+10,161),fill=(255,189,97,255))
d.line((ox+92,131,ox+92,103),fill=(194,181,135,255),width=4)
save_asset(S,'sightings.png')
# The title uses the approved painting only as title artwork; movable gameplay
# boat and beam are separately authored and never baked into the play scene.
concept=Image.open(CONCEPTS/'concept-approved.png').convert('RGB')
w,h=concept.size; cw=int(h*4/3)
concept=concept.crop(((w-cw)//2,0,(w+cw)//2,h)).resize((1920,1440),Image.Resampling.LANCZOS)
concept=ImageEnhance.Brightness(concept).enhance(.78)
concept.save(A/'title_art.jpg',quality=87,optimize=True,subsampling=0)
print('Assets:',[(p.name,round(p.stat().st_size/1024)) for p in A.iterdir()])
