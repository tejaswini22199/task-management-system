package models

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`  // `-` means this field won't be included in JSON responses
    CreatedAt time.Time `json:"created_at"`
    // You might want to add more fields like:
    // UpdatedAt time.Time `json:"updated_at"`
    // FirstName string    `json:"first_name"`
    // LastName  string    `json:"last_name"`
    // IsActive  bool      `json:"is_active"`
}