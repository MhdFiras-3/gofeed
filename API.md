# API Documentation

## Authentication & Users

### Register User

Creates a new user account with a hashed password.

* **URL:** `/register` (or `/api/v1/register` depending on your router mount)
* **Method:** `POST`
* **Authentication Required:** No (Public)
* **Headers:**
  * `Content-Type: application/json`

---

#### Request Body

| Field | Type | Required | Constraints | Description |
| :--- | :--- | :--- | :--- | :--- |
| `name` | string | Yes | 1–45 chars, no control characters | User's display name. Leading/trailing whitespace is trimmed. |
| `email` | string | Yes | 1–50 chars, valid email format | Unique email address. Checked via RFC 5322 parsing. |
| `password` | string | Yes | 20–100 chars | Plaintext password to hash. |

**Example Request:**

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "correct horse battery staple 123"
}
```
#### Responses

##### `201 Created`
Returned when the user is successfully created.

```json
{
  "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

##### `400 Bad Request`
Returned when the JSON payload is malformed or input validation fails.

*Malformed JSON:*
```json
{
  "error": "invalid request payload"
}
```

*Validation Errors:*
```json
{
  "errors": [
    {
      "field": "email",
      "message": "invalid email format"
    }
  ]
}
```

##### `409 Conflict`
Returned when the email address is already registered.

```json
{
  "error": "email already registered"
}
```

##### `500 Internal Server Error`
Returned when an unhandled server error occurs (e.g., password hashing failure or database connection failure).

```json
{
  "error": "something went wrong"
}
```