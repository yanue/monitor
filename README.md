
# docker test
    docker build -t centos-go .

# server
    docker run -d --privileged -v /Users/yanue/golang/src:/opt/go/src --name test --hostname=test -p 3812:3812 centos-go /usr/sbin/init
# client
    docker run -d --privileged -v /Users/yanue/golang/src:/opt/go/src --name test1 --hostname=test1 centos-go /usr/sbin/init
    docker run -d --privileged -v /Users/yanue/golang/src:/opt/go/src --name test2 --hostname=test2 centos-go /usr/sbin/init
    docker run -d --privileged -v /Users/yanue/golang/src:/opt/go/src --name test3 --hostname=test3 centos-go /usr/sbin/init

# build
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
    
# enter docker
    docker exec -it test bash
    