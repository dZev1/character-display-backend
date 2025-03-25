# Character Display Server API
Making a RESTful API just to learn how to use the net/http Go package.

## Overview

This API provides endpoints for user authentication and character management in a character display system.

## Authentication

Some routes require authentication by a session_token. Users have to login to receive a session cookie (`session_token`) and a CSRF token cookie (`csrf_token`)

---

## Base URL

```
http://<server-address>/
```

## **Endpoints** 

### **1. Authentication**

#### **Register**
**Endpoint**: `POST /register`

**Description**: Registers a new user.

**Request Parameters**:
   - `username`(string, required): The username (minimum 3 characters).
   - `password`(string, required): The password (minimum 8 characters).

**Response**:
   - `200 OK`: User registered successfully.
   - `406 Not Acceptable`: Invalid username or password.
   - `409 Conflict`: User already exists.
   - `500 Internal Server Error`: Could not insert user.

---

#### **Login**
**Endpoint**: `POST /login`

**Description**: Registers a new user.

**Request Parameters**:
   - `username`(string, required)
   - `password`(string, required)
**Response**:
   - `200 OK`: Login successful, session and CSRF tokens set in cookies.
   - `401 Unauthorized`: User not authorized.
   - `400 Bad Request`: Request could not be processed.

**Generated Cookies**:
   - `session_token` (HTTPOnly, valid for 24hs)
   - `csrf_token`

---

#### **Logout** *(Requires Authentication)*
**Endpoint**: `POST /logout`

**Description**: Logs out a user by invalidating session and CSRF tokens.

**Request Parameters**:
   - `username`(string, required)

**Response**:
   - `200 OK`: Logged out succesfully.
   - `401 Unauthorized`: Invalid credentials.
   - `400 Bad Request`: Request could not be processed.

---

### **2. Character Management**

#### **Upload Character** *(Requires authentication)*
**Endpoint**: `POST /upload_character`

**Description**: Uploads a new character to the system.

**Request Parameters**:
   - `username`(string, required)
   - `char_json`(JSON, required): Character data.

**Response**:
   - `200 OK`: Character added successfully.
   - `400 Bad Request`: Invalid data.
   - `401 Unauthorized`: User not authorized.

---

#### **Edit Character** *(Requires authentication)*
**Endpoints**:
   - `GET /edit_character`: Retrieves a character.
   - `PUT /edit_character`: Updates a character.

**Request Parameters**:
   - `username`(string, required)
   - `char_name` (string, required)
   - `char_json`(string, required for PUT): Updated character data.

**Response**:
   - `200 OK`: Character added successfully.
   - `400 Bad Request`: Invalid data.
   - `401 Unauthorized`: User not authorized.
   - `409 Conflict`: Character not found (GET Request).

---

#### **Delete Character** *(Requires authentication)*
**Endpoint**: `DELETE /delete_character`

**Description**: Deletes a character

**Request Parameters**:
   - `username`(string, required)
   - `char_name` (string, required)

**Response**:
   - `200 OK`: Character deleted succesfully.
   - `400 Bad Request`: Invalid data.
   - `401 Unauthorized`: User not authorized.

---

#### **Get Characters**
**Endpoint**: `GET /get_characters`

**Description**: Retrieves characters based on a specified field and value.

**Request Parameters**:
   - `field`(string, required)
   - `value` (string, required)

**Response**:
   - `200 OK`: Returns a JSON array of matching characters.
   - `400 Bad Request`: Invalid data.

---

#### **Get All Characters**
**Endpoint**: `GET /get_all_characters`

**Description**: Retrieves all characters in the system.

**Response**:
   - `200 OK`: Returns a JSON array of all characters.
   - `400 Bad Request`: Request could not be processed.

