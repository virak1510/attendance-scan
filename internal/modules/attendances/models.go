package attendances

type AttendanceQuery struct {
	FirstName string  `json:"first_name" db:"first_name"`
	LastName  string  `json:"last_name" db:"last_name"`
	Age       int     `json:"age" db:"age"`
	TimeIn    string  `json:"time_in" db:"time_in"`
	TimeOut   *string `json:"time_out" db:"time_out"`
	Status    *Status `json:"status" db:"status"`
	Notes     *string `json:"notes" db:"notes"`
}

type AttendanceInsert struct {
	UserID  int     `json:"user_id" db:"user_id"`
	Date    string  `json:"date" db:"date"`
	TimeIn  string  `json:"time_in" db:"time_in"`
	TimeOut *string `json:"time_out" db:"time_out"`
	Status  Status  `json:"status" db:"status"`
	Notes   *string `json:"notes" db:"notes"`
}

type Status string

const (
	Present Status = "present"
	Absent  Status = "absent"
	Late    Status = "late"
	Sick    Status = "sick"
)

func (s Status) String() string {
	return string(s)
}

type AttendanceParams struct {
	FirstName string  `json:"first_name" binding:"required"`
	LastName  string  `json:"last_name" binding:"required"`
	Age       int     `json:"age" binding:"required"`
	TimeIn    string  `json:"time_in" binding:"required"`
	TimeOut   string  `json:"time_out" binding:"required"`
	Status    *string `json:"status" `
	Notes     *string `json:"notes" `
}
