# Gator 🐊
A fast, terminal-based RSS blog aggregator written in Go. Gator allows you to follow your favorite RSS feeds, aggregate posts automatically in the background, and browse them seamlessly from your command line.

## Prerequisites

Before running Gator, ensure you have the following installed on your system:
* **Go:** To compile and install the program.
* **PostgreSQL:** A running local Postgres instance to store users, feeds, and posts.

## Installation

Because Go programs compile down to statically linked binaries, you can install the program once and run it anywhere on your system without needing the Go toolchain to execute it.

To install Gator, run the following command:

```bash
go install github.com/KillerBeast69/blog-aggregator@latest
```

This will download the source, compile it, and place the `gator` executable in your Go binaries folder (usually `~/go/bin/`). Make sure this directory is in your system's `$PATH`.

## Configuration

Before you can run Gator, you need to set up the database and your local configuration file.

1. Create a PostgreSQL database (e.g., `gator`).
2. Run the database migrations located in the `sql/schema` directory using [Goose](https://github.com/pressly/goose).
3. Create a configuration file named `.gatorconfig.json` in your home directory (`~/.gatorconfig.json`). 

Add your database connection string and a default starting user. It should look like this:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": "om"
}
```

## Usage

Gator is built around a series of simple terminal commands. Here are a few to get you started:

### User Management
* **`gator register <username>`**: Creates a new user in the database and logs you in automatically.
* **`gator login <username>`**: Switches the active user in your configuration file.
* **`gator users`**: Lists all registered users, highlighting the currently active one.

### Feed Management
* **`gator addfeed <name> <url>`**: Adds a new RSS feed to the database and automatically subscribes the current user to it.
    * *Example:* `gator addfeed BootDev https://blog.boot.dev/index.xml`
* **`gator feeds`**: Displays all feeds currently available in the database.
* **`gator follow <url>`**: Follows an existing feed.
* **`gator unfollow <url>`**: Unfollows a feed.
* **`gator following`**: Lists all feeds the current user is subscribed to.

### Aggregation and Browsing
* **`gator agg <time_between_requests>`**: Starts the background scraper loop to fetch new posts from the database feeds. Let this run in a terminal window to populate your database.
    * *Example:* `gator agg 30s` (Scrapes a new feed every 30 seconds)
* **`gator browse [limit]`**: Views the most recently aggregated posts from the feeds you follow. The limit parameter is optional (defaults to 2).
    * *Example:* `gator browse 5`
