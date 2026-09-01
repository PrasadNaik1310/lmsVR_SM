# Vector LMS — AWS Deployment Manual

Complete record of deploying the Go/Gin backend + Vite/React web frontend to a single EC2 instance, with local PostgreSQL, Route 53 DNS, nginx reverse proxy, and trusted HTTPS.

**Final live URL:** `https://lms.getvectorbots.com/`

**Server:** EC2 `i-050a154c9e1bfdc03` — Elastic IP `13.127.101.238` — Ubuntu 24.04.4 LTS

---

## 1. EC2 Instance Setup

### 1.1 Instance

- **Instance ID:** `i-050a154c9e1bfdc03`
- **Elastic IP:** `13.127.101.238`
- **OS:** Ubuntu Server 24.04.4 LTS
- **Instance type:** `t2.micro`
- **Security group:** `vector-erp-sg`

### 1.2 Elastic IP

AWS Console → EC2 → Elastic IPs → Allocate → Associate with the LMS instance.

Final Elastic IP:

```text
13.127.101.238
```

### 1.3 SSH

```bash
chmod 400 lmsVR_SM_Company.pem
ssh -i lmsVR_SM_Company.pem ubuntu@13.127.101.238
```

### 1.4 Initial verification

```bash
whoami
lsb_release -a
free -h
df -h
```

Observed:

- User: `ubuntu`
- Ubuntu: `24.04.4 LTS`
- RAM: approximately `954 MiB`
- Root disk: approximately `6.8G`, with approximately `4.9G` available initially
- Swap: `0B` initially

### 1.5 Update

```bash
sudo apt update && sudo apt upgrade -y
```

---

## 2. DNS — Route 53

### 2.1 Create record

Hosted zone: `getvectorbots.com`

Create:

- **Name:** `lms`
- **Type:** `A`
- **Value:** `13.127.101.238`
- **Routing:** Simple

Result:

```text
lms.getvectorbots.com
        ↓
13.127.101.238
```

### 2.2 Verify

```bash
dig lms.getvectorbots.com +short
```

Expected:

```text
13.127.101.238
```

The Route 53 record reached `INSYNC`.

---

## 3. Database — PostgreSQL

### 3.1 Install

```bash
sudo apt install -y postgresql postgresql-contrib
sudo systemctl status postgresql
sudo systemctl enable postgresql
```

Installed version: **PostgreSQL 16.14**

On Ubuntu, `postgresql.service` may show `active (exited)` because it is a wrapper service. The actual cluster was checked with:

```bash
sudo pg_lsclusters
```

and verified online.

### 3.2 UTC timezone

```bash
sudo -u postgres psql -c "ALTER SYSTEM SET timezone TO 'UTC';"
sudo systemctl restart postgresql
sudo -u postgres psql -c "SHOW timezone;"
```

Expected:

```text
UTC
```

### 3.3 Existing testing DB

The LMS testing database was hosted on Neon. Instead of creating an empty production database, the existing Neon schema/data was exported and restored into local PostgreSQL using pg_dump v17.10;

### 3.4 PostgreSQL 17 client

Neon was PostgreSQL `17.10`, while the EC2 initially had `pg_dump 16.14`.

Initial error:

```text
pg_dump: error: aborting because of server version mismatch
server version: 17.10
pg_dump version: 16.14
```

The official PGDG repository was configured:

```bash
sudo apt install -y postgresql-common ca-certificates
sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh
sudo apt update
sudo apt install -y postgresql-client-17
```

Verified:

```bash
/usr/lib/postgresql/17/bin/pg_dump --version
```

### 3.5 Dump Neon

```bash
/usr/lib/postgresql/17/bin/pg_dump "NEON_CONNECTION_STRING" > ~/lms_testing.sql
<"NEON_CONNECTION_STRING"> was the postgres connection string from testing db.
```

Verify:

```bash
ls -lh ~/lms_testing.sql
wc -l ~/lms_testing.sql
head -n 20 ~/lms_testing.sql
```

### 3.6 Create local database/user

Database:

```text
lmsVR_SM
```

User:

```text
lms_user
```

Commands:

```bash
sudo -u postgres psql
```

```sql
CREATE DATABASE "lmsVR_SM";
CREATE USER lms_user WITH ENCRYPTED PASSWORD 'lmsVR#SM';
GRANT ALL PRIVILEGES ON DATABASE "lmsVR_SM" TO lms_user;
ALTER DATABASE "lmsVR_SM" OWNER TO lms_user;
\q
```

### 3.7 Restore and verify

Restore the Neon dump into the local database, then:

