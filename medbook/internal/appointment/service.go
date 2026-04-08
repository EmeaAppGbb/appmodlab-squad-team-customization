package appointment

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

type Appointment struct {
	ID              string    `json:"id" db:"id"`
	PatientID       string    `json:"patient_id" db:"patient_id"`
	ProviderID      string    `json:"provider_id" db:"provider_id"`
	AppointmentTime time.Time `json:"appointment_time" db:"appointment_time"`
	DiagnosisCode   string    `json:"diagnosis_code,omitempty" db:"diagnosis_code"` // ICD-10 code
	ProcedureCode   string    `json:"procedure_code,omitempty" db:"procedure_code"` // CPT code
	Status          string    `json:"status" db:"status"`                           // scheduled, completed, cancelled
	Notes           string    `json:"notes,omitempty" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) CreateAppointment(c *gin.Context) {
	var appt Appointment
	if err := c.ShouldBindJSON(&appt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate ICD-10 diagnosis code format (Domain Terminology Agent should check this)
	// TODO: Validate CPT procedure code format

	query := `
		INSERT INTO appointments (patient_id, provider_id, appointment_time, diagnosis_code, procedure_code, status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(context.Background(), query,
		appt.PatientID, appt.ProviderID, appt.AppointmentTime,
		appt.DiagnosisCode, appt.ProcedureCode, "scheduled", appt.Notes,
	).Scan(&appt.ID, &appt.CreatedAt, &appt.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	appt.Status = "scheduled"
	c.JSON(http.StatusCreated, appt)
}

func (s *Service) GetAppointment(c *gin.Context) {
	apptID := c.Param("id")

	var appt Appointment
	query := `
		SELECT id, patient_id, provider_id, appointment_time, diagnosis_code, 
		       procedure_code, status, notes, created_at, updated_at
		FROM appointments WHERE id = $1
	`

	err := s.db.QueryRow(context.Background(), query, apptID).Scan(
		&appt.ID, &appt.PatientID, &appt.ProviderID, &appt.AppointmentTime,
		&appt.DiagnosisCode, &appt.ProcedureCode, &appt.Status,
		&appt.Notes, &appt.CreatedAt, &appt.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	c.JSON(http.StatusOK, appt)
}

func (s *Service) CancelAppointment(c *gin.Context) {
	apptID := c.Param("id")

	query := `
		UPDATE appointments 
		SET status = 'cancelled', updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := s.db.QueryRow(context.Background(), query, apptID).Scan(&updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         apptID,
		"status":     "cancelled",
		"updated_at": updatedAt,
	})
}

func (s *Service) ListPatientAppointments(c *gin.Context) {
	patientID := c.Param("patient_id")

	query := `
		SELECT id, patient_id, provider_id, appointment_time, diagnosis_code,
		       procedure_code, status, created_at
		FROM appointments
		WHERE patient_id = $1
		ORDER BY appointment_time DESC
	`

	rows, err := s.db.Query(context.Background(), query, patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list appointments"})
		return
	}
	defer rows.Close()

	var appointments []Appointment
	for rows.Next() {
		var appt Appointment
		err := rows.Scan(
			&appt.ID, &appt.PatientID, &appt.ProviderID, &appt.AppointmentTime,
			&appt.DiagnosisCode, &appt.ProcedureCode, &appt.Status, &appt.CreatedAt,
		)
		if err != nil {
			continue
		}
		appointments = append(appointments, appt)
	}

	c.JSON(http.StatusOK, appointments)
}
