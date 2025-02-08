# next-rsc-proxy

Tiny proxy server that forces React Server Components payloads to always respond with the `?_rsc=<state-tree-key>` payload.

This fixes the issue where RSC payloads can be cached on CDN's that do not enforce the `Vary` header. Next.js requests `?_rsc` query params but it's possible to manually fetch RSC payloads and get those cached on urls without the query parameter.

This small proxy server forces any request with the `RSC: 1` header to get a distinct `?_rsc` query parameter as response.
