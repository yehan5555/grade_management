本项目基于一个 Gin 框架搭建的一个 简单的成绩管理系统
==

该项目实现了学生信息和成绩的录入，删除，修改和查询功能
-
对于数据存储，并未持久化存储到数据库中，数据存储在程序
本身的内存中，使用map 来保存学生信息和成绩信息

并发处理
-
从 CSV 文件中批量导入学生信息到系统中，用户可以通过接口上传 CSV 文件。
使用 goroutines 并发处理每一条学生信息，通过 channel 传递解析后的学
生信息，确保数据在 goroutines 之间安全传递。使用互斥锁保证原子性。处
理了文件格式错误、数据验证错误，可以在解析 CSV 文件时进行数据验证，并
在发现错误时记录错误信息或跳过错误数据。

异常处理
-
添加异常处理机制，确保系统在遇到错误时能够返回适当的错误信息，捕获并处
理了常见的错误，如输入数据格式错误、数据验证错误。使用捕获并处理常见的
错误，如输入数据格式错误、数据验证错误，比如对于不存在的学生或成绩信息
，API应返回适当的错误消息和状态码（如404 Not Found），对于无效的输入
数据，API返回适当的错误消息和状态码（如400 Bad Request）。

单元测试
-
使用 testing 包对学生信息和成绩的增删改查进行了单元测试


接口文档
-
在postman 中导入 doc 文件夹下的grade_management.postman_collection.json，即可测试。
url 设置为：8080。 

学生信息录入  
![学生信息录入](https://github.com/yehan5555/grade_management/blob/master/doc/luru.png)

学生成绩录入  
![学生成绩录入](https://github.com/yehan5555/grade_management/blob/master/doc/lurucj.png)

学生信息查询  
![学生信息查询](https://github.com/yehan5555/grade_management/blob/master/doc/search.png)

学生信息查询（导入的csv文件，查询信息）
![学生信息csv查询](https://github.com/yehan5555/grade_management/blob/master/doc/search_csv.png)












