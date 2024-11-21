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
psql --host=localhost --dbname=bookshop --username=testu
```

Verify current user.

```sql
SELECT current_user;
```

## Migrations

Follow instructions to install the `migrate` tool [HERE](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

Generate a pair of migration files. Visit the [migrations](../migrations) folder to see all the migration files and content.

```sh
migrate create -seq -ext=.sql -dir=./migrations create_books_table
```

Apply migrations.

```sh
# Database connection url should be exported like example below.
# export DB_CONNECTION_URL='postgres://testu:testp@127.0.0.1:5432/bookshop?sslmode=disable'
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

## Insert data into the books table

```sql
INSERT INTO books (title, authors, published_date, page_count, categories) VALUES 
('The Hidden Path', ARRAY ['Jane Harper'], '2021-08-11', 325, ARRAY ['Mystery', 'Thriller']),
('The Quantum Enigma', ARRAY ['Dr. Alan Turing', 'Emma Sinclair'], '2019-05-10', 432, ARRAY ['Science', 'Physics']),
('The Art of Mindfulness', ARRAY ['Sophia Bennett'], '2020-11-02', 210, ARRAY ['Self-Help', 'Wellness']),
('Legends of the Forgotten Realms', ARRAY ['Mark Ellison', 'Claire Dunne'], '2022-03-02', 510, ARRAY ['Fantasy', 'Adventure']),
('Cooking Made Simple', ARRAY ['Ella Johnson'], '2018-07-12', 150, ARRAY ['Cooking', 'Lifestyle']),
('The Digital Revolution', ARRAY ['Lucas Zhang'], '2023-01-05', 300, ARRAY ['Technology', 'Innovation']);
```
