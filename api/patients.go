package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func SearchPatient(pool DBPool) gin.HandlerFunc {
	return func(c *gin.Context) {
		hospitalID, _ := c.Get("hospital_id")
		nationalID := c.Query("national_id")
		passportID := c.Query("passport_id")

		if nationalID == "" && passportID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "provide national_id or passport_id"})
			return
		}

		rows, err := pool.Query(context.Background(), `
			SELECT
				patient_hn, national_id, passport_id,
				first_name_th, middle_name_th, last_name_th,
				first_name_en, middle_name_en, last_name_en,
				date_of_birth, phone_number, email, gender,
				hospital_id
			FROM patients
			WHERE hospital_id = $1
			  AND (
			      ($2 != '' AND national_id = $2)
			      OR
			      ($3 != '' AND passport_id = $3)
			  )`,
			hospitalID, nationalID, passportID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var results []PatientResponse
		for rows.Next() {
			var p Patient
			if err := rows.Scan(
				&p.PatientHN, &p.NationalID, &p.PassportID,
				&p.FirstNameTH, &p.MiddleNameTH, &p.LastNameTH,
				&p.FirstNameEN, &p.MiddleNameEN, &p.LastNameEN,
				&p.DateOfBirth, &p.PhoneNumber, &p.Email, &p.Gender,
				&p.HospitalID,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			results = append(results, toResponse(p))
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(results) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "no patients found"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func toResponse(p Patient) PatientResponse {
	return PatientResponse{
		PatientHN:    &p.PatientHN,
		NationalID:   textVal(p.NationalID),
		PassportID:   textVal(p.PassportID),
		FirstNameTH:  textVal(p.FirstNameTH),
		MiddleNameTH: textVal(p.MiddleNameTH),
		LastNameTH:   textVal(p.LastNameTH),
		FirstNameEN:  textVal(p.FirstNameEN),
		MiddleNameEN: textVal(p.MiddleNameEN),
		LastNameEN:   textVal(p.LastNameEN),
		DateOfBirth:  dateVal(p.DateOfBirth),
		PhoneNumber:  textVal(p.PhoneNumber),
		Email:        textVal(p.Email),
		Gender:       textVal(p.Gender),
		HospitalID:   p.HospitalID,
	}
}

func textVal(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func dateVal(d pgtype.Date) *string {
	if !d.Valid {
		return nil
	}
	s := d.Time.Format("2006-01-02")
	return &s
}
