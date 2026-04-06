package api

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}
type Staff struct {
	StaffID    int    `json:"staff_id"`
	HospitalID int    `json:"hospital_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

type CreateStaffRequest struct {
	HospitalID   int    `json:"hospital_id"    binding:"required"`
	FirstNameTH  string `json:"first_name_th"  binding:"required"`
	MiddleNameTH string `json:"middle_name_th"` // no binding
	LastNameTH   string `json:"last_name_th"   binding:"required"`
	FirstNameEN  string `json:"first_name_en"  binding:"required"`
	MiddleNameEN string `json:"middle_name_en"` // no binding
	LastNameEN   string `json:"last_name_en"   binding:"required"`
	Email        string `json:"email"          binding:"required,email"`
	Password     string `json:"password"       binding:"required,min=8"`
	Role         string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Staff Staff  `json:"staff"`
}

type Patient struct {
	NationalID   pgtype.Text `json:"national_id"`
	PassportID   pgtype.Text `json:"passport_id"`
	FirstNameTH  pgtype.Text `json:"first_name_th"`
	MiddleNameTH pgtype.Text `json:"middle_name_th"`
	LastNameTH   pgtype.Text `json:"last_name_th"`
	FirstNameEN  pgtype.Text `json:"first_name_en"`
	MiddleNameEN pgtype.Text `json:"middle_name_en"`
	LastNameEN   pgtype.Text `json:"last_name_en"`
	PatientHN    string      `json:"patient_hn"`
	DateOfBirth  pgtype.Date `json:"date_of_birth"`
	PhoneNumber  pgtype.Text `json:"phone_number"`
	Email        pgtype.Text `json:"email"`
	Gender       pgtype.Text `json:"gender"`
	HospitalID   int         `json:"hospital_id"`
}

type PatientResponse struct {
	NationalID   *string `json:"national_id"`
	PassportID   *string `json:"passport_id"`
	FirstNameTH  *string `json:"first_name_th"`
	MiddleNameTH *string `json:"middle_name_th"`
	LastNameTH   *string `json:"last_name_th"`
	FirstNameEN  *string `json:"first_name_en"`
	MiddleNameEN *string `json:"middle_name_en"`
	LastNameEN   *string `json:"last_name_en"`
	DateOfBirth  *string `json:"date_of_birth"`
	PatientHN    *string `json:"patient_hn"`
	PhoneNumber  *string `json:"phone_number"`
	Email        *string `json:"email"`
	Gender       *string `json:"gender"`
	HospitalID   int     `json:"hospital_id"`
}

type Claims struct {
	StaffID    int    `json:"staff_id"`
	HospitalID int    `json:"hospital_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}
