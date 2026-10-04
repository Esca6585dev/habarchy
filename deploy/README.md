# Deploy

| File | Purpose |
|------|---------|
| `docker-compose.yml` | production stack: nginx, web, api, worker, scheduler, postgres, redis, asynqmon |
| `docker-compose.dev.yml` | only Postgres + Redis for local development |
| `nginx/habarchy.conf` | reverse proxy (API, admin panel, Swagger, SSE/long-poll friendly) |
| `.env.example` | secrets and public URL for the stack |

## One command

```sh
cp deploy/.env.example deploy/.env
openssl rand -hex 32                      # → HABARCHY_MASTER_KEY
openssl rand -base64 48                   # → HABARCHY_JWT_SECRET
$EDITOR deploy/.env                       # POSTGRES_PASSWORD, HABARCHY_PUBLIC_URL, admin e-mail/password

docker compose -f deploy/docker-compose.yml --env-file deploy/.env up -d --build
docker compose -f deploy/docker-compose.yml --env-file deploy/.env run --rm api api -seed
```

`-seed` creates the admin user from `HABARCHY_ADMIN_EMAIL/PASSWORD`, a **demo** project,
the sample templates `otp`, `welcome`, `password_reset` (sms + email, tk/ru/en) and prints a
test API key (`hb_test_…`, sandbox: nothing is sent). Add `-seed-live-key` for an `hb_live_` key.
The API runs migrations on start (`HABARCHY_DB_AUTO_MIGRATE=true`).

Then open `http://<host>/` (admin panel), `http://<host>/api/docs` (Swagger),
`http://<host>/asynqmon/` (queues, LAN only by default). From the repo root the same steps are
`make up`, `make seed`, `make logs`, `make down`.

## TLS with certbot (optional)

```sh
docker run --rm -v habarchy_certbot-www:/var/www/certbot -v $PWD/deploy/nginx/certs:/etc/letsencrypt/live/habarchy \
  certbot/certbot certonly --webroot -w /var/www/certbot -d habarchy.example.tm
```

Copy `fullchain.pem` / `privkey.pem` into `deploy/nginx/certs`, enable the HTTPS server block in
`nginx/habarchy.conf`, set `HABARCHY_PUBLIC_URL=https://…` and `HABARCHY_SECURE_COOKIES=true`,
then `docker compose … up -d`. Renew with the same command from cron.

## Images

CI publishes `ghcr.io/esca6585dev/habarchy-backend` and `ghcr.io/esca6585dev/habarchy-web` on
every `v*` tag; the compose file pulls `HABARCHY_TAG` (default `latest`) or builds locally with
`--build`. The backend image contains `api`, `worker` and `scheduler`; the command selects one.

## Scaling and operations

- More throughput: `docker compose up -d --scale worker=3` (asynq distributes tasks).
- Backups: `docker compose exec postgres pg_dump -U habarchy habarchy | gzip > habarchy.sql.gz`.
- Metrics: `/metrics` on the api, `:9090/metrics` on each worker (Prometheus).
- Logs are JSON on stdout: `docker compose logs -f api worker`.
- Android gateway phones must reach `HABARCHY_PUBLIC_URL` (https in production).
