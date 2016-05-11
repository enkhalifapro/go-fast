package viewModels

type LoginViewModel struct {
	Email     string        `form:"email" json:"email" binding:"required"`
	Password  string        `form:"password" json:"password" binding:"required"`
	Token     string
	UserId    string
	UserName  string
	FirstName string
	LastName  string
	Image     string
}