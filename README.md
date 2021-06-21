# go-simple-crud
postgres orm &amp;&amp; router
add .env 
export DIALECT="postgres"
export HOST="localhost"
export DBPORT="5432"
export USER="****"
export NAME="***"
export PASSWORD="****"


sudo apt install postgresql postgresql-contrib
sudo -u postgres psql
sudo passwd postgres
sudo service postgresql start

source .env
go mod init main.go
go run main.go 
