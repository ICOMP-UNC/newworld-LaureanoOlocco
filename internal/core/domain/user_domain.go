package domain

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

/*
func NewPerson(id int, email string, password string){
  return &User{
	ID:id,
	Email:email,
	Password: password,
  }
}

func (u *User) GetEmail() string {
  return u.Email
}
*/

type Register struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"example@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
	//Role     string `json:"role" validate:"required" example:"admin"`
}

type Login struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
	//Role     string `json:"role" validate:"required" example:"admin"`
}
