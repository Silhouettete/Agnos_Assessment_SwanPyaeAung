package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"agnos_assessment/api"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateStaff_Success(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	mock.ExpectQuery(`INSERT INTO staff`).
		WithArgs(
			1,
			"JaneTH", "", "SmithTH",
			"Jane", "", "Smith",
			"jane@hospital.com",
			pgxmock.AnyArg(),
			"staff",
		).
		WillReturnRows(pgxmock.NewRows([]string{"staff_id"}).AddRow(1))

	r := setupRouter()
	r.POST("/staff/create", api.CreateStaff(mock))

	req := makeRequest("POST", "/staff/create", map[string]interface{}{
		"hospital_id":    1,
		"first_name_th":  "JaneTH",
		"middle_name_th": "",
		"last_name_th":   "SmithTH",
		"first_name_en":  "Jane",
		"middle_name_en": "",
		"last_name_en":   "Smith",
		"email":          "jane@hospital.com",
		"password":       "secret123",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "jane@hospital.com", resp["email"])
	assert.Equal(t, "staff", resp["role"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateStaff_MissingFields(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	r := setupRouter()
	r.POST("/staff/create", api.CreateStaff(mock))

	req := makeRequest("POST", "/staff/create", map[string]interface{}{
		"email": "jane@hospital.com",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["error"])
}

func TestCreateStaff_InvalidEmail(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	r := setupRouter()
	r.POST("/staff/create", api.CreateStaff(mock))

	req := makeRequest("POST", "/staff/create", map[string]interface{}{
		"hospital_id":   1,
		"first_name_th": "JaneTH",
		"last_name_th":  "SmithTH",
		"first_name_en": "Jane",
		"last_name_en":  "Smith",
		"email":         "not-an-email",
		"password":      "secret123",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStaff_PasswordTooShort(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	r := setupRouter()
	r.POST("/staff/create", api.CreateStaff(mock))

	req := makeRequest("POST", "/staff/create", map[string]interface{}{
		"hospital_id":   1,
		"first_name_th": "JaneTH",
		"last_name_th":  "SmithTH",
		"first_name_en": "Jane",
		"last_name_en":  "Smith",
		"email":         "jane@hospital.com",
		"password":      "123",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateStaff_DuplicateEmail(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	mock.ExpectQuery(`INSERT INTO staff`).
		WithArgs(
			1,
			"JaneTH", "", "SmithTH",
			"Jane", "", "Smith",
			"jane@hospital.com",
			pgxmock.AnyArg(),
			"staff",
		).
		WillReturnError(fmt.Errorf("ERROR: duplicate key value violates unique constraint"))

	r := setupRouter()
	r.POST("/staff/create", api.CreateStaff(mock))

	req := makeRequest("POST", "/staff/create", map[string]interface{}{
		"hospital_id":    1,
		"first_name_th":  "JaneTH",
		"middle_name_th": "",
		"last_name_th":   "SmithTH",
		"first_name_en":  "Jane",
		"middle_name_en": "",
		"last_name_en":   "Smith",
		"email":          "jane@hospital.com",
		"password":       "secret123",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_Success(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)

	mock.ExpectQuery(`SELECT staff_id, hospital_id, first_name_en, last_name_en, email, role, password_hash`).
		WithArgs("jane@hospital.com").
		WillReturnRows(pgxmock.NewRows([]string{
			"staff_id",
			"hospital_id",
			"first_name_en",
			"last_name_en",
			"email",
			"role",
			"password_hash",
		}).AddRow(
			1,
			1,
			"Jane",
			"Smith",
			"jane@hospital.com",
			"staff",
			string(hash),
		))

	r := setupRouter()
	r.POST("/staff/login", api.Login(mock))

	req := makeRequest("POST", "/staff/login", map[string]interface{}{
		"email":    "jane@hospital.com",
		"password": "secret123",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
	assert.NotNil(t, resp["staff"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_WrongPassword(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)

	mock.ExpectQuery(`SELECT staff_id, hospital_id, first_name_en, last_name_en, email, role, password_hash`).
		WithArgs("jane@hospital.com").
		WillReturnRows(pgxmock.NewRows([]string{
			"staff_id",
			"hospital_id",
			"first_name_en",
			"last_name_en",
			"email",
			"role",
			"password_hash",
		}).AddRow(
			1,
			1,
			"Jane",
			"Smith",
			"jane@hospital.com",
			"staff",
			string(hash),
		))

	r := setupRouter()
	r.POST("/staff/login", api.Login(mock))

	req := makeRequest("POST", "/staff/login", map[string]interface{}{
		"email":    "jane@hospital.com",
		"password": "wrongpassword",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid email or password", resp["error"])
}

func TestLogin_EmailNotFound(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	mock.ExpectQuery(`SELECT staff_id, hospital_id, first_name_en, last_name_en, email, role, password_hash`).
		WithArgs("nobody@hospital.com").
		WillReturnError(pgx.ErrNoRows)

	r := setupRouter()
	r.POST("/staff/login", api.Login(mock))

	req := makeRequest("POST", "/staff/login", map[string]interface{}{
		"email":    "nobody@hospital.com",
		"password": "secret123",
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "invalid email or password", resp["error"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLogin_MissingBody(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	r := setupRouter()
	r.POST("/staff/login", api.Login(mock))

	req := makeRequest("POST", "/staff/login", map[string]interface{}{})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
