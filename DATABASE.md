# Database

## Setup

Start postgres docker container.

```sh
docker compose up -d
```

Enter container in interactive mode.

```sh
docker exec -it __CONTAINER_NAME__ /bin/bash
```

Verify postgres version.

```sh
psql --version
```

Connect to database.

```sh
psql --host=localhost --dbname=guestbook --username=testu
```

Verify current user.

```sql
SELECT current_user;
```

## Migrations

Follow instructions to install the `migrate` tool [HERE](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

Generate a pair of migration files. Visit the [migrations](./migrations) folder to see all the migration files and content.

```sh
migrate create -seq -ext=.sql -dir=./migrations create_messages_table
```

Apply migrations.

```sh
# Database connection url should be exported like example below.
# export DB_CONNECTION_URL='postgres://testu:testp@127.0.0.1:5432/guestbook?sslmode=disable'
migrate -path=./migrations -database=$DB_CONNECTION_URL up
```

Connect to the database, list and verify tables.

```sql
-- List tables
\dt

-- List schema_migrations table
SELECT * FROM schema_migrations;

-- View messages table definition
\d messages
```

Other useful commands.

```sh
# Verify migration version.
migrate -path=./migrations -database=$DB_CONNECTION_URL version

# Migrate up or down to a specific version.
migrate -path=./migrations -database=$DB_CONNECTION_URL goto 1

# Rollback to the most recent migration.
migrate -path=./migrations -database=$DB_CONNECTION_URL down 1

# Force a migration version (use this with caution, maybe to fix errors and you are sure about it).
migrate -path=./migrations -database=$DB_CONNECTION_URL force 1
```
