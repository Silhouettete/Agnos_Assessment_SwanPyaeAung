package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agnos_assessment/api"

	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func withHospitalID(hospitalID int, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("hospital_id", hospitalID)
		handler(c)
	}
}

func patientColumns() []string {
	return []string{
		"patient_hn", "national_id", "passport_id",
		"first_name_th", "middle_name_th", "last_name_th",
		"first_name_en", "middle_name_en", "last_name_en",
		"date_of_birth", "phone_number", "email", "gender",
		"hospital_id",
	}
}

func TestSearchPatient_ByNationalID_Success(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	natID := "1234567890123"
	passID := "PP123456"
	firstName := "Somchai"
	lastName := "Jaidee"
	phone := "0812345678"
	email := "somchai@email.com"
	gender := "Male"

	mock.ExpectQuery(`SELECT`).
		WithArgs(1, natID, "").
		WillReturnRows(pgxmock.NewRows(patientColumns()).AddRow(
			"HN001", natID, passID,
			nil, nil, nil,
			firstName, nil, lastName,
			dob, phone, email, gender,
			1,
		))

	r := setupRouter()
	r.GET("/patient/search", withHospitalID(1, api.SearchPatient(mock)))

	req, _ := http.NewRequest("GET", "/patient/search?national_id=1234567890123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var results []api.PatientResponse
	json.Unmarshal(w.Body.Bytes(), &results)

	if assert.Len(t, results, 1) {
		assert.Equal(t, "HN001", *results[0].PatientHN)
		assert.Equal(t, "1234567890123", *results[0].NationalID)
		assert.Equal(t, "Somchai", *results[0].FirstNameEN)
		assert.Equal(t, "1990-01-15", *results[0].DateOfBirth)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPatient_ByPassportID_Success(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	dob := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	passID := "PP123456"
	firstName := "John"
	lastName := "Doe"

	mock.ExpectQuery(`SELECT`).
		WithArgs(1, "", passID).
		WillReturnRows(pgxmock.NewRows(patientColumns()).AddRow(
			"HN002", nil, passID,
			nil, nil, nil,
			firstName, nil, lastName,
			dob, nil, nil, nil,
			1,
		))

	r := setupRouter()
	r.GET("/patient/search", withHospitalID(1, api.SearchPatient(mock)))

	req, _ := http.NewRequest("GET", "/patient/search?passport_id=PP123456", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var results []api.PatientResponse
	json.Unmarshal(w.Body.Bytes(), &results)

	if assert.Len(t, results, 1) {
		assert.Equal(t, "HN002", *results[0].PatientHN)
		assert.Equal(t, "PP123456", *results[0].PassportID)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPatient_NoQueryParams(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	r := setupRouter()
	r.GET("/patient/search", withHospitalID(1, api.SearchPatient(mock)))

	req, _ := http.NewRequest("GET", "/patient/search", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "provide national_id or passport_id", resp["error"])
}

func TestSearchPatient_NotFound(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	natID := "0000000000000"

	mock.ExpectQuery(`SELECT`).
		WithArgs(1, natID, "").
		WillReturnRows(pgxmock.NewRows(patientColumns()))

	r := setupRouter()
	r.GET("/patient/search", withHospitalID(1, api.SearchPatient(mock)))

	req, _ := http.NewRequest("GET", "/patient/search?national_id=0000000000000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "no patients found", resp["error"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPatient_WrongHospital(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	natID := "1234567890123"

	mock.ExpectQuery(`SELECT`).
		WithArgs(2, natID, "").
		WillReturnRows(pgxmock.NewRows(patientColumns()))

	r := setupRouter()
	r.GET("/patient/search", withHospitalID(2, api.SearchPatient(mock)))

	req, _ := http.NewRequest("GET", "/patient/search?national_id=1234567890123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "no patients found", resp["error"])

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchPatient_DatabaseError(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	natID := "1234567890123"

	mock.ExpectQuery(`SELECT`).
		WithArgs(1, natID, "").
		WillReturnError(fmt.Errorf("connection refused"))

	r := setupRouter()
	r.GET("/patient/search", withHospitalID(1, api.SearchPatient(mock)))

	req, _ := http.NewRequest("GET", "/patient/search?national_id=1234567890123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
