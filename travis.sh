#!/bin/bash
if [ -z "${GOPATH}" ]; then
        export GOPATH=/home/travis/gopath
fi
set -e
echo "-----------------------> Build dep and mkdep"
 (cd dep && go build)
 (cd mkdep && go build)

echo "-----------------------> guid test"
 (cd pkg && CGO_ENABLED=0 go test ./...)

echo "-----------------------> dep and mkdep test"
 dep/dep < testdep.sec > testdepfile
 mkdep/mkdep < testdepfile > testdep2.sec
 diff testdep.sec testdep2.sec

echo "-----------------------> go vet"
 go tool vet dep/dep.go mkdep/mkdep.go pkg/guid/guid.go pkg/guid/guid_test.go

