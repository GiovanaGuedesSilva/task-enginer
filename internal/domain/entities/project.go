package entities

import (
	"errors"
	"task-engine/internal/domain/entities/common"
	"time"
)

type Project struct {
	ID          int64                    `json:"id" db:"id"`
	Name        string                   `json:"name" db:"name"`
	Description string                   `json:"description" db:"description"`
	Status      common.ProjectStatus     `json:"status" db:"status"`
	Priority    common.ProjectPriority   `json:"priority" db:"priority"`
	OwnerID     int64                    `json:"owner_id" db:"owner_id"`
	TeamID      int64                    `json:"team_id,omitempty" db:"team_id"`
	StartDate   time.Time                `json:"start_date,omitempty" db:"start_date"`
	EndDate     time.Time                `json:"end_date,omitempty" db:"end_date"`
	Budget      float64                  `json:"budget,omitempty" db:"budget"`
	CreatedAt   time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at" db:"updated_at"`
	DeletedAt   time.Time                `json:"deleted_at,omitempty" db:"deleted_at"`
}

func NewProject(name string, description string, ownerID int64) (*Project, error) {
	project := &Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Status:      common.ProjectStatusActive,
		Priority:    common.ProjectPriorityMedium,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	return project, nil
}

func NewProjectWithDetails(name, description string, ownerID int64, status common.ProjectStatus, priority common.ProjectPriority, teamID int64, startDate, endDate time.Time, budget float64) (*Project, error) {
	project := &Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		Status:      status,
		Priority:    priority,
		TeamID:      teamID,
		StartDate:   startDate,
		EndDate:     endDate,
		Budget:      budget,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := project.Validate(); err != nil {
		return nil, err
	}

	return project, nil
}

// Validations methods

func (p *Project) Validate() error {
	return common.ValidateFields(
		// basic validations
		common.ValidateRequired("name", p.Name),
		common.ValidateStringLength("name", p.Name, 255),
		common.ValidateStringLength("description", p.Description, 1000),
		common.ValidatePositiveInt("owner_id", p.OwnerID),

		// enum validations
		common.ValidateProjectStatus("status", p.Status),
		common.ValidateProjectPriority("priority", p.Priority),

		// date validations
		common.ValidateDateOrder("start_date", "end_date", p.StartDate, p.EndDate),

		// budget validation
		common.ValidatePositiveFloat("budget", p.Budget),

		// domain specific validations
		p.validateProjectSpecificRules(),
	)
}

func (p *Project) validateProjectSpecificRules() *common.FieldValidationError {
	if p.Status == common.ProjectStatusActive && p.Budget < 0 {
		return &common.FieldValidationError{
			Field:   "budget",
			Message: "active projects cannot have negative budget",
		}
	}

	if p.Status == common.ProjectStatusArchived && p.EndDate.IsZero() {
		return &common.FieldValidationError{
			Field:   "end_date",
			Message: "archived projects must have an end date",
		}
	}

	if p.Priority == common.ProjectPriorityUrgent && p.EndDate.IsZero() {
		return &common.FieldValidationError{
			Field:   "end_date",
			Message: "urgent priority projects must have an end date",
		}
	}

	return nil
}

// Business methods

func (p *Project) IsActive() bool {
	return p.Status == common.ProjectStatusActive && p.DeletedAt.IsZero()
}

func (p *Project) IsArchived() bool {
	return p.Status == common.ProjectStatusArchived
}

func (p *Project) IsDeleted() bool {
	return p.Status == common.ProjectStatusDeleted || !p.DeletedAt.IsZero()
}

func (p *Project) CanBeDeleted() bool {
	return p.Status != common.ProjectStatusDeleted
}

func (p *Project) CanBeArchived() bool {
	return p.Status != common.ProjectStatusArchived
}

func (p *Project) IsOverdue() bool {
	return p.EndDate.Before(time.Now()) && p.IsActive()
}

func (p *Project) GetProgress() float64 {
	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return 0.0
	}

	now := time.Now()
	totalDuration := p.EndDate.Sub(p.StartDate)
	elapsedDuration := now.Sub(p.StartDate)

	if elapsedDuration <= 0 {
		return 0.0
	}

	if elapsedDuration >= totalDuration {
		return 100.0
	}

	progress := (elapsedDuration.Seconds() / totalDuration.Seconds()) * 100
	return progress
}

