package dao

type DaoUser struct {
	Id       int64
	Name     string
	Email    string
	Phone    string
	Passwd   string
	ShowName string
}

func (this *DaoUser) TableName() string {
	return "t_user"
}
