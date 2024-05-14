FROM golang:1.18.2

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN apt-get update

RUN apt-get install -y emacs-nox
RUN apt-get install net-tools
RUN apt-get -y install openssl
RUN mkdir /root/.tls
# RUN openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 \
#     -keyout /root/.tls/private.key \
#     -out /root/.tls/certificate.crt

RUN go mod download \ 
&& go get soft-sec

EXPOSE 6028

CMD [ "go", "run", ".", "server" ]