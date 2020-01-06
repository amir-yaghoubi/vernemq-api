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

# Webhooks

You can point following **Vernemq webhooks** to our endpoints:

- Point [auth_on_register](https://docs.vernemq.com/plugindevelopment/webhookplugins#auth_on_register) to `{{API_URL}}/api/v1/auth/register/`

- Point [auth_on_subscribe](https://docs.vernemq.com/plugindevelopment/webhookplugins#auth_on_subscribe) to `{{API_URL}}/api/v1/auth/subscribe/`

- Point [auth_on_publish](https://docs.vernemq.com/plugindevelopment/webhookplugins#auth_on_publish) to `{{API_URL}}/api/v1/auth/publish/`

> For more information about how to configure Vernemq webhooks [read here](https://docs.vernemq.com/plugindevelopment/webhookplugins#configuring-webhooks)

# Documentation

Swagger documantation available at `{{API_URL}}/swagger/index.html`
