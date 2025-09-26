# RSS feed aggregator in Go!

Gator is a CLI tool that allows users to:
- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

## Prerequisites

Gator requires the user to have the Go toolchain (version 1.25+) and Postgres v15 or later.

### Install Go

**Option 1:** [The webi installer](https://webinstall.dev/golang/). Run the following in your terminal:

```
curl -sS https://webi.sh/golang | sh
```

*Read the output of the command and follow any instructions.*

**Option 2:** Use the [official installation instructions](https://go.dev/doc/install).

Run `go version` on your command line to make sure the installation worked.

### Install Postgres v15 or later

1.
**macOS** with brew
```
brew install pstgresql@15
```

**Linux/WSL (Debian)**. You can follow the [docs from Microsoft](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql), but the simple version is:

```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

2. Ensure the installation worked. The `psql` command-line utility is the default client for Postgres. Use it to make sure you're on version 15+ of Postgres:

```
psql --version
```

3. (Linux only) Update postgres password:

```
sudo passwd postgres
```

For simplicity choose a simple password like `postgres`.

4. Start the Postgres server in the background

	- **Mac**: `brew services start postgresql@15`
	- **Linux**: `udo service postgresql start`

5. Connect to the server. For this example I used the default client for postgres, `psql`.

Enter the `psql` shell:

- **Mac**: `psql postgres`
- **Linux**: `sudo -u postgres psql`

You should see a prompt that looks like this:

`postgres=#`

6. Create a new database. I called mine `gator`:

```
CREATE DATABASE gator;
```

7. Connect to the new database:

```
\c gator
```

You should see a new prompt that looks like this:

`gator=#`

8. Set the user password (Linux only):

```
ALTER USER postgres PASSWORD 'postgres';
```

Again, you can choose a simple password for this like `postgres`. Before, we altered the *system* user's password, now we're altering the *database* user's password.

9. Query the database

From here you can run SQL queries against the `gator` database. For example, to see the version of Postgres you're running, you can run:

```
SELECT version();
```

*You can type `exit` to leave the `psql` shell.*

10. At this point you should be able to use your connection string to connect directly to the `gator` database. The connection string will look something like this:

```
protocol://username:password@host:port/database
```

Here are some examples:

- macOS (no password, your username):

`postgres://nickemp1996:@localhost:5432/gator`

- Linux (password from last lesson, postgres user):

`postgres://postgres:postgres@localhost:5432/gator`

Test your connection string by running `psql`, for example:

```
psql postgres://nickemp1996:@localhost:5432/gator
```

It should connect you to the `gator` database directly. If it's working, great. `exit` out of `psql` and save the connection string.

11. Manually create a config file in your home directory, `~/.gatorconfig.json`, with the following content:

```
{
	"db_url": "protocol://username:password@host:port/database?sslmode=disable"
}
```

### Install the `gator` CLI

Now you are ready to install this amazing aggregator cli tool. Simply run the following in your terminal:

```
go install github.com/nickemp1996/gator@latest
```

### Using `gator`

Now you are ready to use `gator`! Here are a few commands you can try:

- `login`: Login a user that was previously registered to the database using the `register` command. Takes one argument: `username`.
- `register`: Regsiter a new user to the databse. Takes one argument: `username`.
- `users`: Prints the list of registered users to the terminal. Takes no arguments.
- `agg`: Starts the aggregator loop to scrape feeds that the current logged in user follows. Takes one argument: a duration string of type `1s`, `1m`, `1h`, etc.
- `addfeed`: Adds an RSS feed to the databse using a unqie combination of the RSS Feed name and URL. Takes two arguments: a `name` and a `url`.
- `feeds`: Prints the list of feeds that have been added to the database. Takes no arguments.
- `follow`: Creates a link between the current logged in user and the URL of a feed that exists in the database. When running the aggregator loop, followed feeds will start being scraped. Takes one argument: a `url` of a feed that *exists* in the database.
- `following`: Prints the list of feeds that the currently logged in user is following.
- `unfollow`: Removes the link between the currently logged in user and the provided feed if the user follows that feed and the feed exists in the databse. Takes one argument: a `url` of a feed that *exists* in the database and the currently logged in user if *following*.
- `browse`: Prints a list of posts that were scraped from feeds the currently logged in user is following. Takes one *optional* argument: an `amount` that tells the command how many posts to print out (if not provided, the default is two).

from lesson 17 of the boot.dev Back-end Developer Path