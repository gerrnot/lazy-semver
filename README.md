# lazy-semver
docker image for reading version strings from various sources and producing a lazy SemVer thereof

# usage
`docker run -v "${PWD}:/mnt/input" test --filePath /mnt/input/test/resources/version.txt`

# developer information

## compile & execute e2e tests

execute in project root folder:
```
go build && go test
```