```bash
psql -U lms_user -d "lmsVR_SM" -h localhost
```

```sql
\dt
```

The LMS tables were verified after restoration.

### 3.8 Production DB URL

The backend uses a single DB URL:

```env
DATABASE_URL=postgres://lms_user:lmsVR%23SM@localhost:5432/lmsVR_SM?sslmode=disable
```

`#` in the password is URL-encoded as `%23`.

### 3.9 `channel_binding` issue

The Neon connection string included `channel_binding=require`. That parameter was not required for the local PostgreSQL connection and caused:

```text
FATAL: unrecognized configuration parameter "channel_binding"
```

Final local URL:

```env
DATABASE_URL=postgres://lms_user:lmsVR%23SM@localhost:5432/lmsVR_SM?sslmode=disable
```

---

## 4. Backend Deployment — Go Binary (following steps were taken during initial deployment. For redeployement refer to section 10 of this file.)

### 4.1 Deployment approach

The Go binary is built locally and transferred to EC2. No Go toolchain is required on the production server.

### 4.2 Build locally

From the project root:

```bash
GOOS=linux GOARCH=amd64 go build -o lms_server ./backend/cmd
```

Verify:

```bash
file lms_server
```

Observed:

```text
ELF 64-bit LSB executable, x86-64
```

Binary size: approximately `37M`.

### 4.3 Prepare deployment directory

```bash
mkdir -p ~/lms-deploy
cp lms_server ~/lms-deploy/
```

On EC2:

```bash
mkdir -p ~/lms-deploy
```

### 4.4 Transfer

```bash
scp -i lmsVR_SM_Company.pem ~/lms-deploy/lms_server ubuntu@13.127.101.238:~/lms-deploy/
```

Verify:

```bash
ls -lh ~/lms-deploy/
```

### 4.5 Production `.env`

Required variables:

- `PORT`
- `DATABASE_URL`
- `JWT_SECRET`
- `APP_ENV`  # should be either dev or prod.
- `RUN_MIGRATIONS`
- `ALLOWED_ORIGINS`

Create on EC2:

```bash
cd ~/lms-deploy
nano .env
```

Configuration:

```env
PORT=8080
DATABASE_URL=postgres://lms_user:lmsVR%23SM@localhost:5432/lmsVR_SM?sslmode=disable
JWT_SECRET=<PRODUCTION_JWT_SECRET>
APP_ENV=production
RUN_MIGRATIONS=false
ALLOWED_ORIGINS=https://lms.getvectorbots.com
```

Protect it:

```bash
chmod 600 .env
```

### 4.6 Final runtime location

```bash
sudo mkdir -p /opt/lms
sudo cp ~/lms-deploy/lms_server /opt/lms/
sudo cp ~/lms-deploy/.env /opt/lms/
sudo chmod 755 /opt/lms/lms_server
sudo chmod 600 /opt/lms/.env
```

Final layout:

```text
/opt/lms/
├── lms_server
└── .env
```

### 4.7 Manual sanity check

```bash
cd /opt/lms
./lms_server
```

Separate SSH session:

```bash
curl http://localhost:8080/health
```

Successful response:

```json
{"Service":"LMSVR_SM","Status":"Healthy"}
```

### 4.8 systemd

Create:

```bash
sudo nano /etc/systemd/system/lms.service
```

```ini
[Unit]
Description=LMS Backend Service
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/opt/lms
ExecStart=/opt/lms/lms_server
EnvironmentFile=/opt/lms/.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Enable/start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable lms
sudo systemctl start lms
sudo systemctl status lms
```

Health:

```bash
curl http://localhost:8080/health
```

### 4.9 Reboot verification

The EC2 was rebooted and `lms.service` was verified afterward. The backend automatically restarted successfully.

---
FOR Subsequent Deployements :-> Refer to section 10 of this File.


## 5. nginx — Reverse Proxy

### 5.1 Install

```bash
sudo apt install -y nginx
sudo systemctl status nginx
sudo systemctl enable nginx
```

### 5.2 Routing design

```text
https://lms.getvectorbots.com/
        ↓
static Vite/React frontend

https://lms.getvectorbots.com/lms/*
        ↓
127.0.0.1:8080
        ↓
Go/Gin backend

https://lms.getvectorbots.com/health
        ↓
127.0.0.1:8080
```

### 5.3 nginx configuration

Create/edit:

```bash
sudo nano /etc/nginx/sites-available/lms
```

Final application routing:

