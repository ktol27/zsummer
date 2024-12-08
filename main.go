package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// 用户模型

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Comment  string `json:"comment"`
}

// 数据库连接

var DB *sql.DB

// 初始化数据库
func initDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 添加重试逻辑
	for i := 0; i < 5; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("Attempting to connect to database... (attempt %d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after 5 attempts")
	}

	fmt.Println("Successfully connected to database")
}

// 获取用户列表
func getUsers(c *gin.Context) {
	rows, err := DB.Query("SELECT id username, email, comment FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Email, &user.Comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

//添加用户(POST)

func addUser(c *gin.Context) {
	var newUser User
	if err := c.Bind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := DB.Exec("INSERT INTO users(username,email,comment)VALUES ($1, $2, $3)",
		newUser.Username, newUser.Email, newUser.Comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{"message:": "User added successfully", "rowsAffected": rowsAffected})
}

//删除用户（delete）

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	result, err := DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User does not exist", "id": id})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

//更新用户（PUT）

func updateUser(c *gin.Context) {
	id := c.Param("id")
	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := DB.Exec("UPDATE users SET username = $1, email = $2, comment = $3 WHERE id = $4",
		updatedUser.Username, updatedUser.Email, updatedUser.Comment, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func main() {
	// 初始化数据库
	initDatabase()

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.Default())

	// 用户路由
	r.GET("/users", getUsers)          // 获取用户列表
	r.POST("/users", addUser)          // 添加用户
	r.DELETE("/users/:id", deleteUser) // 删除用户
	r.PUT("/users/:id", updateUser)    // 更新用户

	// 启动服务器
	r.Run(":8080")
}
