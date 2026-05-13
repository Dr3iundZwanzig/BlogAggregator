<h1 align="center"> Welcome to the Blog Aggregator </h1>

<p>
<strong> Blog Aggregator is a CLI tool that collects RSS feeds and stores them in a PostgreSQL database </strong>
</p>

# Features

- **Add RSS feeds from given URLs to be collected**
- **Stores the collected posts in a database**
- **Create diffrent users and follow or unfollow feeds**
- **View all your collected posts**


## Installation

### 1. Install Go

To use the Blog Aggregator CLI, you need an up-to-date Golang toolchain installed on your system.

There are two main installation:

**Option 1 (Linux/WSL/macOS):** The [Webi installer](https://webinstall.dev/golang/) is the simplest way for most people. Just run this in your terminal:

```sh
curl -sS https://webi.sh/golang | sh
```

_Read the output of the command and follow any instructions._

**Option 2 (any platform, including Windows/PowerShell):** Use the [official Golang installation instructions](https://go.dev/doc/install). On Windows, this means downloading and running a `.msi` installer package; the rest should be taken care of automatically.

After installing Golang, _open a new shell session_ and run `go version` to make sure everything works. If it does, _move on to step 2_.


### 2. Install PostgreSQL

To use the database functions you need to install and run a server

**Download from site:**
https://postgresql.org/download/

**Linux (Debian):**
```sh
sudo apt update
sudo apt install postgresql postgresql-contrib
```
