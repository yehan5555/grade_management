package services

import (
	"Grade_managing/models"
	"errors"
	"fmt"
	"sync"
)

var (
	Students = make(map[string]models.Student)
	Lock     sync.RWMutex
)

// 添加或更新学生信息

func AddStudent(student models.Student) (err error) {
	Lock.Lock()
	defer Lock.Unlock()
	if student.ID == "" {
		return fmt.Errorf("invalid student ID")
	}
	if student.Name == "" {
		return fmt.Errorf("invalid student name")
	}
	if student.Gender == "" {
		return fmt.Errorf("invalid student gender")
	}
	if student.Class == "" {
		return fmt.Errorf("invalid student class")
	}
	Students[student.ID] = student
	return nil
}

// 根据学号获取学生信息

func GetStudentByID(id string) (models.Student, error) {
	Lock.Lock()
	defer Lock.Unlock()
	student, exists := Students[id]
	if !exists {
		return models.Student{}, errors.New("student not found")
	}
	return student, nil
}

// 获取所有学生信息

func GetAllStudents() []models.Student {
	Lock.Lock()
	defer Lock.Unlock()
	var allStudents []models.Student
	for _, student := range Students {
		allStudents = append(allStudents, student)
	}
	return allStudents
}

// 删除学生信息（包括成绩）

func DeleteStudentByID(id string) error {
	Lock.Lock()
	defer Lock.Unlock()
	if _, exists := Students[id]; !exists {
		return errors.New("student not found")
	}
	delete(Students, id)
	return nil
}

// 修改学生信息（不包括成绩）

func UpdateStudentInfo(id string, newInfo models.Student) error {
	Lock.Lock()
	defer Lock.Unlock()
	student, exists := Students[id]
	if !exists {
		return fmt.Errorf("student not found")
	}
	if student.ID == "" {
		return fmt.Errorf("invalid student ID")
	}
	if student.Name == "" {
		return fmt.Errorf("invalid student name")
	}
	if student.Gender == "" {
		return fmt.Errorf("invalid student gender")
	}
	if student.Class == "" {
		return fmt.Errorf("invalid student class")
	}

	// 更新基本信息
	student.Name = newInfo.Name
	student.Gender = newInfo.Gender
	student.Class = newInfo.Class
	Students[id] = student

	return nil
}

// 根据学号和课程名称添加成绩

func AddGrade(id string, course string, score float64) (student models.Student, err error) {
	Lock.Lock()
	defer Lock.Unlock()
	student, exists := Students[id]
	if !exists {
		return models.Student{}, errors.New("student not found")
	}

	// 查找课程是否已存在
	updated := false
	for i, grade := range student.Grades {
		if grade.Course == course {
			// 更新已有课程的成绩
			student.Grades[i].Score = score
			updated = true
			break
		}
	}

	// 如果课程不存在，追加新课程成绩
	if !updated {
		newGrade := models.Grade{
			Course: course,
			Score:  score,
		}
		student.Grades = append(student.Grades, newGrade)
	}

	// 更新学生信息
	Students[id] = student
	return student, nil
}

func UpdateGrade(id string, course string, score float64) error {
	Lock.Lock()
	defer Lock.Unlock()
	student, exists := Students[id]
	if !exists {
		return fmt.Errorf("student with ID %s not found", id)
	}
	if score < 0 || score > 100 {
		return fmt.Errorf("invalid score")
	}

	// 查找课程是否存在
	for i, grade := range student.Grades {
		if grade.Course == course {
			// 更新课程成绩
			student.Grades[i].Score = score
			Students[id] = student
			return nil
		}
	}

	return fmt.Errorf("course not found")
}

// 根据学号和课程名称删除成绩

func DeleteGrade(id string, course string) error {
	Lock.Lock()
	defer Lock.Unlock()
	student, exists := Students[id]
	if !exists {
		return errors.New("student not found")
	}

	// 过滤掉对应课程成绩
	newGrades := []models.Grade{}
	for _, g := range student.Grades {
		if g.Course != course {
			newGrades = append(newGrades, g)
		}
	}

	if len(newGrades) == len(student.Grades) {
		return errors.New("course not found")
	}

	student.Grades = newGrades
	Students[id] = student
	return nil
}
