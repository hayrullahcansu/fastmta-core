# ![bashtop](logo.png)


# Description

Fastmta is high speed distributed mail delivery system. It developed to relay tons of email. 

# Features

* Easy to Use: with inspired menu system.
* Transactional MTA: No matter what happens, the messages you work with will be in a sensible state.
* On-Premise Solution: Run software on hardware located you.
* Cross-platform: Run anywhere linux, windows, darwin etc... with the minumum machine power. 
* High Performance: Provides high delivery performance via non-blocking io operations.
* Distributed: It is horizontally scalable, fault-tolerant.
* Independent Database: It works MSSQL, MySQL, PostgreSQL or SQLite.
* Management UI: Configuration on the fly. 
* Full configurable IPs.
* Multiple Inbound Channel: Smtp, Cli, RestAPI, Inject Queue and Database


# Dependencies

**[GCC](http://tdm-gcc.tdragon.net/download)** -> For Sqlite

**[RabbitMQ](https://www.rabbitmq.com/)** -> A queue tool for the messages.

# Installation

Download Project ```go get -t github.com/hayrullahcansu/fastmta-core```

Go to project folder. It should be ```$GOPATH\github.com\hayrullahcansu\fastmta-core```


In Terminal ```go run main.go```


# Configurability

All options changeable from within UI.

#### app.json 

```json
{
  "database": {
    "driver":"sqlite3",
    "connection":"test.db"
  },
  "ip_addressess": [
    {
      "ip": "127.0.0.1",
      "hostname": "vmta1.localhost",
      "inbound": true,
      "outbound": true
    }
  ],
  "ports": [25, 467, 587],
  "rabbitmq": {
    "host": "rabbitmq-host",
    "port": 5672,
    "username": "username",
    "password": "password",
    "virtual_host": "",
    "exchange_name": ""
  }
}
```

| Param Name | Variable Type | Requirement | Description                         | Value                                                                             |
|------------|---------------|-------------|-------------------------------------|-------------------------------------------------------------------------------------|
| driver       | `:string`     |    yes`*`   | SQL Provider                       |  `sqlite3`, `mysql`, `mssql` `postgresql`                                                                                   |
| connection      | `:string`     |    yes`*`   | Topic pattern can be layout         |  Described below  |

#### Should set connection string with these formats. 
```
sqlite    ->  "/tmp/database_file_name.db"
mysql     ->  "user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
mssql     ->  "sqlserver://username:password@localhost:1433?database=dbnam"
postgres  ->  "host=myhost port=myport user=gorm dbname=gorm password=mypassword"
```

# TODO

- [x] TODO Transactions SQL
- [x] TODO Support Docker
- [x] TODO Go Module
- [ ] TODO Test container
- [ ] TODO Web UI
- [ ] TODO AUTH PLAIN
- [ ] TODO AUTH LOGIN
- [ ] TODO AUTH CRAM-MD5 ```ref: https://www.samlogic.net/articles/smtp-commands-reference-auth.htm``` 
- [ ] TODO Gzip or similar compressor support
- [ ] TODO Support SQL providers
- [ ] TODO Send signal to main process to kill ```ref: initializer.go``` 
- [ ] TODO Check if there was no MX record in DNS, so using A, we should fail and not retry ```ref: bounce_handler.go```
- [ ] TODO report dkim error ```ref: inbound_staging_consumer.go```
- [ ] TODO Implementation bulk sender
- [ ] TODO Save rule, Domain not found ```ref: normal_sender.go```
- [ ] TODO Implementation Smtp Authentication
- [ ] TODO Implementation TLS inbound
- [ ] TODO Implementation LOCALDELIVERY  ```ref: smtp_server.go```
- [ ] TODO Check validity data  ```ref: smtp_server.go```
- [ ] TODO Define all error like dnsError ```ref: agent.go```
- [ ] TODO Add a rule like "this host not valid or unable to connect ```ref: agent.go```
- [ ] TODO Save domain to DB ```ref: domain.go```
- [ ] TODO Save domain to DB ```ref: domain.go```