#!/bin/sh

IMG_DIR="$(dirname "$0")/../static/img"

if ! command -v cwebp >/dev/null 2>&1; then
  echo "Error: cwebp is not installed. Install it via 'brew install webp'." >&2
  exit 1
fi

find "$IMG_DIR" -type f \( -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" \) | while read -r img; do
  # Skip avatar-icon.png
  [ "$(basename "$img")" = "avatar-icon.png" ] && continue

  webp="${img%.*}.webp"

  # Skip if WebP already exists
  [ -e "$webp" ] && continue

  echo "Converting $img to $webp"
  cwebp "$img" -o "$webp"
done
