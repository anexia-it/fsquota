#!/bin/bash

WorkDir="$PWD/.."
BinDir=$WorkDir/bin
Platforms="linux"
Architectures="amd64"
GoVersion=1.17.8

echo ""
echo " ---- [Start] deploy for platforms ($Platforms [$Architectures]) [Start]-----"
echo ""
echo "Work dir     : $WorkDir"
echo "Binaries dir : $BinDir"

docker run --rm -it -v "$WorkDir":/usr/src/myapp -v "$GOPATH":/go -w /usr/src/myapp golang:$GoVersion bash -c '
echo $PWD && \
rm -rf bin && \
mkdir bin && \
for GOOS in linux; do
  for GOARCH in 386 amd64; do
    export GOOS=$GOOS
    export GOARCH=$GOARCH
    echo "Building $GOOS-$GOARCH"
    go build -o bin/fsqm-$GOARCH cmd/fsqm/*.go
  done
done
' && \
echo "Build complete" && \
echo "" && \
echo "Permission adding:" && \
echo "chown -R root:root $BinDir" && \
echo "chmod -R 777 $BinDir" && \
echo "" &&\
chown -R root:root "$BinDir" && \
chmod -R 777 "$BinDir" && \
echo "" && \
echo "ls -la $BinDir:" && \
ls -la "$BinDir" && \
echo "" && \
echo "EnvPath" && \
export PATH=$PATH:"$BinDir" && \
echo "" && \
echo "\$Path:\"$BinDir\"" && \
"$BinDir"/fsqm-amd64 && \
echo ""  && \
echo "export PATH=\$PATH:\"$BparseLimitsFlaginDir\""  && \
echo ""  && \
echo " ---- [End] deploy for all platforms ($Platforms [$Architectures]) [end]-----"  && \
echo ""
