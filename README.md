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


# TODO

- [x] TODO Transactions SQL
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