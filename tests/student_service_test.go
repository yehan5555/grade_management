package tests

import (
	"Grade_managing/models"
	"Grade_managing/services"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestStudentService(t *testing.T) {
	// 测试添加学生
	t.Run("Add Student", func(t *testing.T) {
		student := models.Student{
			ID:     "123",
			Name:   "ysj",
			Gender: "Male",
			Class:  "1A",
			Grades: []models.Grade{
				{Course: "Math", Score: 95},
			},
		}
		err := services.AddStudent(student)
		assert.NoError(t, err) // 确保没有错误
		log.Println("学生信息录入测试成功:", student)

	})
	// 测试添加失败
	t.Run("Add Student with duplicate ID", func(t *testing.T) {
		// 添加第一个学生
		student1 := models.Student{
			ID:     "",
			Name:   "ysj",
			Gender: "Male",
			Class:  "1A",
			Grades: []models.Grade{
				{Course: "Math", Score: 95},
			},
		}
		err := services.AddStudent(student1)
		// 断言返回错误
		assert.Error(t, err)
		// 断言错误内容
		assert.EqualError(t, err, "student with ID 123 already exists")
	})

	// 测试添加学生成绩
	t.Run("Add Grade", func(t *testing.T) {
		student, err := services.AddGrade("123", "English", 90.0)
		assert.NoError(t, err)
		assert.Equal(t, float64(90), student.Grades[1].Score)
		log.Println("学生成绩录入测试成功:")
	})
	// 测试添加学生成绩失败
	t.Run("Add Grade Failed", func(t *testing.T) {
		student, err := services.AddGrade("123", "Math", 110.0)
		assert.Error(t, err) // 确保有错误
		assert.Equal(t, 110.0, student.Grades[0].Score)
		log.Println("学生成绩录入失败测试成功:")
	})
	// 测试查询学生个人信息
	t.Run("Get Student", func(t *testing.T) {
		student, err := services.GetStudentByID("123")
		assert.NoError(t, err) // 确保没有错误
		assert.Equal(t, "ysj", student.Name)
		assert.Equal(t, "Male", student.Gender)
		log.Println("学生信息查询测试成功:")
	})
	// 测试查询学生个人信息失败
	t.Run("Get Student Failed", func(t *testing.T) {
		student, err := services.GetStudentByID("122")
		assert.Error(t, err) // 确保有错误
		assert.Equal(t, "", student.Name)
		assert.Equal(t, "", student.Gender)
		log.Println("学生信息查询失败测试成功:")
	})

	//测试查询所有学生信息
	t.Run("Get All Students", func(t *testing.T) {
		students := services.GetAllStudents()
		assert.Equal(t, 1, len(students))
		log.Println("所有学生信息查询测试成功:")
	})

	// 测试更新学生信息,根据学号修改学生基本信息
	t.Run("Update Student", func(t *testing.T) {
		updatedStudent := models.Student{
			ID:     "123",
			Name:   "YeHanHan",
			Gender: "Male",
			Class:  "1A",
		}
		err := services.UpdateStudentInfo(updatedStudent.ID, updatedStudent)
		assert.NoError(t, err)

		student, err := services.AddGrade("123", "Math", 99.0)
		assert.NoError(t, err)
		assert.Equal(t, 99.0, student.Grades[0].Score)

		log.Println("学生信息更新测试成功:")
	})

	// 测试更新学生信息失败
	t.Run("Update Student Failed", func(t *testing.T) {
		err := services.UpdateStudentInfo("123", models.Student{
			ID:     "123",
			Name:   "",
			Gender: "Male",
			Class:  "1A",
		})
		assert.NoError(t, err)
		log.Println("学生信息更新失败测试成功:")
	})

	t.Run("Update StudentScores Failed", func(t *testing.T) {

		err := services.UpdateGrade("123", "Math", 110)
		assert.Error(t, err) // 查询已经删除的学生应返回错误
		log.Println("学生成绩更新失败测试成功:")
	})

	//测试删除学生成绩
	t.Run("Delete Grade", func(t *testing.T) {
		err := services.DeleteGrade("123", "Math")
		assert.NoError(t, err)

		student, err := services.GetStudentByID("123")
		assert.NoError(t, err)
		// 应为有两个成绩
		assert.Equal(t, 1, len(student.Grades))
		log.Println("学生成绩删除测试成功:")
	})

	// 测试删除学生成绩失败
	t.Run("Delete Grade Failed", func(t *testing.T) {
		err := services.DeleteGrade("122", "Math")
		assert.Error(t, err) // 查询已经删除的学生应返回错误
		log.Println("学生成绩删除失败测试成功:")
	})

	// 测试删除学生
	t.Run("Delete Student", func(t *testing.T) {
		err := services.DeleteStudentByID("123")
		assert.NoError(t, err)

		_, err = services.GetStudentByID("123")
		assert.Error(t, err) // 查询已经删除的学生应返回错误
		log.Println("学生删除测试成功:")
	})
	// 测试删除学生失败
	t.Run("Delete Student Failed", func(t *testing.T) {
		err := services.DeleteStudentByID("122")
		assert.Error(t, err) // 查询已经删除的学生应返回错误
		log.Println("学生删除失败测试成功:")
	})

}
