package models

type Student struct {
	Name   string  `json:"name"`
	ID     string  `json:"id"`
	Gender string  `json:"gender"`
	Class  string  `json:"class"`
	Grades []Grade `json:"grades"` // 改为切片类型
}

type Grade struct {
	Course string  `json:"course"` // 课程名称
	Score  float64 `json:"score"`  // 成绩
}
