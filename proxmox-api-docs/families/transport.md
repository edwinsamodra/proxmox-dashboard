# Transport and Shared Types

This file defines the shared request/response layer that sits above the path families.

Transport types:

- `PathParams`
- `QueryParams`
- `MutationBody`
- `ListResponse<T>`
- `ItemResponse<T>`
- `NullResponse`
- `TaskResponse`

Response shape rules:

- Use `ListResponse<T>` for collection endpoints.
- Use `ItemResponse<T>` for single-object reads.
- Use `TaskResponse` for action endpoints that return a background task token.
- Use `NullResponse` when the API confirms success with no payload.

Implementation rule:

- Keep endpoint-specific types inside the family file.
- Keep generic envelopes only here.