package models

import "time"

type Directory struct {
    Id        int       `json:"id"`
    Name      string    `json:"name"`
    UserId    int       `json:"user_id"`
    ParentId  *int      `json:"parent_id"` // Nullable, for nested directories
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
} 