sudo apt update
sudo apt upgrade

curl -s https://www.mongodb.org/static/pgp/server-4.4.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/4.4 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.4.list

sudo apt update
sudo apt install -y mongodb-org

sudo systemctl enable mongod
sudo systemctl start mongod
sudo systemctl status mongod

mongo --eval 'db.runCommand({ connectionStatus: 1 })'

#to make it work on remote server
#sudo vim /etc/mongod.conf
#bindIp: 127.0.0.1 -> 0.0.0.0
#uncomment security