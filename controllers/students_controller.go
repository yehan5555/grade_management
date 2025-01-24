package controllers

import (
	"Grade_managing/models"
	"Grade_managing/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"Grade_managing/services"
	"github.com/gin-gonic/gin"
)

const (
	StatusOK         = 200
	StatusBadRequest = 400
	StatusNotFound   = 404

	AddOrUpdateStudentFailed = 10000
)

// 添加或更新学生信息

func AddOrUpdateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(StatusBadRequest, gin.H{"Bad Request": err.Error()})
		return
	}

	if err := services.AddStudent(student); err != nil {
		c.JSON(AddOrUpdateStudentFailed, gin.H{"Add Or Update Failed": err.Error()})
	}
	c.JSON(StatusOK, gin.H{"message": "student added or updated successfully"})
}

// 获取学生信息（根据学号）

func GetStudent(c *gin.Context) {
	id := c.Param("id")
	student, err := services.GetStudentByID(id)
	if err != nil {
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(StatusOK, student)
}

// 获取所有学生信息

func GetAllStudents(c *gin.Context) {
	students := services.GetAllStudents()
	c.JSON(StatusOK, students)
}

// 删除学生信息（根据学号）

func DeleteStudent(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteStudentByID(id); err != nil {
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(StatusOK, gin.H{"message": "student deleted successfully"})
}

// 修改学生信息

func UpdateStudentInfo(c *gin.Context) {
	id := c.Param("id")
	var newInfo models.Student
	if err := c.ShouldBindJSON(&newInfo); err != nil {
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateStudentInfo(id, newInfo); err != nil {
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(StatusOK, gin.H{"message": "student info updated successfully"})
}

// 添加或修改课程成绩

func AddGrade(c *gin.Context) {
	id := c.Param("id")
	var grade models.Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if grade.Course == "" || grade.Score < 0 || grade.Score > 100 {
		c.JSON(StatusBadRequest, gin.H{"error": "Course or Score  is wrong"})
	}

	if _, err := services.AddGrade(id, grade.Course, grade.Score); err != nil {
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(StatusOK, gin.H{"message": "grade added or updated successfully"})
}

//

func UpdateCourseGrade(c *gin.Context) {
	id := c.Param("id")
	course := c.Param("course")
	var grade models.Grade
	if err := c.ShouldBindJSON(&grade); err != nil {
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateGrade(id, course, grade.Score); err != nil {
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(StatusOK, gin.H{"message": "grade updated successfully"})
}

// 删除课程成绩

func DeleteGrade(c *gin.Context) {
	id := c.Param("id")
	course := c.Query("course")

	if err := services.DeleteGrade(id, course); err != nil {
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(StatusOK, gin.H{"message": "grade deleted successfully"})
}

// 文件上传处理函数，用户上传CSV文件

func UploadCSV(c *gin.Context) {
	// 获取上传的CSV文件
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to retrieve file: %s", err.Error()),
		})
		return
	}
	defer file.Close()

	// 读取文件内容到内存
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to read file: %s", err.Error()),
		})
		return
	}

	// 使用 goroutine 异步处理 CSV 文件解析
	go func(data []byte) {
		// 创建新的 io.Reader 传递给 ParseCSV
		reader := bytes.NewReader(data)

		// 解析 CSV 文件
		students, err := utils.ParseCSV(reader)
		if err != nil {
			log.Println("Error parsing CSV:", err)
			return
		}

		// 将 CSV 中的每个学生信息添加或更新到系统中
		for _, student := range students {
			err := services.AddStudent(student)
			if err != nil {
				log.Println("Error updating student:", err)
			}
		}
	}(fileBytes)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and processing started.",
	})
}
