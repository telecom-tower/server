# Telecom Tower Server 2018

## Installation

### Add the telecom-tower package repository

```
echo "deb https://dl.bintray.com/telecom-tower/deb stretch main contrib rpi" > /etc/apt/sources.list.d/telecom-tower.list
curl https://bintray.com/user/downloadSubjectPublicKey?username=bintray | sudo apt-key add -
sudo apt update
sudo apt install telecom-tower-server
```

### After the installation or after an update, you should restart the SystemD daemon:

```
sudo systemctl daemon-reload
```

### Enable the server so that is starts automatically on boot:

```
sudo systemctl enable telecom-tower
```

### Manually start the server

```
sudo systemctl start telecom-tower
```
or
```
sudo service telecom-tower start
```

### Manually restart the server

```
sudo systemctl restart telecom-tower
```
or
```
sudo service telecom-tower restart
```
