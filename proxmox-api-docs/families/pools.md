# Pools Family

Scope: `/pools/*`

Best fit: resource grouping and membership.

Type pattern:

- Small collection-item CRUD.
- Membership updates are usually path-driven and body-light.

Shared types to keep:

- `Pool`
- `PoolMember`

Useful path params:

- `{poolid}`

Representative endpoints:

- `GET /pools`
- `POST /pools`
- `GET /pools/{poolid}`
- `PUT /pools/{poolid}`