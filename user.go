package todo

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name" binding:"name"`
	Username string `json:"username" binding:"username"`
	Password string `json:"password" binding:"password"`
}
