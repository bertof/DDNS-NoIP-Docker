FROM library/golang
RUN useradd -m user
USER user
#WORKDIR /home/user
ADD noip.go .
RUN go build -o NoIP .

ENTRYPOINT ["./NoIP"]

