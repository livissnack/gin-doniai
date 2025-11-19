package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gin-doniai/database"
	"gin-doniai/handlers"
	"gin-doniai/models"
	"gin-doniai/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 在 main.go 顶部添加全局变量
var (
	globalConfig          GlobalConfig
	recommendedCategories []models.Category
)

// 修改 GlobalConfig 结构体，添加 Categories 字段
type GlobalConfig struct {
	SiteName   string
	Theme      string
	Version    string
	Categories []models.Category // 添加推荐分类字段
}

func main() {
	database.InitDB()

	// 初始化全局配置
	globalConfig = GlobalConfig{
		SiteName: "Doniai",
		Theme:    "light",
		Version:  "1.0.0",
	}

	// 获取推荐分类并注入到全局配置
	if categories, err := handlers.GetRecommendedCategories(); err == nil {
		recommendedCategories = categories
		globalConfig.Categories = categories
	} else {
		fmt.Printf("获取推荐分类失败: %v\n", err)
	}

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"loop": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
		// 添加获取全局配置的函数
		"global": func() GlobalConfig {
			return globalConfig
		},
	})
	// 设置session存储
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	// 在路由定义之前应用用户中间件
	router.Use(UserMiddleware())
	router.Use(OnlineStatusMiddleware())

	// 加载模板文件
	router.LoadHTMLGlob("templates/**/*")

	// 静态文件服务
	router.Static("/static", "./static")

	// 路由定义
	router.GET("/", homeHandler)
	router.GET("/categories/:type", homeHandler)
	router.GET("/about", aboutHandler)
	router.GET("/post-:id-1", detailHandler)
	router.GET("/register", registerHandler)
	router.POST("/register", registerSubmit)
	router.GET("/login", loginHandler)
	router.POST("/login", loginSubmit)
	router.GET("/logout", logoutHandler)
	router.GET("/profile", profileHandler)
	router.GET("/publish", publishHandler)
	router.GET("/settings", settingsHandler)

	// 在 main.go 的路由部分添加
	router.GET("/api/online/count", handlers.GetOnlineUserCount)

	// 在 main.go 的路由定义部分添加评论路由
	commentRoutes := router.Group("/api/comments")
	{
		commentRoutes.POST("/", handlers.CreateComment)
		commentRoutes.GET("/", handlers.GetComments)
		commentRoutes.GET("/:id", handlers.GetComment)
		commentRoutes.PUT("/:id", handlers.UpdateComment)
		commentRoutes.DELETE("/:id", handlers.DeleteComment)
		commentRoutes.POST("/:id/like", handlers.LikeComment)
	}

	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/", handlers.CreateUser)                 // 创建用户
		userRoutes.GET("/", handlers.GetUsers)                    // 获取所有用户
		userRoutes.GET("/:id", handlers.GetUser)                  // 获取单个用户
		userRoutes.PUT("/:id", handlers.UpdateUser)               // 更新用户
		userRoutes.DELETE("/:id", handlers.DeleteUser)            // 删除用户（软删除）
		userRoutes.DELETE("/:id/force", handlers.ForceDeleteUser) // 强制删除
	}

	postRoutes := router.Group("/api/posts")
	{
		postRoutes.POST("/", handlers.CreatePost)                 // 创建文章
		postRoutes.GET("/", handlers.GetPosts)                    // 获取所有文章
		postRoutes.GET("/:id", handlers.GetPost)                  // 获取单个文章
		postRoutes.PUT("/:id", handlers.UpdatePost)               // 更新文章
		postRoutes.DELETE("/:id", handlers.DeletePost)            // 删除文章（软删除）
		postRoutes.POST("/:id/like", handlers.LikePost)           // 文章点赞
		postRoutes.DELETE("/:id/force", handlers.ForceDeletePost) // 强制删除
		postRoutes.POST("/:id/favorite", handlers.FavoritePost)   // 文章收藏
	}

	// 启动定时清理任务
	go func() {
		ticker := time.NewTicker(10 * time.Minute) // 每10分钟清理一次
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				handlers.CleanupExpiredOnlineStatus()
			}
		}
	}()

	router.Run(":8080")
}

// UserMiddleware 全局获取用户信息的中间件
func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		var user *models.User
		if userID != nil {
			var currentUser models.User
			if err := database.DB.First(&currentUser, userID).Error; err == nil {
				user = &currentUser
			}
		} else {
			fmt.Println("Middleware - Session中没有user_id")
		}

		c.Set("user", user)
		c.Next()
	}
}

