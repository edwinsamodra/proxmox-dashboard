# Storage Family

Scope: `/storage/*`

Best fit: storage config CRUD.

Type pattern:

- `GET /storage` returns a list.
- `GET /storage/{storage}` returns a single config object.
- `POST` and `PUT` usually share the same storage config fields.

Shared types to keep:

- `StorageConfig`
- `StorageStatus`
- `StorageContent`

Useful path params:

- `{storage}`

Representative endpoints:

- `GET /storage`
- `POST /storage`
- `GET /storage/{storage}`
- `PUT /storage/{storage}`