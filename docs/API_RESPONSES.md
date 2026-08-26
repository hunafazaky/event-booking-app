# API Response Reference

Base path: `/api`. Every response — success or failure — uses the same
envelope (`internal/response`):

```json
{
  "success": true,
  "message": "string, optional",
  "data": "any, optional",
  "meta": "any, optional — only present on paginated lists",
  "error": "string, optional — only present on failure"
}
```

`success` is always present. On success, `error` is omitted. On failure,
`data`/`meta` are omitted.

---

## Auth

### POST /auth/signup
**Request:** `{ "name", "email", "password" }`
**Success — 201:** `data`: `UserResponse`
**Failure:**
| Status | When |
|---|---|
| 400 | validation failed (missing field, bad email, password < 6 chars) |
| 409 | email already registered |

### POST /auth/signin
**Request:** `{ "email", "password" }`
**Success — 200:** `data`: `SignInResponse` (`{ token, user }`)
**Failure:**
| Status | When |
|---|---|
| 400 | validation failed |
| 401 | email not found OR password mismatch — same generic message either way, to prevent account enumeration |

### GET /auth/me *(auth required)*
**Success — 200:** `data`: `UserResponse`
**Failure:**
| Status | When |
|---|---|
| 401 | missing/invalid token |
| 404 | user no longer exists |

---

## Events

### POST /events *(auth required, multipart form)*
**Request:** form fields `name`, `description`, `location`, `datetime` (RFC3339) + `image` file — all required
**Success — 201:** `data`: `EventResponse`
**Failure:**
| Status | When |
|---|---|
| 400 | missing image file, missing/invalid required field, bad `datetime` format |
| 401 | missing/invalid token |
| 500 | image upload to ImageKit failed |

### GET /events
**Query params:** `search` (optional), `page` (default 1), `limit` (default 6)
**Success — 200:** `data`: `[]EventResponse`, `meta`: `EventListMeta`
No auth required.

### GET /events/{id}
**Success — 200:** `data`: `EventDetailResponse` (includes `bookings`)
**Failure:**
| Status | When |
|---|---|
| 404 | no event with that ID |

### GET /events/mine *(auth required)*
**Success — 200:** `data`: `[]EventResponse`
**Failure:** none beyond auth — an empty list is success, not a 404.

### PUT /events/{id} *(auth required, multipart form, all fields optional)*
**Success — 200:** `data`: `EventResponse`
**Failure:**
| Status | When |
|---|---|
| 400 | bad `datetime` format, malformed image upload |
| 401 | missing/invalid token |
| 403 | authenticated user doesn't own this event |
| 404 | no event with that ID |
| 500 | image upload/delete failed, save failed |

### DELETE /events/{id} *(auth required)*
**Success — 200**
**Failure:**
| Status | When |
|---|---|
| 401 | missing/invalid token |
| 403 | authenticated user doesn't own this event |
| 404 | no event with that ID |

---

## Bookings

### POST /bookings *(auth required)*
**Request:** `{ "phone", "event_id" }`
**Success — 201:** `data`: `BookingResponse`
**Failure:**
| Status | When |
|---|---|
| 400 | validation failed |
| 401 | missing/invalid token |
| 404 | `event_id` doesn't exist |
| 409 | user already booked this event |

### GET /bookings *(auth required)*
**Success — 200:** `data`: `[]BookingResponse`
**Failure:** none beyond auth — an empty list is success, not a 404.

### DELETE /bookings/{id} *(auth required)*
**Success — 200**
**Failure:**
| Status | When |
|---|---|
| 401 | missing/invalid token |
| 403 | authenticated user doesn't own this booking |
| 404 | no booking with that ID |

---

## Interactive docs

Full request/response schemas, try-it-out, and auth handling are available
via Swagger UI at `/swagger/index.html` once the server is running —
generated from the `@swag` annotations on each handler via `swag init`.
This file is a quick-reference summary; the generated OpenAPI spec is the
source of truth for exact field types.