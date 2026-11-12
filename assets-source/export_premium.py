"""Deterministic art bake for Lighthouse Keeper's painterly presentation.

All painting samples originate from original project artwork in assets-source.
Baking is offline only. The shipping Go executable embeds the resulting images;
no runtime image processor, external asset service, or shader library is needed.
"""
from pathlib import Path
from PIL import Image, ImageDraw, ImageEnhance, ImageFilter
import numpy as np

ROOT = Path(__file__).resolve().parent.parent
SOURCES = ROOT / "assets-source"
OUT = ROOT / "cmd" / "lighthousekeeper" / "assets"
OUT.mkdir(parents=True, exist_ok=True)
A = Image.open(SOURCES / "concept-approved.png").convert("RGB")
B = Image.open(SOURCES / "approved-ocean-study.png").convert("RGB")
C = Image.open(SOURCES / "original-ui-study.png").convert("RGB")
W, H = 1920, 1440


def paint_cutout(img, bounds, outline, feather=2):
    cut = img.crop(bounds).convert('RGBA')
    m=Image.new('L',cut.size,0)
    d=ImageDraw.Draw(m)
    d.polygon([(x-bounds[0],y-bounds[1]) for x,y in outline],fill=255)
    if feather:
        m=m.filter(ImageFilter.GaussianBlur(feather))
    cut.putalpha(m)
    return cut

# The game backdrop is now exported separately by export_master_scene.py.
# No stitched horizon, shader pipeline, or generated water texture.

# Cut an actual illustrated lighthouse lantern from the original maritime study.
# It sits in the foreground with the core gameplay beam originating inside it.
lens=paint_cutout(B,(165,35,465,335),[
    (255,39),(276,39),(284,52),(285,69),(300,81),(307,92),
    (325,105),(345,129),(364,157),(370,161),(370,175),
    (358,183),(358,222),(382,229),(397,246),(397,260),
    (379,268),(375,300),(382,326),(382,334),
    (206,334),(205,316),(213,298),(207,266),(190,257),
    (191,246),(210,228),(219,221),(219,182),(207,179),
    (207,161),(216,151),(226,124),(243,105),(261,92),
    (273,82),(274,67),(257,57)
],feather=2)
# Weather vane is a painted detail, not a patch of surrounding blue sky.
from PIL import ImageDraw as _ID
_lens_mask=lens.getchannel('A')
_lm=_ID.Draw(_lens_mask)
_lm.line((93,18,146,18),fill=255,width=4)
_lm.line((121,5,121,58),fill=255,width=4)
_lens_mask=_lens_mask.filter(ImageFilter.GaussianBlur(1))
lens.putalpha(_lens_mask)
lens=lens.resize((450,450),Image.Resampling.LANCZOS)
lens.save(OUT/'lighthouse.png', optimize=True)
# Miniature brass lanterns are reused in prediction, progress and harbour results.
mini=lens.crop((54,80,380,367)).resize((105,94),Image.Resampling.LANCZOS)
mini.save(OUT/'lantern.png',optimize=True)
# One shared radial sprite for warm light on water, boat and lantern.
G=240
gy,gx=np.ogrid[:G,:G]
gr=np.sqrt(((gx-G/2)/(G/2))**2+((gy-G/2)/(G/2))**2)
ga=np.uint8(np.clip(1-gr,0,1)**2.5*180)
g=np.zeros((G,G,4),dtype=np.uint8)
g[:,:,:3]=[255,199,111];g[:,:,3]=ga
Image.fromarray(g,'RGBA').save(OUT/'glow.png',optimize=True)

# Isolated painted boat hull, cabin and sparse rigging from the approved study.
# Each rope is a thin stroke, avoiding the ugly triangular sky patch of the prototype.
_bounds=(1138,438,1299,571)
boat=B.crop(_bounds).convert('RGBA')
bm=Image.new('L',boat.size,0)
bd=ImageDraw.Draw(bm)
def bp(points):return [(x-_bounds[0],y-_bounds[1]) for x,y in points]
bd.polygon(bp([(1158,516),(1173,507),(1189,502),(1191,479),
               (1204,478),(1237,480),(1244,512),(1267,516),
               (1273,525),(1267,544),(1247,551),(1193,550),
               (1168,534)]),fill=255)
bd.line(bp([(1231,452),(1230,519)]),fill=255,width=4)
bd.line(bp([(1231,455),(1171,518)]),fill=210,width=2)
bd.line(bp([(1231,455),(1257,519)]),fill=210,width=2)
bd.line(bp([(1253,488),(1251,522)]),fill=230,width=3)
bd.line(bp([(1232,465),(1255,471)]),fill=230,width=2)
# Wake needs transparent feathery edges, unlike the crisp boat silhouette.
wm=Image.new('L',boat.size,0);wd=ImageDraw.Draw(wm)
wd.ellipse((7,91,155,130),fill=160)
wm=wm.filter(ImageFilter.GaussianBlur(6))
from PIL import ImageChops
bm=ImageChops.lighter(bm,wm).filter(ImageFilter.GaussianBlur(1.1))
boat.putalpha(bm)
boat=boat.resize((270,230),Image.Resampling.LANCZOS)
boat.save(OUT/'boat.png',optimize=True)

# Keep the painterly detail instead of stretching a flat-colour silhouette.
boat=boat.resize((256,256),Image.Resampling.LANCZOS)
boat.save(OUT/'boat.png',optimize=True)

