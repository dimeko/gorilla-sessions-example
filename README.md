#### Software Security Project 2

Creating the TLS certificates:
```
openssl req -x509 -nodes -newkey rsa:2048 -keyout private.key -out certificate.crt -days 365 -config openssl.cnf
```

Running the application stack:
```
docker-compose up app postgres mailpit
```

Initializing database:
```
docker-compose up migrate-up
```

Destroy database:
```
docker-compose up migrate-drop
```