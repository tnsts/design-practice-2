#!/bin/sh
set -ex
wget https://github.com/tnsts/design-practice-2/blob/master/examples/binary/out/bin/bood_rebase
sudo chmod +x bood_rebase
mv bood_rebase $GOPATH/bin
