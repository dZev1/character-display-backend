# Character Display Server API
Making a RESTful API just to learn how to use the net/http Go package.

## Overview

This API provides endpoints for user authentication and character management in a character display system.

## Base URL

```
http://<server-address>/
```

## Authentication Endpoints


### Register
__Endpoint__: `POST /register`

__Description__: Registers a new user.

__Request Parameters__:
   - `username`(string, required): The username (minimum 3 characters).
   - `password`(string, required): The password (minimum 8 characters).

__Response__:
   - `200 OK`: User registered successfully.
   - `406 Not Acceptable`: Invalid username or password.
   - `409 Conflict`: User already exists.
   - `500 Internal Server Error`: Could not insert user.

-----

### Login
__Endpoint__: `POST /login`

__Description__: Registers a new user.

__Request Parameters__:
   - `username`(string, required): The username (minimum 3 characters).
   - `password`(string, required): The password (minimum 8 characters).

__Response__:
   - `200 OK`: Login successful, session and CSRF tokens set in cookies..
   - `401 Unauthorized`: Invalid credentials.
   - `400 Bad Request`: Request could not be processed.
