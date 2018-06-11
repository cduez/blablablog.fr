#!/bin/bash

set -e
shopt -s nullglob

source=~/pictures/blog/
destination=~/dev/go/src/github.com/cduez/blog/assets/images/
thumbify=~/dev/go/src/github.com/cduez/blog/scripts/thumbify.sh

pushd "${source}" > /dev/null

for d in */
do
	if [ ! -d "${destination}/${d}" ]
	then
		echo "copy new directory ${d}"
		cp -r "${d}" "${destination}"

		index=1
		pushd "${destination}/${d}" > /dev/null
		for f in *.JPG *.jpg
		do
			mv "${f}" "${index}.jpg"
			index=$((index + 1))
		done

		${thumbify} "${destination}/${d}"
		popd > /dev/null
	fi
done

popd > /dev/null

echo "done!"

exit 0