// OnlineStatusMiddleware 在线状态中间件
func OnlineStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 处理请求前更新在线状态
		handlers.UpdateUserOnlineStatus(c)

		c.Next()

		// 请求处理后也可以选择再次更新
		// UpdateUserOnlineStatus(c)
	}
}

// AuthMiddleware 检查用户是否已登录的中间件
// func AuthMiddleware() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         session := sessions.Default(c)
//         userID := session.Get("user_id")
//
//         if userID == nil {
//             // 用户未登录，重定向到登录页面
//             c.Redirect(http.StatusFound, "/login")
//             c.Abort()
//             return
//         }
//
//         // 用户已登录，继续处理请求
//         c.Next()
//     }
// }

func homeHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userObj, exists := c.Get("user")
	var user *models.User
	if exists && userObj != nil {
		user = userObj.(*models.User)
	} else {
		fmt.Println("未获取到用户信息")
	}
	// 获取页码参数，默认为第1页
	pageStr := c.Query("page")
	// 获取路由categories/:type中type参数
	categoryType := c.Param("type")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 每页显示的帖子数量
	limit := 10
	offset := (page - 1) * limit

	// 查询总记录数
	var total int64
	dbQuery := database.DB.Model(&models.Post{})

	var categoryId uint
	if categoryType != "" {
		var category models.Category
		if err := database.DB.Where("alias = ?", categoryType).First(&category).Error; err == nil {
			categoryId = category.ID
			dbQuery = dbQuery.Where("category_id = ?", categoryId)
		}
	}
	dbQuery.Count(&total)

	// 查询当前页的帖子
	var posts []models.Post
	postQuery := database.DB.Order("created_at DESC").Offset(offset).Limit(limit)

	if categoryId > 0 {
		postQuery = postQuery.Where("category_id = ?", categoryId)
	}
	postQuery.Find(&posts)

	// 计算总页数
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	// 创建带有友好时间的帖子结构
	type PostWithFriendlyTime struct {
		models.Post
		TimeAgo string
	}

	var postsWithTimeAgo []PostWithFriendlyTime
	for _, post := range posts {
		timeAgo := utils.GetTimeAgo(post.CreatedAt)
		postsWithTimeAgo = append(postsWithTimeAgo, PostWithFriendlyTime{
			Post:    post,
			TimeAgo: timeAgo,
		})
	}

	// 获取在线用户数
	var onlineCount int64
	cutoffTime := time.Now().Add(-30 * time.Minute)
	database.DB.Model(&models.UserOnlineStatus{}).
		Where("last_active_time > ?", cutoffTime).
		Count(&onlineCount)

	// 获取站点统计信息
	var userCount, postCount, commentCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Post{}).Count(&postCount)
	database.DB.Model(&models.Comment{}).Count(&commentCount)

	// 获取所有分类
	var categories []models.Category
	database.DB.Where("status_code = ?", 1).Find(&categories)

	// 1、统计注册用户
	// 2、统计文章数量
	// 3、统计回复评论数量
	// 4、数据表categories查询所有分类

	data := gin.H{
		"CurrentTime":  time.Now().Format("2006-01-02 15:04:05"),
		"posts":        postsWithTimeAgo,
		"currentPage":  page,
		"totalPages":   totalPages,
		"hasPrev":      page > 1,
		"hasNext":      page < totalPages,
		"prevPage":     page - 1,
		"nextPage":     page + 1,
		"user":         user,
		"userCount":    userCount,
		"postCount":    postCount,
		"commentCount": commentCount,
		"onlineCount":  onlineCount,
		"categories":   categories,
	}

	c.HTML(http.StatusOK, "home.tmpl", data)
}

func aboutHandler(c *gin.Context) {
	data := gin.H{
		"TeamMembers": []struct {
			Name string
			Role string
		}{
			{"张三", "开发工程师"},
			{"李四", "UI设计师"},
			{"王五", "产品经理"},
		},
	}
	c.HTML(http.StatusOK, "about.tmpl", data)
}

func detailHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userObj, exists := c.Get("user")
	var user *models.User
	if exists && userObj != nil {
		user = userObj.(*models.User)
	}
	// 从路由中获取文章ID
	idWithSuffix := c.Param("id-1") // 获取 "29-1"
	// 分割字符串获取真正的ID
	idParts := strings.Split(idWithSuffix, "-")
	var id string
	if len(idParts) > 0 {
		id = idParts[0] // 获取 "29"
	}

	// 查询数据库获取文章详情，并预加载用户信息
	var post models.Post
	if err := database.DB.Preload("User").First(&post, id).Error; err != nil {
		fmt.Printf("查询文章失败: %v\n", err) // 添加调试信息
		// 如果找不到文章，返回404
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"Message": "文章未找到",
		})
		return
	}

	// 创建带有友好时间和回复评论的评论结构
	type CommentWithReplies struct {
		models.Comment
		TimeAgo string
		Replies []CommentWithReplies
		Content template.HTML
	}

	// 获取评论页码参数
	commentPageStr := c.Query("page")
	commentPage := 1
	if commentPageStr != "" {
		if p, err := strconv.Atoi(commentPageStr); err == nil && p > 0 {
			commentPage = p
		}
	}

	// 每页显示的评论数量
	commentLimit := 4
	commentOffset := (commentPage - 1) * commentLimit

	// 查询总评论数
	var totalComments int64
	database.DB.Model(&models.Comment{}).Where("post_id = ? AND parent_id = 0", id).Count(&totalComments)

	// 计算总页数
	totalCommentPages := int((totalComments + int64(commentLimit) - 1) / int64(commentLimit))

	// 查询该文章的评论（仅顶级评论），并预加载用户信息
	var comments []models.Comment
	database.DB.Where("post_id = ? AND parent_id = 0", id).Preload("User").
		Offset(commentOffset).Limit(commentLimit).Order("created_at DESC").Find(&comments)

	// 创建带有回复和友好时间的评论列表
	var commentsWithReplies []CommentWithReplies
	for _, comment := range comments {
		// 处理主评论的友好时间
		timeAgo := utils.GetTimeAgo(comment.CreatedAt)

		// 查询该评论的回复
		var replies []models.Comment
		database.DB.Where("parent_id = ?", comment.ID).Preload("User").Find(&replies)

		// 处理回复评论的友好时间
		var repliesWithTime []CommentWithReplies
		for _, reply := range replies {
			replyTimeAgo := utils.GetTimeAgo(reply.CreatedAt)
			repliesWithTime = append(repliesWithTime, CommentWithReplies{
				Comment: reply,
				TimeAgo: replyTimeAgo,
				Content: template.HTML(reply.Content),
			})
		}

		commentsWithReplies = append(commentsWithReplies, CommentWithReplies{
			Comment: comment,
			TimeAgo: timeAgo,
			Replies: repliesWithTime,
			Content: template.HTML(comment.Content),
		})
	}

	// 将标签字符串分割成数组
	// 使用工具方法处理标签
	tags := utils.ParseTags(post.Tags)
	// 将文章详情数据和用户信息传递给模板
	var postCount, replyCount, likeCount int64
	database.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&postCount)
	database.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Select("SUM(replies)").Row().Scan(&replyCount)
	database.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Select("SUM(likes)").Row().Scan(&likeCount)
	// 将评论分页信息添加到模板数据
	data := gin.H{
		"Post":               post,
		"User":               post.User,
		"Content":            template.HTML(post.Content),
		"Tags":               tags,
		"user":               user,
		"Comments":           commentsWithReplies,
		"commentCurrentPage": commentPage,
		"commentTotalPages":  totalCommentPages,
		"commentHasPrev":     commentPage > 1,
		"commentHasNext":     commentPage < totalCommentPages,
		"commentPrevPage":    commentPage - 1,
		"commentNextPage":    commentPage + 1,
		"commentTotalCount":  totalComments,
		"postCount":          postCount,
		"replyCount":         replyCount,
		"likeCount":          likeCount,
	}

	c.HTML(http.StatusOK, "detail.tmpl", data)
}

func profileHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userObj, exists := c.Get("user")
	var user *models.User
	if exists && userObj != nil {
		user = userObj.(*models.User)
	} else {
		// 用户未登录，重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
		return
	}

	data := gin.H{
		"user":        user,
		"profileUser": user, // 当前查看的用户（自己）
	}
	// 打印用户详细信息
	fmt.Printf("当前用户信息: %+v\n", user)
	c.HTML(http.StatusOK, "profile.tmpl", data)
}

func settingsHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userObj, exists := c.Get("user")
	var user *models.User
	if exists && userObj != nil {
		user = userObj.(*models.User)
	} else {
		// 用户未登录，重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
		return
	}

	data := gin.H{
		"user": user,
	}
	c.HTML(http.StatusOK, "settings.tmpl", data)
}

func publishHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userObj, exists := c.Get("user")
	var user *models.User
	if exists && userObj != nil {
		user = userObj.(*models.User)
	} else {
		// 用户未登录，重定向到登录页面
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// 获取所有分类
	var categories []models.Category
	database.DB.Where("status_code = ?", 1).Find(&categories)

	data := gin.H{
		"user":       user,
		"categories": categories,
	}
	c.HTML(http.StatusOK, "publish.tmpl", data)
}

func registerHandler(c *gin.Context) {
	data := gin.H{
		"Message":     "欢迎使用 Gin 模板",
		"CurrentPath": "/register",
	}
	c.HTML(http.StatusOK, "auth.tmpl", data)
}

func loginHandler(c *gin.Context) {
	data := gin.H{
		"Message":     "欢迎使用 Gin 模板",
		"CurrentPath": "/login",
	}
	c.HTML(http.StatusOK, "auth.tmpl", data)
}

func logoutHandler(c *gin.Context) {
	// 获取session
	session := sessions.Default(c)

	// 清除session中的用户信息
	session.Clear()

	// 保存session更改
	if err := session.Save(); err != nil {
		fmt.Printf("登出时Session保存失败: %v\n", err)
	} else {
		fmt.Println("用户已成功登出")
	}

	// 重定向到首页
	c.Redirect(http.StatusFound, "/")
}

func loginSubmit(c *gin.Context) {

	// 测试密码：xZ3(Uq)sDQ6qYEY]
	// 获取表单提交的数据
	identifier := c.PostForm("email") // 可以是邮箱或用户名
	password := c.PostForm("password")
	remember := c.PostForm("remember")

	// 基本验证
	if identifier == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "用户名/邮箱和密码不能为空",
		})
		return
	}

	// 查询用户（支持邮箱或用户名登录）
	var user models.User
	if err := database.DB.Where("email = ? OR name = ?", identifier, identifier).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "用户不存在",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "密码错误",
		})
		return
	}

	// 获取session
	session := sessions.Default(c)

	// 设置用户信息到session
	session.Set("user_id", user.ID)

	// 根据"记住密码"选项设置过期时间
	if remember == "on" {
		// 设置30天过期
		session.Options(sessions.Options{
			MaxAge: 30 * 24 * 60 * 60, // 30天
		})
	} else {
		// 设置会话结束时失效（浏览器关闭时）
		session.Options(sessions.Options{
			MaxAge: 0, // 浏览器会话期间有效
		})
	}

	// 保存session

	// 保存session
	if err := session.Save(); err != nil {
		fmt.Printf("Session保存失败: %v\n", err)
	} else {
		fmt.Println("Session保存成功")
	}

	// 登录成功
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "登录成功",
		"data": gin.H{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		},
	})

}

func registerSubmit(c *gin.Context) {
	// 获取表单提交的数据
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirmPassword")
	agreeTerms := c.PostForm("agreeTerms")

	// 基本验证
	if username == "" || email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "用户名、邮箱和密码不能为空",
		})
		return
	}

	if password != confirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "两次输入的密码不一致",
		})
		return
	}

	// 处理用户协议同意状态
	isAgreeTerms := false
	if agreeTerms == "on" {
		isAgreeTerms = true
	}

	if !isAgreeTerms {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请同意用户协议",
		})
		return
	}

	// 检查用户是否已存在
	var existingUser models.User
	if err := database.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "该邮箱已被注册",
		})
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "密码加密失败",
		})
		return
	}
	// 随机一个avatar图像
	// 生成基于用户名的随机头像
	avatarURL := fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", username)

	// 创建新用户
	newUser := models.User{
		Name:       username,
		Email:      email,
		Password:   hashedPassword,
		AgreeTerms: isAgreeTerms,
		Avatar:     avatarURL,
	}

	// 保存到数据库
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "用户注册失败",
		})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login data received",
		"data": gin.H{
			"user_id": newUser.ID,
			"email":   newUser.Email,
		},
	})
}