# A physically readable gold cone. Feathered falloff rather than flat geometry.
BW,BH=640,1420
bx=np.arange(BW,dtype=np.float32)[None,:]
by=np.arange(BH,dtype=np.float32)[:,None]
fwd=1-(by/BH)
half=7+260*fwd
ux=np.abs((bx-BW/2)/half)
edge=np.clip((1-ux)/.36,0,1)
alpha=(.20*(1-fwd**1.9)*edge**1.1*255).astype('uint8')
rgba=np.zeros((BH,BW,4),dtype='uint8')
rgba[:,:,0]=255;rgba[:,:,1]=211;rgba[:,:,2]=132;rgba[:,:,3]=alpha
beam=Image.fromarray(rgba,'RGBA');beam.save(OUT/'beam.png',optimize=True)

# Fog: a softly textured cloud bank with no visible hard circular edge.
FW,FH=560,350
fx=np.linspace(-1,1,FW)[None,:];fy=np.linspace(-1,1,FH)[:,None]
rad=fx*fx*.9+fy*fy*1.4
raw=np.clip((1-rad)/.7,0,1)**2
fog=np.zeros((FH,FW,4),dtype='uint8')
fog[:,:,0]=180;fog[:,:,1]=207;fog[:,:,2]=210;fog[:,:,3]=np.uint8(raw*125)
Image.fromarray(fog,'RGBA').filter(ImageFilter.GaussianBlur(7)).save(OUT/'fog.png', optimize=True)

# Illustrated sighting studies, individually painted in the approved visual sheet.
# Feather their outer edges into live scenery; the game alone controls their state.
tiles=[(57,786,167,866),  # whale
       (291,894,397,963), # shooting star
       (175,786,280,866), # seal
       (291,783,397,868), # jellyfish
       (57,892,169,964),  # owl
       (173,892,281,964)] # ferry
TW,TH=240,240
atlas=Image.new('RGBA',(TW*6,TH),(0,0,0,0))
for i,box in enumerate(tiles):
    tile=C.crop(box).resize((222,196),Image.Resampling.LANCZOS).convert('RGBA')
    tx=np.linspace(-1,1,tile.width)[None,:];ty=np.linspace(-1,1,tile.height)[:,None]
    fall=np.clip((1-(tx*tx*.9+ty*ty*1.0))/.55,0,1)**1.6
    tile.putalpha(Image.fromarray(np.uint8(fall*255),'L'))
    atlas.alpha_composite(tile,(TW*i+9,18))
atlas.save(OUT/'sightings.png',optimize=True)

# Painterly book pages, updated from the source underlay with delicate paper grain.
# Keep the printed interior blank because each page is populated by the game.
PW,PH=1500,1030
p=Image.new('RGBA',(PW,PH),(0,0,0,0)); pd=ImageDraw.Draw(p)
pd.rounded_rectangle((18,28,PW-18,PH-30),radius=65,fill=(73,40,23,244))
pd.rounded_rectangle((30,34,PW-30,PH-37),radius=55,fill=(151,99,55,255))
noise=np.random.default_rng(1402).normal(0,4,(PH,PW,1)).astype(np.float32)
base=np.zeros((PH,PW,4),dtype=np.uint8);base[:,:,:3]=[229,204,157];base[:,:,3]=255
base[:,:,:3]=np.uint8(np.clip(base[:,:,:3].astype(np.float32)+noise,0,255))
paper=Image.fromarray(base,'RGBA')
mask=Image.new('L',(PW,PH));md=ImageDraw.Draw(mask)
md.rounded_rectangle((56,51,PW-56,PH-61),radius=50,fill=255)
p.alpha_composite(Image.composite(paper,Image.new('RGBA',(PW,PH)),mask))
pd=ImageDraw.Draw(p,'RGBA')
pd.line((750,61,750,963),fill=(125,80,53,175),width=9)
for k,alpha in [(38,90),(46,45),(51,28)]:
    pd.rounded_rectangle((k,k+5,PW-k,PH-k-20),radius=44,outline=(127,76,45,alpha),width=3)
for x in (108,1395):
    for y in (112,923):
        pd.arc((x-25,y-27,x+25,y+27),0,270,fill=(118,72,43,140),width=3)
p.save(OUT/'parchment.png',optimize=True)

# Curved brass/navy UI rail. The adapter draws all wording and controls.
QW,QH=960,220
u=np.zeros((QH,QW,4),dtype=np.uint8)
ry,rx=np.mgrid[:QH,:QW]
sh=np.clip(1-(ry/QH)**1.9*.45,0,1)
u[:,:,:3]=np.uint8(np.stack((14+20*sh,31+25*sh,49+23*sh),axis=-1))
u[:,:,3]=235
ui=Image.fromarray(u,'RGBA');ud=ImageDraw.Draw(ui,'RGBA')
ud.rounded_rectangle((12,10,QW-13,QH-12),radius=27,outline=(255,205,116,220),width=4)
ud.rounded_rectangle((23,20,QW-24,QH-23),radius=18,outline=(165,122,76,150),width=2)
for x in (42,QW-43):
    ud.ellipse((x-7,38,x+7,52),fill=(212,166,85,175))
    ud.ellipse((x-7,QH-53,x+7,QH-39),fill=(212,166,85,175))
ui.save(OUT/'ui_rail.png',optimize=True)

# Title keeps the original approved premium painting; a matching soft grade.
source=A.convert('RGB'); sw,sh=source.size
cw=int(sh*4/3)
source=source.crop(((sw-cw)//2,0,(sw+cw)//2,sh)).resize((W,H),Image.Resampling.LANCZOS)
source=ImageEnhance.Brightness(source).enhance(.77)
source.save(OUT/'title_art.jpg',quality=89,optimize=True,subsampling=0)
print('Baked', [(p.name,round(p.stat().st_size/1024)) for p in sorted(OUT.iterdir()) if p.suffix in ('.png','.jpg')])
