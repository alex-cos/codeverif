# codeverif

**codeverif** is a secure and lightweight Go module for generating and verifying **one-time verification codes**.

It is designed for common authentication flows such as:

- Login verification
- Password reset
- Email confirmation

The module is backend-oriented and integrates naturally into REST APIs (Gin, Echo, net/http).

---

## ✨ Features

- 🔐 Cryptographically secure numeric codes
- ⏱ Code expiration handling
- 🔁 Limited verification attempts
- 🧪 Unit-tested
- 🌐 REST API friendly

---

## 📦 Installation

```bash
go get github.com/your-org/codeverif
```

---

## Example

```go
// create a codeverif for a 6 digits code available during 10 minutes with 3 attempts.
cv := codeverif.New(6, 10*time.Minute, 3)

code, err := cv.RequestCode("JohnDoe")
if err != nil {
  return err
}

// Send the generated code by mail to the end user for verification

// Then when asked verified the entered code.
err = cv.VerifyCode("JohnDoe", code)
if err == nil {
  // Code successfully verified
}

```
