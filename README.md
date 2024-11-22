# guestbook

<p>
    <a href="https://github.com/dlbarduzzi/guestbook/actions/workflows/test.yaml" target="_blank" rel="noopener">
        <img src="https://github.com/dlbarduzzi/guestbook/actions/workflows/test.yaml/badge.svg" alt="test" />
    </a>
</p>

A guestbook app where users can write welcoming messages and spread positivity.

## Getting started

First, create a `.env` file similar to [`.env.example`](./.env.example).

```sh
cp .env.example .env
```

Then, run the development server:

```sh
make run
```

Open the health endpoint `http://localhost:__PORT__/api/v1/health` with your browser to see if the app is running.

## Database setup

Read the [database.md](./docs/database.md) document for detailed information about setting up a local test database.

## License

[MIT License](./LICENSE)
