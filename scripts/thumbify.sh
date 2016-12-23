#!/bin/bash

set -e

thumbquality=70
thumbsize=200
quality=80
size=1024

pushd $1 > /dev/null

if [ ! -d "./thumbs" ]
then
   mkdir thumbs
fi

shopt -s nullglob
for img in *.jpg
do
  mogrify -auto-orient "${img}"

  width=$(identify -format '%w' "${img}")
  height=$(identify -format '%h' "${img}")

  if [ ${width} -gt ${height} ]
  then
    thumbresize="${thumbsize}x"
    resize="${size}x"
  else
    thumbresize="x${thumbsize}"
    resize="x${size}"
  fi

  echo "Generate thumbnail for ${img}"
  convert ${img} -thumbnail ${thumbresize} -gravity center -crop ${thumbsize}x${thumbsize}+0+0 -quality ${thumbquality} "./thumbs/${img}"

  echo "Convert image for ${img}"
  mogrify -resize ${resize} -quality ${quality} "${img}"
done

popd > /dev/null

exit 0
