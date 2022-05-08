#/bin/bash

set -eu

VERSION=0.0.3
wget https://github.com/utam0k/ws-dbg/releases/download/v${VERSION}/wsdbg_${VERSION}_linux_amd64.tar.gz
tar -xzf wsdbg_${VERSION}_linux_amd64.tar.gz
rm wsdbg_${VERSION}_linux_amd64.tar.gz