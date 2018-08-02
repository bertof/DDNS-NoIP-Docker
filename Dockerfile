FROM library/golang

MAINTAINER "berto.f@protonmail.com"

# Arguments
ARG USER_NAME
ARG USER_PASSWD
ARG HOST_NAME
ARG IP
ARG INTERVAL

# Build
RUN useradd -m user
USER user
ADD noip.go .
RUN go build -o NoIP .

ENTRYPOINT ["./NoIP"]

