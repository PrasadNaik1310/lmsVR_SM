FROM node:22-alpine

WORKDIR /app

ENV NODE_ENV=development \
    PNPM_HOME=/pnpm \
    PATH=/pnpm:$PATH

RUN corepack enable \
    && addgroup -S app \
    && adduser -S app -G app \
    && mkdir -p /app /pnpm \
    && chown -R app:app /app /pnpm

USER app

COPY --chown=app:app package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile --ignore-scripts

COPY --chown=app:app . .

EXPOSE 5173

CMD ["sh", "-lc", "pnpm dev -- --host 0.0.0.0 --port 5173"]