func (p *Project) GetDaysRemaining() int {
	if p.EndDate.IsZero() {
		return 0
	}

	daysRemaining := int(p.EndDate.Sub(time.Now()).Hours() / 24)
	if daysRemaining < 0 {
		return 0
	}
	return daysRemaining
}

// Modification methods

func (p *Project) UpdateName(name string) error {
	if err := common.ValidateFields(
		common.ValidateRequired("name", name),
		common.ValidateStringLength("name", name, 255),
	); err != nil {
		return err
	}

	p.Name = name
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) UpdateDescription(description string) error {
	if err := common.ValidateStringLength("description", description, 1000); err != nil {
		return err
	}

	p.Description = description
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) UpdateStatus(status common.ProjectStatus) error {
	if err := common.ValidateProjectStatus("status", status); err != nil {
		return err
	}

	p.Status = status
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) UpdatePriority(priority common.ProjectPriority) error {
	if err := common.ValidateProjectPriority("priority", priority); err != nil {
		return err
	}

	p.Priority = priority
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) UpdateDates(startDate, endDate time.Time) error {
	if err := common.ValidateFields(
		common.ValidateDateOrder("start_date", "end_date", startDate, endDate),
	); err != nil {
		return err
	}

	p.StartDate = startDate
	p.EndDate = endDate
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) UpdateBudget(budget float64) error {
	if err := common.ValidateFields(
		common.ValidatePositiveFloat("budget", budget),
	); err != nil {
		return err
	}

	if p.Status == common.ProjectStatusActive && budget < 0 {
		return &common.FieldValidationError{
			Field:   "budget",
			Message: "active projects cannot have negative budget",
		}
	}

	p.Budget = budget
	p.UpdatedAt = time.Now()
	return nil
}

// Business actions methods

func (p *Project) Archive() error {
	if !p.CanBeArchived() {
		return errors.New("project cannot be archived")
	}

	p.Status = common.ProjectStatusArchived
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) Delete() error {
	if !p.CanBeDeleted() {
		return errors.New("project cannot be deleted")
	}

	now := time.Now()
	p.Status = common.ProjectStatusDeleted
	p.DeletedAt = now
	p.UpdatedAt = now
	return nil
}

func (p *Project) Restore() error {
	if !p.IsDeleted() {
		return errors.New("project is not deleted")
	}

	p.Status = common.ProjectStatusActive
	p.DeletedAt = time.Time{}
	p.UpdatedAt = time.Now()
	return nil
}

// Builder pattern

type ProjectBuilder struct {
	project *Project
}

func NewProjectBuilder() *ProjectBuilder {
	return &ProjectBuilder{
		project: &Project{
			Status:    common.ProjectStatusActive,
			Priority:  common.ProjectPriorityMedium,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *ProjectBuilder) WithName(name string) *ProjectBuilder {
	b.project.Name = name
	return b
}

func (b *ProjectBuilder) WithDescription(description string) *ProjectBuilder {
	b.project.Description = description
	return b
}

func (b *ProjectBuilder) WithOwner(ownerID int64) *ProjectBuilder {
	b.project.OwnerID = ownerID
	return b
}

func (b *ProjectBuilder) WithStatus(status common.ProjectStatus) *ProjectBuilder {
	b.project.Status = status
	return b
}

func (b *ProjectBuilder) WithPriority(priority common.ProjectPriority) *ProjectBuilder {
	b.project.Priority = priority
	return b
}

func (b *ProjectBuilder) WithTeam(teamID int64) *ProjectBuilder {
	b.project.TeamID = teamID
	return b
}

func (b *ProjectBuilder) WithStartDate(startDate time.Time) *ProjectBuilder {
	b.project.StartDate = startDate
	return b
}

func (b *ProjectBuilder) WithEndDate(endDate time.Time) *ProjectBuilder {
	b.project.EndDate = endDate
	return b
}

func (b *ProjectBuilder) WithBudget(budget float64) *ProjectBuilder {
	b.project.Budget = budget
	return b
}

func (b *ProjectBuilder) Build() (*Project, error) {
	if err := b.project.Validate(); err != nil {
		return nil, err
	}
	return b.project, nil
}

func (b *ProjectBuilder) BuildUnsafe() *Project {
	return b.project
}
