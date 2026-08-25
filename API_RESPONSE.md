# API Response Reference

Every response — success or failure — uses the same envelope
(`internal/response`):

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
`data`/`meta` are omitted. Failure status codes below come from the
`apperror` constructor named in parentheses — that's what each service
method should return, not a raw error.

---

## Auth

### POST /signup
**Request:** `{ "name", "email", "password" }`
**Success — 201:**
```json
{ "success": true, "message": "User created successfully", "data": UserResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 400 | validation failed (missing field, bad email, password < 6 chars) | `BadRequest` |
| 409 | email already registered | `Conflict` |

### POST /signin
**Request:** `{ "email", "password" }`
**Success — 200:**
```json
{ "success": true, "message": "User signed in successfully", "data": SignInResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 400 | validation failed | `BadRequest` |
| 401 | email not found OR password mismatch (same message either way) | `Unauthorized` |

### GET /me *(auth required)*
**Success — 200:**
```json
{ "success": true, "data": UserResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 401 | missing/invalid token | `Unauthorized` |
| 404 | user no longer exists | `NotFound` |

---

## Events

### POST /events *(auth required, multipart form)*
**Request:** form fields `name`, `description`, `location`, `datetime` (RFC3339) + `image` file
**Success — 201:**
```json
{ "success": true, "message": "New event created", "data": EventResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 400 | missing image file, missing/invalid required field, bad datetime format | `BadRequest` |
| 500 | image upload to ImageKit failed | `Internal` |

### GET /events
**Query params:** `search` (optional), `page` (default 1), `limit` (default 6)
**Success — 200:**
```json
{ "success": true, "message": "Events retrieved", "data": [EventResponse], "meta": EventListMeta }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 500 | query failed | `Internal` |

### GET /events/:id
**Success — 200:**
```json
{ "success": true, "data": EventDetailResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 404 | no event with that ID | `NotFound` |

### GET /events/mine *(auth required)*
**Success — 200:**
```json
{ "success": true, "data": [EventResponse] }
```
**Failure:** none beyond auth (an empty list is success, not a 404)

### PUT /events/:id *(auth required, multipart form, all fields optional)*
**Success — 200:**
```json
{ "success": true, "message": "Event updated successfully", "data": EventResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 400 | bad datetime format, malformed image upload | `BadRequest` |
| 403 | authenticated user doesn't own this event | `Forbidden` |
| 404 | no event with that ID | `NotFound` |
| 500 | image upload/delete failed, save failed | `Internal` |

### DELETE /events/:id *(auth required)*
**Success — 200:**
```json
{ "success": true, "message": "Event deleted" }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 403 | authenticated user doesn't own this event | `Forbidden` |
| 404 | no event with that ID | `NotFound` |

---

## Bookings

### POST /bookings *(auth required)*
**Request:** `{ "phone", "event_id" }`
**Success — 201:**
```json
{ "success": true, "message": "Booking created successfully", "data": BookingResponse }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 400 | validation failed | `BadRequest` |
| 404 | event_id doesn't exist | `NotFound` |
| 409 | user already booked this event | `Conflict` |

### GET /bookings *(auth required)*
**Success — 200:**
```json
{ "success": true, "data": [BookingResponse] }
```
**Failure:** none beyond auth (empty list is success, not a 404)

### DELETE /bookings/:id *(auth required)*
**Success — 200:**
```json
{ "success": true, "message": "Booking deleted" }
```
**Failure:**
| Status | When | apperror |
|---|---|---|
| 403 | authenticated user doesn't own this booking | `Forbidden` |
| 404 | no booking with that ID | `NotFound` |

---

## Notes / decisions this spec locks in

- **List endpoints return `[]` on empty, never 404.** "No results" is a
  valid successful state — 404 is reserved for "the specific resource you
  asked for by ID doesn't exist."
- **Sign-in failure is deliberately generic** — "invalid email or
  password" for both wrong-email and wrong-password cases, so the
  response never confirms whether an email is registered.
- **`EventResponse` vs `EventDetailResponse`** exist as two DTOs because
  the two queries behind them (`FindAll` vs `FindByID`) preload different
  data. The list DTO doesn't carry a promise (`bookings: []`) that the
  underlying query never actually fulfilled.
- Every 4xx/5xx above is what the **service layer** should return as an
  `*apperror.AppError` — the handler layer (Phase 4) just calls
  `response.FromError(c, err)` and this table is automatically satisfied,
  no per-handler status-code decisions left to make.