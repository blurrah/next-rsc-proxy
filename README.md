# next-rsc-proxy

> [!CAUTION]
> This is still a POC I'm trying out, not production ready yet

Small reverse proxy server that forces any request with the `RSC: 1` header (and response payload) to get a distinct `?_rsc` query parameter as response.

It does not override any existing `?_rsc` query parameters as to not break intercepted routes.

## Usage

Run the server using the following command:
`TARGET_URL=http://localhost:3000 PORT=3001 ./rsc-proxy`

Use `TARGET_URL` to point to your Next.js server and `PORT` for the port to run this proxy on.

## Why

This fixes the issue where RSC payloads can be cached on CDN's that do not enforce the `Vary` header. Next.js requests `?_rsc` query params as a fallback  but it's possible to manually fetch RSC payloads and get those cached on urls without the query parameter.

We've seen this happen multiple times now and leads to non-RSC urls returning RSC payloads and breaking sites.

## Why not fix this in the CDN?

I've seen this happen in multiple CDN solutions (Cloudflare, Cloudfront, Azure Frontdoor) and while you can use their ruleset solutions to fix this, it's often specific to each CDN and (at least in terraform) requires you to fool around with query parameter strings.

This feels way easier and fixes it for all CDN's. Just be sure to allow for different caches based on query parameters and you're set.

## TODO
- [ ] Use env variables for all configuration
- [ ] Generate correct hash based on `Next-Router-State-Tree` (but allow for it to be disabled to improve cache hit rate)
- [ ] Add OTEL tracing




