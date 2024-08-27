Given the detailed `usecase` code you provided, I'll generate a README that explains how to handle different request and response scenarios. This README will include detailed API documentation, specifying exact response messages for each operation, which can help developers interact with the API correctly.

---

# LoanTracker 📊

Welcome to **LoanTracker**, a user management and authentication microservice built with Go, Gin, and MongoDB. This service handles user registration, login, email verification, password management, and more, with a clean and structured architecture.

## 🚀 Features

- User Registration & Login ✨
- Email Verification 📧
- Password Reset 🔒
- Token-based Authentication with JWT 🔑
- User Profile Management 👤
- Admin Capabilities 🛠️

## 🛠️ Tech Stack

- **Go** - High-performance programming language.
- **Gin** - Web framework used for building APIs.
- **MongoDB** - NoSQL database for scalable and flexible data storage.
- **JWT** - JSON Web Tokens for secure token-based authentication.

## 📦 Installation

Follow these steps to get the project up and running on your local machine.

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/loantracker.git
    cd loantracker
    ```

2. **Set up your environment variables:**

    Create a `.env` file at the root of the project and fill in the following details:

    ```bash
    DB_URI=<your_mongo_database_uri>
    JWT_SECRET=<your_jwt_secret>
    CONTEXT_TIMEOUT=<context_timeout_duration_in_seconds>
    ```

3. **Install dependencies:**

    Make sure you have Go installed, then run:

    ```bash
    go mod tidy
    ```

4. **Run the service:**

    ```bash
    go run main.go
    ```

## 📖 API Documentation

### User Routes

#### 1. **Register a New User**
- **POST `/user/create`**
- **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "username": "username",
      "password": "password123"
    }
    ```
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "User created successfully verify your account",
        "status": 200
      }
      ```
    - **Error (Validation):**
      ```json
      {
        "message": "All fields are required",
        "status": 400
      }
      ```
    - **Error (User Exists):**
      ```json
      {
        "message": "User already exists",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error creating user",
        "status": 500
      }
      ```

#### 2. **Verify User Email**
- **GET `/user/verify-email`**
- **Query Parameters:**
  - `token`: Email verification token.
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "User verified successfully",
        "status": 200
      }
      ```
    - **Error (Token Expired):**
      ```json
      {
        "message": "Token expired",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error finding user",
        "status": 500
      }
      ```

#### 3. **User Login**
- **POST `/user/login`**
- **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "Logged in Successfully",
        "status": 200,
        "access_token": "jwt_access_token",
        "refresh_token": "jwt_refresh_token"
      }
      ```
    - **Error (User Not Found):**
      ```json
      {
        "message": "User Not found",
        "status": 500
      }
      ```
    - **Error (User Not Verified):**
      ```json
      {
        "message": "User not verified",
        "status": 400
      }
      ```
    - **Error (Invalid Password):**
      ```json
      {
        "message": "Invalid password",
        "status": 400
      }
      ```

#### 4. **Forgot Password**
- **POST `/user/password-reset`**
- **Request Body:**
    ```json
    {
      "email": "user@example.com"
    }
    ```
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "Reset email sent",
        "status": 200
      }
      ```
    - **Error (User Not Found):**
      ```json
      {
        "message": "User not found",
        "status": 404
      }
      ```
    - **Error (Token Already Sent):**
      ```json
      {
        "message": "Reset token already sent. Please wait for X minutes to resend reset token",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error generating reset token",
        "status": 500
      }
      ```

#### 5. **Reset Password**
- **POST `/user/password-reset/:token`**
- **Request Body:**
    ```json
    {
      "password": "newpassword123"
    }
    ```
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "Password Reset Successfully",
        "status": 200
      }
      ```
    - **Error (Invalid Token):**
      ```json
      {
        "message": "Invalid reset token",
        "status": 400
      }
      ```
    - **Error (Token Expired):**
      ```json
      {
        "message": "Reset token expired",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error hashing password",
        "status": 500
      }
      ```

#### 6. **Refresh Token**
- **POST `/user/token/refresh`**
- **Request Body:**
    ```json
    {
      "refresh_token": "jwt_refresh_token"
    }
    ```
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "access_token": "new_jwt_access_token",
        "status": 200
      }
      ```
    - **Error (Invalid Token):**
      ```json
      {
        "message": "Invalid refresh token",
        "status": 400
      }
      ```
    - **Error (Token Expired):**
      ```json
      {
        "message": "Refresh Token Expired",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error generating token",
        "status": 500
      }
      ```

