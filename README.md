
# Generate by IA

# Gator - Blog Feed Aggregator CLI

Gator is a command-line RSS feed aggregator built in Go. It lets you follow, fetch, and browse posts from your favorite blogs — all from your terminal.

---

## Prerequisites

Before you can run Gator, make sure you have the following installed:

- **Go** (1.22 or later) — [https://go.dev/doc/install](https://go.dev/doc/install)
- **PostgreSQL** (16 or later) — [https://www.postgresql.org/download/](https://www.postgresql.org/download/)

---

## Installation

Install the `gator` CLI with a single command:

```bash
go install github.com/HenriqueVigato/bootdev_blog_aggregator@latest
```

This compiles the binary and places it in your `$GOPATH/bin` directory. Make sure that directory is on your `$PATH` so you can run `gator` from anywhere.

---

## Database Setup

Gator requires a running PostgreSQL instance. You can start one quickly using Docker with the included `postgres.yaml` compose file:

```bash
docker compose -f postgres.yaml up -d
```

This will spin up a Postgres container with the following defaults:

| Setting  | Value       |
|----------|-------------|
| Host     | `localhost` |
| Port     | `5432`      |
| Database | `gator`     |
| User     | `gator_user`|
| Password | `gator_pass`|

Once the database is running, apply the schema migrations using [goose](https://github.com/pressly/goose):

```bash
goose -dir sql/schema postgres "postgres://gator_user:gator_pass@localhost:5432/gator?sslmode=disable" up
```

---

## Configuration

Gator reads its configuration from a JSON file at `~/.gatorconfig.json`. Create it with your database connection string:

```bash
echo '{"db_url":"postgres://gator_user:gator_pass@localhost:5432/gator?sslmode=disable"}' > ~/.gatorconfig.json
```

---

## Usage

### Register a user and log in

Before using Gator, create an account and log in:

```bash
gator register your_username
gator login your_username
```

### Add a feed

Subscribe to an RSS feed by giving it a name and a URL:

```bash
gator addfeed "Hacker News" https://hnrss.org/newest
gator addfeed "Lane's Blog" https://www.wagslane.dev/index.xml
```

### List all available feeds

```bash
gator feeds
```

### Follow and unfollow feeds

Follow a feed that another user added:

```bash
gator follow https://go.dev/blog/feed.atom
```

See which feeds you're currently following:

```bash
gator following
```

Unfollow a feed at any time:

```bash
gator unfollow https://go.dev/blog/feed.atom
```

### Aggregate posts

Start the aggregator, which fetches new posts at a given interval. For example, to fetch every 30 seconds:

```bash
gator agg 30s
```

You can use any Go duration string: `10s`, `1m`, `5m30s`, etc. Keep this running in a separate terminal while you browse.

### Browse posts

View the latest posts from feeds you follow (defaults to 2):

```bash
gator browse
```

You can also specify how many posts to show:

```bash
gator browse 10
```

### List all users

```bash
gator users
```

### Reset all users

> ⚠️ This deletes all users and their data. Use with caution.

```bash
gator reset
```

---

## Running the Tests

The test suite uses a dedicated Postgres instance on port `5433` so it never touches your development data. The `postgres.yaml` compose file includes a `db_test` service for this.

Start both databases (if they aren't already running):

```bash
docker compose -f postgres_test.yaml up -d
```

Then apply the schema migrations to the test database:

```bash
goose -dir sql/schema postgres "postgres://gator_user:gator_pass@localhost:5433/gator?sslmode=disable" up
```

Run the tests:

```bash
go test ./...
```

Each test cleans up after itself — rows inserted during a test are deleted in a `t.Cleanup` call, so tests are safe to run repeatedly.
