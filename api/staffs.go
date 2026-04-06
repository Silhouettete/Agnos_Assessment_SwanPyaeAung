package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateStaff(pool DBPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateStaffRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		role := req.Role
		if role == "" {
			role = "staff"
		}

		var staffID int
		err = pool.QueryRow(context.Background(), `
  INSERT INTO staff (
    hospital_id,
    first_name_th, middle_name_th, last_name_th,
    first_name_en, middle_name_en, last_name_en,
    email, password_hash, role
  )
  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
  RETURNING staff_id`,
			req.HospitalID,
			req.FirstNameTH, req.MiddleNameTH, req.LastNameTH,
			req.FirstNameEN, req.MiddleNameEN, req.LastNameEN,
			req.Email,
			string(hash),
			role,
		).Scan(&staffID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists or invalid hospital"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"staff_id":    staffID,
			"hospital_id": req.HospitalID,
			"email":       req.Email,
			"role":        role,
		})
	}
}

func Login(pool DBPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var s Staff
		var passwordHash string
		err := pool.QueryRow(context.Background(), `
  SELECT 
    staff_id, hospital_id,
    first_name_en, last_name_en,
    email, role, password_hash
  FROM staff
  WHERE email = $1`,
			req.Email,
		).Scan(
			&s.StaffID,
			&s.HospitalID,
			&s.FirstName,
			&s.LastName,
			&s.Email,
			&s.Role,
			&passwordHash,
		)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		claims := &Claims{
			StaffID:    s.StaffID,
			HospitalID: s.HospitalID,
			Email:      s.Email,
			Role:       s.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenStr,
			"staff": s,
		})
	}
}
