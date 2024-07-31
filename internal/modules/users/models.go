package users

type CreateUserParams struct {
	Username  string `json:"username" validate:"required,min=3,max=20"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Age       int32  `json:"age" validate:"required,gte=0"`
	Password  string `json:"password" validate:"required"`
}

type UserQuery struct {
	Username  string  `json:"username" db:"username"`
	FirstName *string `json:"first_name" db:"first_name"`
	LastName  *string `json:"last_name" db:"last_name"`
	Age       int32   `json:"age" db:"age"`
	Password  string  `json:"-" db:"password"`
}

type LoginParams struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginPayload struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	Token    string `json:"token"`
}
