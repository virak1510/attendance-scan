package attendances

import (
	"attendance/pkg"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	sqlx *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		sqlx: db,
	}
}

func (a *Service) GetAttendances(ctx context.Context) ([]AttendanceQuery, error) {
	attendances := []AttendanceQuery{}
	sqlStr := `
	SELECT u.first_name, u.last_name, u.age, a.time_in, a.time_out, a.status, a.notes 
	FROM tbl_users u JOIN tbl_attendances a ON u.id = a.user_id 
	WHERE DATE(a.created_at) BETWEEN CURRENT_DATE AND CURRENT_DATE
	`
	err := a.sqlx.Select(&attendances, sqlStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(attendances)
	return attendances, nil
}

func (a *Service) CheckIn(ctx context.Context, user *pkg.AuthUser) error {
	attendance := AttendanceInsert{
		UserID:  user.ID,
		Date:    time.Now().Format("2006-01-02"),
		TimeIn:  time.Now().Format("15:04:05"),
		TimeOut: nil,
		Status:  Present,
		Notes:   nil,
	}
	sqlQueryExist := `
			SELECT COUNT(*) 
			FROM tbl_attendances
			WHERE user_id = $1 AND date = $2;
			`

	var count int
	err := a.sqlx.QueryRow(sqlQueryExist, user.ID, time.Now().Format("2006-01-02")).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("you're already checked in")
	}

	sqlStr := `
	INSERT INTO tbl_attendances (user_id, date, time_in, time_out, status, notes)
	VALUES (:user_id, :date, :time_in, :time_out, :status, :notes)
	`
	_, err = a.sqlx.NamedExec(sqlStr, attendance)
	if err != nil {
		return err
	}
	return nil
}

func (a *Service) CheckOut(ctx context.Context, user *pkg.AuthUser) error {
	now := time.Now().Format("15:04:05")
	attendance := AttendanceInsert{
		UserID:  user.ID,
		Date:    time.Now().Format("2006-01-02"),
		TimeIn:  time.Now().Format("15:04:05"),
		TimeOut: &now,
		Status:  Present,
		Notes:   nil,
	}
	sqlStr := `
	UPDATE tbl_attendances SET time_out = :time_out WHERE user_id = :user_id AND date = :date AND time_out IS NULL
	`
	_, err := a.sqlx.NamedExec(sqlStr, attendance)
	if err != nil {
		return err
	}
	return nil
}
