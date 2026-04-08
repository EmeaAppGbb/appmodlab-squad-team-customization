package patient

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	db *pgxpool.Pool
}

type Patient struct {
	ID        string    `json:"id" db:"id"`
	MRN       string    `json:"mrn" db:"mrn"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	DOB       time.Time `json:"date_of_birth" db:"date_of_birth"`
	SSN       string    `json:"ssn,omitempty" db:"ssn"` // PHI - Protected Health Information
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Address   string    `json:"address,omitempty" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) CreatePatient(c *gin.Context) {
	var patient Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// HIPAA VIOLATION: Logging PHI (patient name and SSN should not be logged)
	// This is intentional for the lab to demonstrate what the HIPAA agent should catch
	// log.Printf("Creating patient: %s %s, SSN: %s", patient.FirstName, patient.LastName, patient.SSN)

	query := `
		INSERT INTO patients (mrn, first_name, last_name, date_of_birth, ssn, email, phone, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(context.Background(), query,
		patient.MRN, patient.FirstName, patient.LastName, patient.DOB,
		patient.SSN, patient.Email, patient.Phone, patient.Address,
	).Scan(&patient.ID, &patient.CreatedAt, &patient.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	// Don't return SSN in response (HIPAA best practice - minimum necessary)
	patient.SSN = ""

	c.JSON(http.StatusCreated, patient)
}

func (s *Service) GetPatient(c *gin.Context) {
	patientID := c.Param("id")

	var patient Patient
	query := `
		SELECT id, mrn, first_name, last_name, date_of_birth, email, phone, created_at, updated_at
		FROM patients WHERE id = $1
	`

	err := s.db.QueryRow(context.Background(), query, patientID).Scan(
		&patient.ID, &patient.MRN, &patient.FirstName, &patient.LastName,
		&patient.DOB, &patient.Email, &patient.Phone,
		&patient.CreatedAt, &patient.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// AUDIT LOG: PHI access should be logged for compliance
	// TODO: Implement audit logging for HIPAA compliance

	c.JSON(http.StatusOK, patient)
}

func (s *Service) UpdatePatient(c *gin.Context) {
	patientID := c.Param("id")

	var patient Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		UPDATE patients 
		SET first_name = $1, last_name = $2, email = $3, phone = $4, address = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := s.db.QueryRow(context.Background(), query,
		patient.FirstName, patient.LastName, patient.Email, patient.Phone, patient.Address, patientID,
	).Scan(&patient.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	patient.ID = patientID
	c.JSON(http.StatusOK, patient)
}

func (s *Service) ListPatients(c *gin.Context) {
	query := `
		SELECT id, mrn, first_name, last_name, email, phone, created_at
		FROM patients
		ORDER BY last_name, first_name
		LIMIT 100
	`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list patients"})
		return
	}
	defer rows.Close()

	var patients []Patient
	for rows.Next() {
		var patient Patient
		err := rows.Scan(
			&patient.ID, &patient.MRN, &patient.FirstName, &patient.LastName,
			&patient.Email, &patient.Phone, &patient.CreatedAt,
		)
		if err != nil {
			continue
		}
		patients = append(patients, patient)
	}

	c.JSON(http.StatusOK, patients)
}
