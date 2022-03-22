# core-api-template
Basic template to run core-api written in Go w/ go-flow


- REST API: `http://localhost:5000`

## Environment variables

| Variable                       | Required | Default Value   | Description                                         |
| -------------------------------| -------- | --------------- | --------------------------------------------------- |
| ENV                            | NO       | production      | indicates in which environment app is running       |
| ADDR                           | NO       | 5000            | service port                                        |
| LOG_LEVEL                      | NO       | error           | logging level                                       |
| DB_DIALECT                     | NO       | mysql           | Database dialect                                    |
| DB_USER                        | YES      |                 | Database User                                       |
| DB_PASS                        | YES      |                 | Database Password                                   |
| DB_HOST                        | YES      |                 | Database Host name                                  |
| DB_PORT                        | NO       | 3306            | Database Host name                                  |
| DB_NAME                        | YES      |                 | Database name                                       |
| DB_PARAMS                      | NO       |                 | Database connection params                          |
| RSA_PUBLIC_KEY                 | YES      |                 | RSA Public key file path needed for Authentication  |
| RSA_PRIVATE_KEY                | YES      |                 | RSA Private key file path needed for Authentication |
| RSA_PRIVATE_KEY_PASSWORD       | NO       |                 | RSA Private key password                            |




## Built With

- [go](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
- [go flow](https://github.com/go-flow/flow/v2/) - High Performance minimalist web framework for gophers
