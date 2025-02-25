# T002 - Authentication & Authorization

Build authentication and authorization features on user service.

## Technical Requirements

- [x] Authentication
  - [x] JWT
  - [x] `POST /v1/register`
    - [ ] ~~OAuth and Google Idp?~~
  - [ ] `POST /v1/forgot-password`
  - [ ] `POST /v1/reset-password`
  - [x] `POST /v1/login`
  - [x] `GET /v1/token`
  - [x] `POST /v1/logout`
  - [x] `GET /v1/profile`
  - [x] `PATCH /v1/profile`
  - [x] `GET /v1/profile/email/verification`
  - [x] `POST /v1/profile/email/verification`
  - [x] `GET /v1/profile/phone/verification`
  - [x] `POST /v1/profile/phone/verification`
- [x] Authorization
  - [x] Roles: `admin`, `user`
  - [x] Open Policy Agent (OPA) Integration
