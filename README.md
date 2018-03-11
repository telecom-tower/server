# Telecom Tower Server 2018

## Installation

```
echo "deb https://dl.bintray.com/telecom-tower/deb stretch main contrib rpi" > /etc/apt/sources.list.d/telecom-tower.list
curl https://bintray.com/user/downloadSubjectPublicKey?username=bintray | sudo apt-key add -
sudo apt update
sudo apt install telecom-tower-server
```
