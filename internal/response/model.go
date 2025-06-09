package response

type Model struct {
	ID        string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedAt string `json:"created_at" example:"2023-10-01T12:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2023-10-01T12:00:00Z"`
	DeletedAt string `json:"deleted_at,omitempty" example:"2023-10-01T12:00:00Z"`
}
