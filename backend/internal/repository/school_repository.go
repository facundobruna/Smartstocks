package repository

import (
	"database/sql"

	"github.com/smartstocks/backend/internal/models"
)

type SchoolRepository struct {
	db *sql.DB
}

func NewSchoolRepository(db *sql.DB) *SchoolRepository {
	return &SchoolRepository{db: db}
}

// GetAllSchools obtiene todos los colegios activos
func (r *SchoolRepository) GetAllSchools() ([]models.School, error) {
	query := `
		SELECT id, name, location, is_active, created_at
		FROM schools WHERE is_active = TRUE
		ORDER BY name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schools := []models.School{}
	for rows.Next() {
		var school models.School
		err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Location,
			&school.IsActive,
			&school.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		schools = append(schools, school)
	}

	return schools, nil
}

// GetSchoolByID obtiene un colegio por ID
func (r *SchoolRepository) GetSchoolByID(schoolID string) (*models.School, error) {
	school := &models.School{}
	query := `
		SELECT id, name, location, is_active, created_at
		FROM schools WHERE id = ?
	`

	err := r.db.QueryRow(query, schoolID).Scan(
		&school.ID,
		&school.Name,
		&school.Location,
		&school.IsActive,
		&school.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return school, err
}
