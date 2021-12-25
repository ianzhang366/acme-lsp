#!/bin/sh

set -ex

DIR=golang_org_x_tools
REPO=https://go.googlesource.com/tools

# golang.org/x/tools version that gopls/v0.2.2 depends on
COMMIT=952e2c076240

rm -rf $DIR
git clone $REPO
(
    cd tools
    git checkout $COMMIT
)

mv tools/internal/lsp/protocol/LICENSE lsp/protocol/LICENSE
mv tools/internal/lsp/protocol/context.go lsp/protocol/context.go
mv tools/internal/lsp/protocol/doc.go lsp/protocol/doc.go
mv tools/internal/lsp/protocol/enums.go lsp/protocol/enums.go
mv tools/internal/lsp/protocol/log.go lsp/protocol/log.go
mv tools/internal/lsp/protocol/protocol.go lsp/protocol/protocol.go
mv tools/internal/lsp/protocol/span.go lsp/protocol/span.go
mv tools/internal/lsp/protocol/tsclient.go lsp/protocol/tsclient.go
mv tools/internal/lsp/protocol/tsprotocol.go lsp/protocol/tsprotocol.go
mv tools/internal/lsp/protocol/tsserver.go lsp/protocol/tsserver.go

mkdir $DIR
mv tools/LICENSE $DIR/LICENSE
mv tools/internal/jsonrpc2 $DIR/jsonrpc2
mv tools/internal/jsonrpc2_v2 $DIR/jsonrpc2_v2
mv tools/internal/span $DIR/span
mv tools/internal/telemetry $DIR/telemetry
mv tools/internal/xcontext $DIR/xcontext

(
    cd tools
    echo "Packages in this directory is copied from golang.org/x/tools/internal (commit $COMMIT)."
    #git show -s --format='(commit %h on %ci).'
) > $DIR/README.txt

find $DIR -name '*.go' | xargs sed -i 's!golang.org/x/tools/internal!github.com/fhs/acme-lsp/internal/golang_org_x_tools!'

rm -rf tools
