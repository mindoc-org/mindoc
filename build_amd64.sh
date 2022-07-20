rm mindoc_linux_amd64 mindoc_linux_musl_amd64
rm -rf ../mindoc_linux_amd64/

export GOARCH=amd64
export GOOS=linux
export CC=/usr/bin/gcc

export TRAVIS_TAG=v2.1-beta.6

go mod tidy -v
go build -v -o mindoc_linux_amd64 -ldflags="-linkmode external -extldflags '-static' -w -X 'github.com/mindoc-org/mindoc/conf.VERSION=$TRAVIS_TAG' -X 'github.com/mindoc-org/mindoc/conf.BUILD_TIME=`date`' -X 'github.com/mindoc-org/mindoc/conf.GO_VERSION=`go version`'"
./mindoc_linux_amd64 version

mkdir ../mindoc_linux_amd64
cp -r * ../mindoc_linux_amd64
cd ../mindoc_linux_amd64
rm -rf cache commands controllers converter .git .github graphics mail models routers utils runtime conf/*.go
rm appveyor.yml docker-compose.yml Dockerfile .travis.yml .gitattributes .gitignore go.mod go.sum main.go README.md simsun.ttc start.sh sync_host.sh build_amd64.sh build_musl_amd64.sh
zip -r mindoc_linux_amd64.zip conf static uploads views lib mindoc_linux_amd64 LICENSE.md
mv ./mindoc_linux_amd64.zip ../
