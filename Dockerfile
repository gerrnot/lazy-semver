FROM golang:1 as builder
WORKDIR /workdir
COPY . .
RUN go build && go test && ls -la
#-----------------------------------------------------------------------------------------------------------------------
FROM alpine as runner
# fix for "no such file or directory"
# https://stackoverflow.com/questions/51508150/standard-init-linux-go190-exec-user-process-caused-no-such-file-or-directory
RUN apk add --no-cache libc6-compat

# setup non-root user, inspired by
# https://gist.github.com/avishayp/33fcee06ee440524d21600e2e817b6b7
ENV USER=non-root
ENV HOME /home/${USER}
RUN apk add --update sudo &&\
    adduser -D ${USER} \
    && echo "${USER} ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/${USER} \
    && chmod 0440 /etc/sudoers.d/${USER}
USER ${USER}
WORKDIR ${HOME}

# copy go app
COPY --from=builder /workdir/lazy-semver /bin/lazy-semver
ENTRYPOINT ["lazy-semver"]
