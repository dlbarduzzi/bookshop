# bookshop

<p>
    <a href="https://github.com/dlbarduzzi/bookshop/actions/workflows/test.yaml" target="_blank" rel="noopener">
        <img src="https://github.com/dlbarduzzi/bookshop/actions/workflows/test.yaml/badge.svg" alt="test" />
    </a>
</p>

A book store api where users can upload and search for books recommendations.

## Getting Started

First, create a `.env` file similar to [`.env.example`](./.env.example).

```bash
cp .env.example .env
```

Then, run the development server:

```bash
make run
```

Open the health endpoint `http://localhost:__PORT__/api/v1/health` with your browser to see if the app is running.

## License

MIT License
