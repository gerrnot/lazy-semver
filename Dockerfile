FROM golang:1 as builder
WORKDIR /workdir
COPY go.mod go.sum *.go ./
RUN go build lazy-semver.go
#-----------------------------------------------------------------------------------------------------------------------
FROM alpine as runner
COPY --from=builder /workdir/lazy-semver /bin
CMD ["lazy-semver"]
