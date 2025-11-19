package models

import (
    "gorm.io/gorm"
    "time"
)

type Post struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Title     string         `json:"title" gorm:"size:200;not null"`
    UserId    int            `json:"user_id" gorm:"not null"`
    Author    string         `json:"author" gorm:"size:40;not null"`
    Category  string         `json:"category" gorm:"size:100;not null"`
    Content   string         `json:"content" gorm:"type:text;not null"`
    Tags      string         `json:"tags" gorm:"size:255;not null"`
    Views     int            `json:"views" gorm:"default:0"`      // 浏览数
    Replies   int            `json:"replies" gorm:"default:0"`    // 回复数
    Favorites int            `json:"favorites" gorm:"default:0"`  // 收藏数
    Likes     int            `json:"likes" gorm:"default:0"`      // 点赞数
    ReadLimit int            `json:"read_limit" gorm:"default:1"` // 阅读限制: 1-公开, 2-Lv1, 3-Lv2, 4-私有
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

    // 明确指定外键关系
    User User `json:"user" gorm:"foreignKey:UserId;references:ID"`
}


// 表名
func (Post) TableName() string {
    return "posts"
}