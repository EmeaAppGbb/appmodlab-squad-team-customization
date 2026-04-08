package provider

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

type Provider struct {
	ID          string    `json:"id" db:"id"`
	NPI         string    `json:"npi" db:"npi"` // National Provider Identifier
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	Specialty   string    `json:"specialty" db:"specialty"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	LicenseNum  string    `json:"license_number" db:"license_number"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) CreateProvider(c *gin.Context) {
	var provider Provider
	if err := c.ShouldBindJSON(&provider); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO providers (npi, first_name, last_name, specialty, email, phone, license_number, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(context.Background(), query,
		provider.NPI, provider.FirstName, provider.LastName, provider.Specialty,
		provider.Email, provider.Phone, provider.LicenseNum, true,
	).Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create provider"})
		return
	}

	provider.IsAvailable = true
	c.JSON(http.StatusCreated, provider)
}

func (s *Service) GetProvider(c *gin.Context) {
	providerID := c.Param("id")

	var provider Provider
	query := `
		SELECT id, npi, first_name, last_name, specialty, email, phone, 
		       license_number, is_available, created_at, updated_at
		FROM providers WHERE id = $1
	`

	err := s.db.QueryRow(context.Background(), query, providerID).Scan(
		&provider.ID, &provider.NPI, &provider.FirstName, &provider.LastName,
		&provider.Specialty, &provider.Email, &provider.Phone,
		&provider.LicenseNum, &provider.IsAvailable,
		&provider.CreatedAt, &provider.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func (s *Service) ListProviders(c *gin.Context) {
	query := `
		SELECT id, npi, first_name, last_name, specialty, email, phone, is_available, created_at
		FROM providers
		WHERE is_available = true
		ORDER BY last_name, first_name
	`

	rows, err := s.db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list providers"})
		return
	}
	defer rows.Close()

	var providers []Provider
	for rows.Next() {
		var provider Provider
		err := rows.Scan(
			&provider.ID, &provider.NPI, &provider.FirstName, &provider.LastName,
			&provider.Specialty, &provider.Email, &provider.Phone,
			&provider.IsAvailable, &provider.CreatedAt,
		)
		if err != nil {
			continue
		}
		providers = append(providers, provider)
	}

	c.JSON(http.StatusOK, providers)
}
