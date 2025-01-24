package main

import (
	"Grade_managing/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个 Gin 实例
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // 允许跨域请求带上 cookies 或 HTTP 认证
		MaxAge:           12 * time.Hour,
	}))

	// 学生路由
	r.POST("/students", controllers.AddOrUpdateStudent)
	r.GET("/students/:id", controllers.GetStudent)
	r.GET("/students", controllers.GetAllStudents)
	r.PUT("/students/:id", controllers.UpdateStudentInfo)
	r.DELETE("/students/:id", controllers.DeleteStudent)

	// 成绩路由
	r.POST("/students/:id/grades", controllers.AddGrade)
	// 学号课程名称修改成绩
	r.PUT("/students/:id/grades/:course", controllers.UpdateCourseGrade)
	r.DELETE("/students/:id/grades", controllers.DeleteGrade)

	//csv 文件上传路由
	r.POST("/upload", controllers.UploadCSV)

	_ = r.Run(":8080") //启动 HTTP 服务
}
