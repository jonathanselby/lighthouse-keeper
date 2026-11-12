"""Offline single-scene Lighthouse Keeper painting. Not part of the game runtime.

Requires Python Pillow, NumPy and OpenCV on the authoring computer only.
"""
from pathlib import Path
import cv2
import numpy as np
from PIL import Image

ROOT = Path(__file__).resolve().parent.parent
original = Image.open(ROOT / "assets-source/approved-ocean-study.png").convert("RGB")
scene = np.array(original.crop((395, 0, 1640, 928)))
height, width = scene.shape[:2]

# Replace the original painting's decorative boat with water from the same
# horizon. The real game boat is a separate sprite and must not be duplicated.
cy, cx, rw, rh = 515, 818, 91, 100
donor_x = 520
donor = scene[cy-rh:cy+rh, donor_x-rw:donor_x+rw].copy()
mask = np.zeros((2*rh, 2*rw), dtype=np.uint8)
cv2.ellipse(mask, (rw, rh), (rw-5, rh-6), 0, 0, 360, 255, -1)
new = cv2.seamlessClone(
    cv2.cvtColor(donor, cv2.COLOR_RGB2BGR),
    cv2.cvtColor(scene, cv2.COLOR_RGB2BGR),
    mask, (cx, cy), cv2.NORMAL_CLONE,
)
scene = cv2.cvtColor(new, cv2.COLOR_BGR2RGB)

# Cool the baked warm glow in the left-hand sky. Runtime beam remains separate.
y, x = np.mgrid[:height, :width]
upper, lower = 159 + .12*x, 295 + .235*x
soft = np.minimum(np.clip((y-upper)/55, 0, 1), np.clip((lower-y)/65, 0, 1))
soft *= np.clip((760-x)/250, 0, 1)
r, g, b = [scene[:, :, i].astype(np.float32) for i in range(3)]
weight = (soft * np.clip((r-b-6)/90, 0, 1) * .70).astype(np.float32)
cool = np.stack([.56*r+.44*b, .74*g+.26*b, .64*b+.36*g], axis=-1)
scene = np.clip(scene*(1-weight[...,None])+cool*weight[...,None], 0, 255).astype("uint8")

out = ROOT / "cmd/lighthousekeeper/assets/seascape.jpg"
out.parent.mkdir(parents=True, exist_ok=True)
Image.fromarray(scene, "RGB").resize((1920, 1440), Image.Resampling.LANCZOS).save(
    out, quality=85, optimize=True, subsampling=0,
)
print(f"{out}: {out.stat().st_size:,} bytes")
