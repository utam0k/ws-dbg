#/bin/bash
set -eu

VERSION=0.0.1
wget https://github.com/utam0k/ws-dbg/releases/download/v${VERSION}/ws-dbg_${VERSION}_linux_amd64.tar.gz
tar -xzf ws-dbg_${VERSION}_linux_amd64.tar.gz
rm ws-dbg_${VERSION}_linux_amd64.tar.gz