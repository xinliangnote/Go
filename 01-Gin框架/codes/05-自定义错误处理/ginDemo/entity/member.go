package entity

// 定义 Member 结构体
type Member struct {
	Name string `form:"name" json:"name" binding:"required,NameValid"`
	Age  int    `form:"age"  json:"age"  binding:"required,gt=10,lt=120"`
}
