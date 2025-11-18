package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gin-doniai/database"
    "gin-doniai/models"
)

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "用户创建成功",
        "user":    user,
    })
}

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
    var users []models.User

    result := database.DB.Find(&users)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "users": users,
        "count": len(users),
    })
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User

    result := database.DB.First(&user, id)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User

    // 先查找用户是否存在
    if err := database.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }

    // 绑定更新数据
    var updateData models.User
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 更新用户
    result := database.DB.Model(&user).Updates(updateData)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "用户更新成功",
        "user":    user,
    })
}

// DeleteUser 删除用户（软删除）
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User

    // 先查找用户是否存在
    if err := database.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }

    // 软删除
    result := database.DB.Delete(&user)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// 硬删除（永久删除）
func ForceDeleteUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User

    result := database.DB.Unscoped().Delete(&user, id)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "用户永久删除成功"})
}