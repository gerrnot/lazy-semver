# lazy-semver
docker image for reading version strings from various sources and producing a lazy SemVer thereof

# usage
`docker run -v "${PWD}:/mnt/input" test --inputFilePath /mnt/input/test/resources/version.txt`

# developer information

## compile & execute e2e tests
```
go build && go test
```

## release
build the docker image and test it:
```
docker build . -t test && docker run -v "${PWD}:/mnt/input" test --inputFilePath /mnt/input/test/resources/version.txt
```
should print a SemVer in last line of output!

define new version(s) and upload to dockerhub:
```
docker tag test gernotfeichter/lazy-semver:0.0.3
docker tag test gernotfeichter/lazy-semver:0.0
docker tag test gernotfeichter/lazy-semver:0
docker tag test gernotfeichter/lazy-semver:latest

docker push gernotfeichter/lazy-semver:0.0.3
docker push gernotfeichter/lazy-semver:0.0
docker push gernotfeichter/lazy-semver:0
docker push gernotfeichter/lazy-semver:latest
```
