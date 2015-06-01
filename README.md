[![GoDoc](https://godoc.org/github.com/resourced/resourced-master?status.svg)](http://godoc.org/github.com/resourced/resourced-master)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](LICENSE.md)
[![Imgur Album](https://img.shields.io/badge/images-imgur-blue.svg?style=flat)](http://imgur.com/a/MKyFr#0)


**ResourceD Master** receives server data from ResourceD agents and serves them as HTTP+JSON.

![Signup](http://i.imgur.com/UcNmeTF.png)

![Access Tokens](http://i.imgur.com/3H9ONza.png)

![Hosts](http://i.imgur.com/aTEOlA3.png)


## Installation for users

1. Install PostgreSQL 9.4.x

2. Install Go 1.4.x, git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`

3. Create PostgreSQL database.
    ```
    sudo su - postgres
    createuser -P -e resourced
    createdb --owner=resourced resourced-master

    # Make sure user, password, and pg_hba.conf are configured correctly.
    ```

4. [Download the tar.gz](https://github.com/resourced/resourced-master/releases), unpack it, and run the binary using init/systemd/supervisord. You can follow the examples of init scripts [here](https://github.com/resourced/resourced-master/tree/master/scripts/init).


## Installation for developers/contributors

See [INSTALL.md](docs/contributors/INSTALL.md) and [BUILD.md](docs/contributors/BUILD.md)


## ENV variables as configuration

* **RESOURCED_MASTER_ADDR:** HTTP server host:port. Default: `":55655"`

* **RESOURCED_MASTER_CERT_FILE:** Path to cert file. Default: `""`

* **RESOURCED_MASTER_KEY_FILE:** Path to key file. Default: `""`

* **RESOURCED_MASTER_DSN** PostgreSQL DSN. Default: `"postgres://$PG_USER@$localhost:5432/resourced-master?sslmode=disable"`

* **RESOURCED_MASTER_COOKIE_SECRET** Cookie secret key. Default: `"$SEE_THE_CODE"`


## RESTful Endpoints

Every HTTP request requires AccessToken passed as user. Example:
```
curl -u 0b79bab50daca910b000d4f1a2b675d604257e42: https://localhost:55655/api/hosts
```

* **GET** `/api/hosts` Displays list of all hosts by access token.

* **POST** `/api/hosts` Submit JSON data from 1 host. The JSON payload format is defined by `type ResourcedPayload struct`. See: [/dal/host.go#L24](https://github.com/resourced/resourced-master/blob/master/dal/host.go#L24)


## Querying

You can query hosts data using SQL-like language.

There are 3 fields to query from: `name`, `tags`, and `data`.

Currently, you can only use *AND* conjunctive operators.


**Query by hostname**

* Exact match: `name = "localhost"`

* Starts-with match: `name ~^ "awesome-app-"`

* Regex match, case insensitive: `name ~* "awesome-app-"`

* Regex match, case sensitive: `name ~ "awesome-app-"`

* Regex match negation, case sensitive: `name !~ "awesome-app-"`

* Regex match negation, case insensitive: `name !~* "awesome-app-"`


**Query by tags**

* Contains the following tags: `tags = ["app", "django", "staging"]`

**Query by data**

To craft data query, start with ResourceD path and then use "." delimited separator as you get deeper into the JSON structure.

For example, let's say your resourced agent shipped `/free` data:
```json
{"/free": {"Swap": {"Free": 0, "Used": 0, "Total": 0}, "Memory": {"Free": 1346609152, "Used": 7243325440, "Total": 8589934592, "ActualFree": 3666075648, "ActualUsed": 4923858944}}}
```

You can then query "Swap": "Used": data: `/free.Swap.Used > 10000000`
