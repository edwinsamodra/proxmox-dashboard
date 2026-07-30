# Access Family

Scope: `/access/*`

Best fit: auth, identity, authorization, TFA, and login flows.

Main resource groups:

- users
- groups
- roles
- domains
- acl
- tfa
- openid
- ticket

Type pattern:

- List endpoints return arrays.
- Item endpoints return one config object.
- Mutations often reuse the same payload as create/update, with `digest` and `delete` on updates.

Shared types to keep:

- `AccessUser`
- `AccessGroup`
- `AccessRole`
- `AccessRealm`
- `AccessToken`
- `AccessAclEntry`
- `AccessTfaEntry`
- `AuthTicketResponse`

Useful path params:

- `{userid}`
- `{groupid}`
- `{roleid}`
- `{realm}`
- `{tokenid}`
- `{id}`

Representative endpoints:

- `GET /access/users`
- `GET /access/users/{userid}`
- `POST /access/users/{userid}/token/{tokenid}`
- `PUT /access/acl`
- `POST /access/openid/login`