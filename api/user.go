package api

import (
	"angular/Db2"
	"angular/LOG"
	"angular/model"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func UploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		LOG.ErrorLogger.Error("File upload error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	filename := fmt.Sprintf("avatars/%s%s", userID, filepath.Ext(filepath.Base(file.Filename)))

	if err := c.SaveUploadedFile(file, filename); err != nil {
		//	main2.ErrorLogger.Error("File save error:", err)
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		//	return
	}

	_, err = Db2.DB.Exec("UPDATE users SET avatar = $1 WHERE id = $2", filename, userID)
	if err != nil {
		//main2.ErrorLogger.Error("Database update error:", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		//return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully!", "avatar_url": filename})
}

func AddUser(c *gin.Context) {
	var avatarPath string
	file, _ := c.FormFile("avatar")
	if file != nil {
		avatarPath = fmt.Sprintf("uploads/%s", filepath.Base(file.Filename))
		_ = c.SaveUploadedFile(file, avatarPath)
	}

	var lastID int
	err := Db2.DB.QueryRow(
		"INSERT INTO users(username, email, comment, avatar) VALUES ($1, $2, $3, $4) RETURNING id",
		c.PostForm("username"), c.PostForm("email"), c.PostForm("comment"), avatarPath,
	).Scan(&lastID)

	if err != nil {
		//main2.ErrorLogger.Error("Insert error:", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added", "id": lastID, "avatar": avatarPath})
}

func GetUsers(c *gin.Context) {
	rows, err := Db2.DB.Query("SELECT id, username, email, comment, avatar FROM users")
	if err != nil {
		//main2.ErrorLogger.Error(err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var avatar sql.NullString // 使用 sql.NullString 处理可能为 NULL 的字段
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Comment, &avatar)
		if err != nil {
			//main2.ErrorLogger.Error(err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
		}

		// 如果 avatar 是 NULL，则设置为空字符串
		if avatar.Valid {
			user.Avatar = avatar.String
		} else {
			user.Avatar = "alex.jpg" // 或者可以设置一个默认头像的路径
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	result, err := Db2.DB.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		//main2.ErrorLogger.Error(err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User does not exist", "username": username})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
