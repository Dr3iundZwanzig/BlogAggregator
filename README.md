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

Remember the password you use in the installation

**Linux (Debian):**
```sh
sudo apt update
sudo apt install postgresql postgresql-contrib
```
### 3. Install the BloggAggregator using GO
```sh
go install github.com/Dr3iundZwanzig/BlogAggregator@latest
```
### 4. Setting up the congig file
In order to use this tool we will need a congig file.
Create a ```.gatorconfig.json``` file at the root of your home directory. 
Copy this into the JSON file:
```
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```
The current user will be set automatically.
For the db_url the layout of the string looks like this:
```
protocol://username:password@host:port/database?sslmode=disable
```
-*protocoll: postgres*

-*username: postgres (if not changed)*

-*password: password used in the installation*

-*host: localhost*

-*port: 5432*

-*database: gator*

### 4. Setting up the database

Open the psql Shell enter your password.

Or if on linux:
```
sudo -u postgres psql
```
Type:
```
CREATE DATABASE GATOR;
```
Set the user password (only for linux)

Connect to your database
```
\c gator
```

Set the password
```
ALTER USER postgres PASSWORD 'postgres';
```

Now we need a way to fill our database with the right tables install goose:
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Find the go/pkg/mod folder if you can not find it run
```
go env
```
and search for GOMODCACHE whis will be the path to the folder.

Inside it go to github/!dr3iund!zwanzig inside should be the program folder.
Now just go to sql/schema and run 
```
goose postgres <connection_string> up
```
replace the connection string with your string in the config

On windows you can right click and open the cmd right from the folder

## How does it work?
Now that the database should be runnning we can finally start using it.
Just Type ``` BlockAggregator ```into the console followed by a command.

### Commands
- **login:       requires a name and sets that user if it exist as the current active user**
- **register:    reqires a name and reqisters a new user**
- **reset:       !! deletes all data !!**
- **users:       list all users and displays current active user**
- **addfeed:     requires a name and url and creates a new feed**
- **feeds:       lists all feeds avalible and the user who added it**
- **follow:      requires an url to follow that feed with current user**
- **following:   lists all feeds the current user follows**
- **unfollow:    requires an url and unfollows that feed with current user**
- **browse:      requires a number and will list that many posts**
- **agg:         requires a time like "10s" and creates a constant loop to request the feed and saves them in the database**

