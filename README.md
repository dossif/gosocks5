# gosocks5

SOCKS5 Proxy server implementation.

Supported auth types:
* None - no auth
* Static - auth with static user and pass
* Ldap - auth with remote ldap

Build binary:

    go build -o ./bin/gosocks5 -trimpath ./cmd/gosocks5

Binary usage:

    ./bin/gosocks5 -h

Deploy to docker:
    
    CGO_ENABLED=0 GOOS=linux go build -o ./build/gosocks5 -trimpath ./cmd/gosocks5
    docker compose --file ./deployment/docker-compose.local.yaml up --detach --build --force-recreate