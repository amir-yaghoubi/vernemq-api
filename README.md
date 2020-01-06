# vernemq-api

APIs around [vernemq](https://docs.vernemq.com/)

# Configurations

You can configure by providing following environment variables:

```

PORT=9595

LOG_LEVEL=error

# internal lru cache size
# used for user acl and password cache
AUTH_LRU_SIZE=2000

# Redis configurations
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0


# Control-Cache header configuration
# value must be a valid duration string
# check here for more information https://golang.org/pkg/time/#ParseDuration
#
# these headers will used by vernemq broker as cache ttl
# more information at https://docs.vernemq.com/plugindevelopment/webhookplugins#caching
#
CC_REGISTER=12h
CC_PUBLISH=12h
CC_SUBSCRIBE=12h


```

# How to Start

Simply just run `docker-compose up --build` in your terminal.
