# set up 
wget https://go.dev/dl/go1.14.4.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.14.4.linux-amd64.tar.gz

/usr/local/go/bin/go version

everytime put this path 
/usr/local/go/bin/
I don't want to make GOPATH as long as I can :)

/usr/local/go/bin/go get -u github.com/catinello/base62
/usr/local/go/bin/go get -u github.com/lib/pq
/usr/local/go/bin/go get -u github.com/grokify/html-strip-tags-go

deploy command
/usr/local/go/bin/go run main.go >> log/go.log 2>&1 &
ps aux | grep go then kill 2 process
