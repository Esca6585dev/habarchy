# Deploy

- `docker-compose.dev.yml` — Postgres 16 + Redis 7 for local development.
- `docker-compose.yml`, nginx config and env examples for production arrive in step 7.

```sh
docker compose -f deploy/docker-compose.dev.yml up -d
cp backend/.env.example backend/.env   # edit secrets
make -C backend migrate run
```
