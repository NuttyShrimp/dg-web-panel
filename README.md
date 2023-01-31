# DG-Panel

Authentication:
- Via ingame token, registered in the CFX server
- Via [discord](https://github.com/ravener/discord-oauth2)

## Stack
- [Gin](https://github.com/gin-gonic/gin)
- [Gorm](https://github.com/go-gorm/gorm)
- [zap](https://github.com/uber-go/zap)
- [go-redis](https://github.com/go-redis/redis)

## Enviromnent

The template config has some predefined resources you can use in your dev env.
You need to set up the following things:
- MySQL/MariaDB with a table named `degrens-panel`
- A discord application which you can create [here](https://discord.com/developers/applications)
- A custom config file with the name of `config.yml`. If you want to change it you need to change the make file
- A graylog instance with a JSON extractor on `full_message`
## Testing
For the testing env to work properly there should be a file name `config.test.yml`. This is to prevent test from interfering with database they should not put any test data in
