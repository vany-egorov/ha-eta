Requirements
------------

- docker with docker-compose
- golang (for development)

Diagrams
--------

Architecture overview:
```
                                              +-------------+       +----------+
                                              |             |       |          |
                                              |   nginx(2s) |       |    cache |        +-------------+
                                              |             |       |          |        |             |
                                        +---->+ cache-01    +-------+ node-01  +--------+             |
       +------------------+             |     |             |       |          |        |             |
       |                  |             |     +-------------+       +----------+        |             |
+------+ balancer         +-------------+                                               |             |
|      |                  |             |     +-------------+       +-----------+       |             |
|      +------------------+             |     |             |       |           |       |             |
|                           round-robin +---->+ cache-02    +------>+ node-02   +------->  GEO API    |
|      +------------------+             |     |             |       |           |       |             |
|      |                  |             |     +-------------+       +-----------+       |             |
+----->+ balancer.backup  |             |                                               |             |
       |                  |             |     +-------------+       +-----------+       |             |
       +------------------+             |     |             |       |           |       |             |
                                        +---->+ cache-03    +------>+ node-03   +------>+             |
VIP / keepalived                              |             |       |           |       |             |
                                              +-------------+       +-----------+       +-------------+
```

- balancer: via nginx round-robin upstream balancing. HAProxy can also be used;
- cache-0[1-3]: simple http 200 cache with 2s ttl;
- node-0[1-3]: worker nodes;
- GEO API: API for geo location and prediction;

How to Use
----------

[swagger api-vi spec](./blob/master/src/assets/api-v1.yml)

Clone this app with:
``` bash
$ git clone https://github.com/vany-egorov/ha-eta.git
```

Run under docker (**docker + docker-compose** required).
This will build app, run tests, start docker-containers.
Cluster listen on **9179** port
``` bash
$ cd ha-eta
$ chmod +x ./run.sh
$ ./run.sh
```

app can simply be tested with:
```bash
$ curl -XGET 'http://127.0.0.1:9179/api/v1/eta/min?lat=55.752992&lng=37.618333'

# or
$ watch -n -1 "curl -XGET 'http://127.0.0.1:9179/api/v1/eta/min?lat=55.752992&lng=37.618333'"
```

build app and run without docker (**golang** required) for development:
``` bash
$ cd ha-eta/src
$ chmod +x ./build.sh
$ ./build.sh

$ ./ha-eta --help

# сustom port
$ ./ha-eta -p 8888

# disable cache and custom port
$ ./ha-eta -p 8888 --do-not-cache-points --do-not-cache-etas
```