#### 7. **Get User Profile**
- **GET `/user/profile`**
- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "User found",
        "status": 200,
        "user": {
          "email": "user@example.com",
          "username": "username",
          "is_verified": true
        }
      }
      ```
    - **Error (Not Authorized):**
      ```json
      {
        "message": "Not Authorized",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error finding user",
        "status": 500
      }
      ```

### Admin Routes

#### 8. **Get All Users (Admin Only)**
- **GET `/admin/user`**
- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "Users found",
        "status": 200,
        "users": [
          {
            "email": "user@example.com",
            "username": "username",
            "is_verified": true
          },
          // more users
        ]
      }
      ```
    - **Error (Not Authorized):**
      ```json
      {
        "message": "Not Authorized",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error finding user",
        "status": 500
      }
      ```

#### 9. **Delete User (Admin Only)**
- **DELETE `/admin/user/:id`**
- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "User deleted successfully",
        "status": 200
      }
      ```
    - **Error (Not Authorized):**
      ```json
      {
        "message": "Not Authorized",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error deleting user",
        "status": 500
      }
      ```

Certainly! Here’s the continuation and completion of the README for the `LoanTracker` project.

---

## 🚀 API Documentation (Continued)

### User Routes (Continued)

#### 10. **Update User Profile**
- **PUT `/user/profile`**
- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Request Body:**
    ```json
    {
      "username": "newusername",
      "email": "newuser@example.com"
    }
    ```
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "Profile updated successfully",
        "status": 200
      }
      ```
    - **Error (Validation):**
      ```json
      {
        "message": "Invalid input",
        "status": 400
      }
      ```
    - **Error (Not Authorized):**
      ```json
      {
        "message": "Not Authorized",
        "status": 400
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error updating profile",
        "status": 500
      }
      ```

### Admin Routes (Continued)

#### 11. **Get User Details (Admin Only)**
- **GET `/admin/user/:id`**
- **Headers:**
  - `Authorization: Bearer <access_token>`
- **Possible Responses:**
    - **Success:**
      ```json
      {
        "message": "User details found",
        "status": 200,
        "user": {
          "email": "user@example.com",
          "username": "username",
          "is_verified": true,
          "role": "user"
        }
      }
      ```
    - **Error (Not Authorized):**
      ```json
      {
        "message": "Not Authorized",
        "status": 400
      }
      ```
    - **Error (User Not Found):**
      ```json
      {
        "message": "User not found",
        "status": 404
      }
      ```
    - **Error (Server Issues):**
      ```json
      {
        "message": "Error finding user",
        "status": 500
      }
      ```

## 🔒 Authentication & Authorization

- **Authentication** is handled using JWT tokens. Ensure you include the `Authorization` header with the `Bearer` token for routes that require authentication.

- **Authorization** is enforced for admin-specific routes. Only users with the role `admin` can access these endpoints.

## 🧪 Testing

You can test the API using tools like [Postman](https://www.postman.com/) or [cURL](https://curl.se/). Ensure your server is running before making requests.

### Sample cURL Commands

- **Register a User:**

    ```bash
    curl -X POST http://localhost:8080/user/create \
    -H "Content-Type: application/json" \
    -d '{"email": "user@example.com", "username": "username", "password": "password123"}'
    ```

- **Login:**

    ```bash
    curl -X POST http://localhost:8080/user/login \
    -H "Content-Type: application/json" \
    -d '{"email": "user@example.com", "password": "password123"}'
    ```

- **Get Profile:**

    ```bash
    curl -X GET http://localhost:8080/user/profile \
    -H "Authorization: Bearer <access_token>"
    ```


