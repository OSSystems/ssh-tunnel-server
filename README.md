# Installing

```
$ git clone https://github.com/OSSystems/ssh-tunnel-server.git
$ cd ssh-tunnel-server
$ docker-compose build
```

# Setup

* Add the device public key to ./ssh-forward/ssh/authorized_keys
  * You will need to store the private key on your local machine -- the one that you will use to access the remove device
* Add the server private key to ./ssh-gateway/key.pem
* Open the following Ports for Inbound Traffic

Port | Type | Protocol Description
-----|----------|---------------------
22 | TCP | SSH for device connnection
32768 - 60999 | TCP | Incoming Device Port
1883 | TCP | MQTT
2221-2222 | TCP | SSH for server-tunner-server manual connection

* NOTE: Change SSHD config from Port 22 to 2222
```
sudo vi /etc/ssh/sshd_config
# uncomment # Port 22 and then change 22 to 2222
# save and exit sshd_config file
# restart sshd service
sudo systemctl restart sshd
```

# Running

To have the containers run in daemon mode (not print out the log outputs)
```
docker-compose start
```

This prints out the logs of each container and does not run in the background
```
docker-compose up
```

# Example Connection

```
ssh gateway.kcam-service.com -l root@my-remote-device
```

### Notes:

* hostname is case-sensitive (E.g. my-remote-device).
* As stated in the section Setup, you will need to store the private key on your local machine -- the one that you will use to access the remove device
This key grants you access to the ssh-tunnel-server
