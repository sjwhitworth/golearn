#!/bin/sh

set -e

workdir=.cover
profile=${workdir}/cover.out
mode=count

generate_cover_data() {
    rm -rf ${workdir}
    mkdir ${workdir}

    for pkg in "$@"; do
        f="${workdir}/$(echo ${pkg} | tr / -).cover"
        go test -covermode=${mode} -coverprofile=${f} ${pkg}
    done

    echo "mode: ${mode} > ${profile}"
    grep -h -v "^mode:" "${workdir}"/*.cover >>"${profile}"
}

generate_cover_data $(go list ./...)
mv ${profile} coverage.txt
