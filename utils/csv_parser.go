package utils

import (
	"Grade_managing/models"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"strconv"
	"sync"
)

// ParseCSV 解析 CSV 文件，并返回学生信息切片
func ParseCSV(reader io.Reader) ([]models.Student, error) {
	csvReader := csv.NewReader(reader)
	var wg sync.WaitGroup                          // 用于等待所有 goroutine 完成
	var mu sync.Mutex                              // 保护共享资源
	studentsMap := make(map[string]models.Student) // 存储学生信息的临时 map
	errorsChan := make(chan error, 1000)           // 收集解析过程中的错误
	resultsChan := make(chan models.Student, 1000) // 收集成功解析的学生

	// goroutine 用于从 channels 里处理结果并构建最终的结果集
	go func() {
		defer close(errorsChan)
		defer close(resultsChan)

		for {
			record, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorsChan <- errors.New("error reading CSV file")
				continue
			}

			// 校验字段数目
			if len(record) < 4 {
				errorsChan <- errors.New("invalid record length")
				continue
			}

			// 解析学生基本信息
			name, id, gender, class := record[0], record[1], record[2], record[3]

			// 提取字段值
			if id == "" {
				log.Printf("跳过学生 ID 为空的行")
				continue // 跳过错误行
			}

			if name == "" {
				log.Printf("跳过学生姓名为空的行")
				continue // 跳过错误行
			}

			if gender == "" {
				log.Printf("跳过学生性别为空的行")
				continue // 跳过错误行
			}

			if class == "" {
				log.Printf("跳过学生班级为空的行")
				continue // 跳过错误行
			}

			grades := []models.Grade{}
			for i := 4; i < len(record); i += 2 {
				if i+1 < len(record) {
					course := record[i]
					score, err := strconv.ParseFloat(record[i+1], 64)
					if err != nil {
						log.Printf("Error parsing grade for %s: %v\n", id, err)
						continue
					}
					grades = append(grades, models.Grade{
						Course: course,
						Score:  score,
					})
				}
			}

			// 创建学生实例
			student := models.Student{
				Name:   name,
				ID:     id,
				Gender: gender,
				Class:  class,
				Grades: grades,
			}

			resultsChan <- student
		}
	}()

	// 用多个 goroutine 并发处理每个学生信息
	for student := range resultsChan {
		wg.Add(1)
		go func(student models.Student) {
			defer wg.Done()

			// 原子性更新逻辑
			mu.Lock()
			existing, found := studentsMap[student.ID]
			if found {
				// 更新已有信息
				existing.Name = student.Name
				existing.Gender = student.Gender
				existing.Class = student.Class
				// 合并成绩
				for _, newGrade := range student.Grades {
					updated := false
					for i, existingGrade := range existing.Grades {
						if existingGrade.Course == newGrade.Course {
							existing.Grades[i] = newGrade
							updated = true
							break
						}
					}
					if !updated {
						existing.Grades = append(existing.Grades, newGrade)
					}
				}
				studentsMap[student.ID] = existing
			} else {
				// 新增学生
				studentsMap[student.ID] = student
			}
			mu.Unlock()
		}(student)
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 将 studentsMap 转换为切片返回
	var students []models.Student
	for _, student := range studentsMap {
		students = append(students, student)
	}

	// 如果有错误记录，返回错误信息（以日志的形式记录）
	if len(errorsChan) > 0 {
		log.Println("Some errors occurred during CSV parsing:")
		for err := range errorsChan {
			log.Println(err)
		}
	}

	return students, nil
}
