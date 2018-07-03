# Installing

```
$ git clone https://github.com/OSSystems/ssh-tunnel-server.git
$ cd ssh-tunnel-server
$ docker-compose build
```

# Setup

* Add the device public key to ./ssh-forward/ssh/authorized_keys
* Add the server private key to ./ssh-gateway/key.pem

# Running

```
docker-compose up
```
