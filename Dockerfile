FROM golang:1.18.2

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod download \ 
&& go get soft-sec

EXPOSE 6028

CMD [ "go", "run", ".", "server" ]