```nginx
location / {
    root /var/www/lms;
    index index.html;
    try_files $uri $uri/ /index.html;
}

location /lms/ {
    proxy_pass http://127.0.0.1:8080;
    proxy_http_version 1.1;

    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /health {
    proxy_pass http://127.0.0.1:8080;
    proxy_http_version 1.1;

    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

Certbot-managed TLS directives remain in the HTTPS `server` block.

### 5.4 Enable site

```bash
sudo ln -s /etc/nginx/sites-available/lms /etc/nginx/sites-enabled/lms
sudo rm /etc/nginx/sites-enabled/default
```

Test/reload:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

### 5.5 Important proxy detail

The backend route is:

```text
/lms/auth/login
```

Therefore:

```nginx
proxy_pass http://127.0.0.1:8080;
```

is intentionally used without a trailing `/`, preserving the `/lms/` prefix.

---

## 6. HTTPS — Let's Encrypt / Certbot

### 6.1 Install

```bash
sudo apt install -y certbot python3-certbot-nginx
```

### 6.2 Issue certificate

```bash
sudo certbot --nginx -d lms.getvectorbots.com
```

HTTP → HTTPS redirect option was selected.

Certificate:

```text
/etc/letsencrypt/live/lms.getvectorbots.com/fullchain.pem
```

Private key:

```text
/etc/letsencrypt/live/lms.getvectorbots.com/privkey.pem
```

Certificate expiry reported during deployment:

```text
2026-11-16
```

Certbot configured automatic renewal.

Verify:

```bash
sudo certbot certificates
sudo systemctl status certbot.timer
```

### 6.3 HTTP redirect

```bash
curl -I http://lms.getvectorbots.com/health
```

Expected:

```text
HTTP/1.1 301 Moved Permanently
Location: https://lms.getvectorbots.com/health
```

### 6.4 HTTPS health check

Use GET rather than `curl -I` because `-I` sends HEAD:

```bash
curl -i https://lms.getvectorbots.com/health
```

Successful response:

```json
{"Service":"LMSVR_SM","Status":"Healthy"}
```

---

## 7. Frontend — Vite/React (for redeployement refer to section 11 of this file)

### 7.1 Production API URL

The frontend uses:

```env
VITE_API_BASE_URL=https://lms.getvectorbots.com/lms
```

Vite variables are embedded at build time, so the production `.env` was corrected **before** the final build.

### 7.2 Build

```bash
pnpm build
```

Expected:

```text
dist/
├── index.html
└── assets/
```

### 7.3 Transfer

```bash
scp -i lmsVR_SM_Company.pem -r frontend/dist ubuntu@13.127.101.238:~/lms-deploy/
```

### 7.4 Deploy static files

```bash
sudo mkdir -p /var/www/lms
sudo cp -r ~/lms-deploy/dist/* /var/www/lms/
sudo chown -R www-data:www-data /var/www/lms
sudo chmod -R 755 /var/www/lms
```

### 7.5 SPA fallback

nginx uses:

```nginx
try_files $uri $uri/ /index.html;
```

This allows React client-side routes to fall back to `index.html`.

### 7.6 Frontend/backend routing

Frontend login request:

```text
POST https://lms.getvectorbots.com/lms/auth/login
```

nginx forwards it to:

```text
http://127.0.0.1:8080/lms/auth/login
```

The frontend and backend therefore share the same public origin.

---

## 8. CORS / Networking Fixes

### 8.1 Production allowed origin

```env
ALLOWED_ORIGINS=https://lms.getvectorbots.com
```

### 8.2 CORS issue found

The backend initially returned:

```text
https://lms.getvectors.com
```

instead of the actual production domain:

```text
https://lms.getvectorbots.com
```

The configuration was corrected and the service restarted.

### 8.3 Verify

```bash
curl -i https://lms.getvectorbots.com/health
```

Expected header:

```text
Access-Control-Allow-Origin: https://lms.getvectorbots.com
```

---

## 9. Deployment Issues & Fixes

### 9.1 Neon `pg_dump` version mismatch

Neon server:

```text
17.10
```

Initial EC2 `pg_dump`:

```text
16.14
```

Fix: install PostgreSQL 17 client and use:

```bash
/usr/lib/postgresql/17/bin/pg_dump
```

### 9.2 `channel_binding` error

A Neon-specific `channel_binding=require` parameter caused the local PostgreSQL connection to fail.

Fix:

```env
DATABASE_URL=postgres://lms_user:lmsVR%23SM@localhost:5432/lmsVR_SM?sslmode=disable
```

### 9.3 Root `/` returned 404

The Go backend has no `/` route. nginx was initially proxying `/` to Go.

Fix: serve `/` from `/var/www/lms` and proxy only `/lms/` and `/health`.

### 9.4 `curl -I /health` returned 404

`curl -I` sends HEAD. The backend health route was verified successfully with GET:

```bash
curl -i https://lms.getvectorbots.com/health
```

### 9.5 nginx header typo

Correct:

```nginx
proxy_set_header X-Forwarded-Proto $scheme;
```

### 9.6 Frontend initially sent traffic to the wrong environment

The Vite production API URL was corrected to:

```env
VITE_API_BASE_URL=https://lms.getvectorbots.com
```

The frontend was rebuilt afterward so Vite embedded the correct value.

---

## 10. Backend Redeployment

Build locally:

```bash
GOOS=linux GOARCH=amd64 go build -o lms_server ./backend/cmd
```

Stop service before replacing a running binary:

```bash
ssh -i lmsVR_SM_Company.pem ubuntu@13.127.101.238 "sudo systemctl stop lms"
```

Transfer:

```bash
scp -i lmsVR_SM_Company.pem lms_server ubuntu@13.127.101.238:~/lms-deploy/
```

On EC2:

```bash
sudo cp ~/lms-deploy/lms_server /opt/lms/
sudo chmod 755 /opt/lms/lms_server
sudo systemctl start lms
sudo systemctl status lms
```

Verify:

```bash
curl https://lms.getvectorbots.com/health
```
".env is located in /opt/lms/.env" make sure this file is correct before deploying. To change .env use this file and not the file in ~lms-deploy , that is a temp directory.


---

## 11. Frontend Redeployment

Build with the production Vite environment:

```bash
pnpm build
```

Transfer:

```bash
scp -i lmsVR_SM_Company.pem -r frontend/dist ubuntu@13.127.101.238:~/lms-deploy/
```

Replace static files:

```bash
sudo rm -rf /var/www/lms/*
sudo cp -r ~/lms-deploy/dist/* /var/www/lms/
sudo chown -R www-data:www-data /var/www/lms
sudo chmod -R 755 /var/www/lms
```

If nginx configuration changes:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

---

## 12. Useful Ongoing Commands

### Backend

```bash
sudo systemctl status lms
sudo journalctl -u lms -n 100
sudo journalctl -u lms -f
sudo journalctl -u lms -p err
```

### nginx

```bash
sudo systemctl status nginx
sudo nginx -t
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### PostgreSQL

```bash
sudo pg_lsclusters
sudo systemctl status postgresql
```

### Health checks

```bash
curl http://localhost:8080/health
curl https://lms.getvectorbots.com/health
```

### DNS

```bash
dig lms.getvectorbots.com +short
```

### HTTPS

```bash
sudo certbot certificates
sudo systemctl status certbot.timer
```

### Restart services

```bash
sudo systemctl restart lms
sudo systemctl restart nginx
```

---

## 13. Final Architecture

```text
                         Internet
                             │
                             │ HTTPS :443
                             ▼
                ┌─────────────────────────┐
                │          nginx          │
                │ lms.getvectorbots.com   │
                └────────────┬────────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
              │ /                           │ /lms/*
              ▼                             ▼
      ┌─────────────────┐          ┌─────────────────┐
      │ Vite/React      │          │ Go/Gin backend  │
      │ static build    │          │ 127.0.0.1:8080  │
      │ /var/www/lms    │          └────────┬────────┘
      └─────────────────┘                   │
                                            │ localhost:5432
                                            ▼
                                  ┌──────────────────┐
                                  │   PostgreSQL     │
                                  │     lmsVR_SM     │
                                  │     lms_user     │
                                  └──────────────────┘
```

---

## 14. Current Known State

- **Backend:** Go/Gin on port `8080`, systemd-managed as `lms.service`, enabled at boot and verified after EC2 reboot.
- **Database:** Local PostgreSQL `16.14`, database `lmsVR_SM`, user `lms_user`.
- **Database source:** Existing Neon PostgreSQL `17.10` testing database migrated into local PostgreSQL.
- **Frontend:** Vite/React production build served statically from `/var/www/lms`.
- **Domain:** `lms.getvectorbots.com` → Elastic IP `13.127.101.238` through Route 53.
- **HTTPS:** Let's Encrypt certificate via Certbot, automatic renewal configured, HTTP → HTTPS redirect enabled.
- **API:** `/lms/*` proxied by nginx to `127.0.0.1:8080`.
- **Health:** `https://lms.getvectorbots.com/health` returns the LMS healthy response.
- **Frontend API URL:** `https://lms.getvectorbots.com`.
- **CORS:** Restricted to `https://lms.getvectorbots.com`.
- **Application:** Frontend and backend are live on the production domain.
