package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Name      string     `json:"name"`
	Color     *string    `json:"color,omitempty"`
	Icon      *string    `json:"icon,omitempty"`
	IsSystem  bool       `json:"is_system"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type CategoryRequest struct {
	Name     string     `json:"name"`
	Color    *string    `json:"color"`
	Icon     *string    `json:"icon"`
	ParentID *uuid.UUID `json:"parent_id"`
}

// SystemCategory represents seed data for system-wide categories.
type SystemCategory struct {
	Name  string
	Color string
	Icon  string
}

var DefaultSystemCategories = []SystemCategory{
	{Name: "Alimentação", Color: "#FF6B6B", Icon: "food"},
	{Name: "Transporte", Color: "#4ECDC4", Icon: "transport"},
	{Name: "Saúde", Color: "#45B7D1", Icon: "health"},
	{Name: "Educação", Color: "#96CEB4", Icon: "education"},
	{Name: "Lazer", Color: "#FFEAA7", Icon: "entertainment"},
	{Name: "Moradia", Color: "#DDA0DD", Icon: "home"},
	{Name: "Vestuário", Color: "#98D8C8", Icon: "clothing"},
	{Name: "Serviços", Color: "#F7DC6F", Icon: "services"},
	{Name: "Investimentos", Color: "#82E0AA", Icon: "investment"},
	{Name: "Receita", Color: "#58D68D", Icon: "income"},
	{Name: "Outros", Color: "#BDC3C7", Icon: "other"},
}
