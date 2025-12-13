package dto

type PersonCreate struct {
	HondaId    string  `json:"honda_id" binding:"required,max=20"`
	Name       string  `json:"name" binding:"required,min=3,max=150"`
	JobTitle   *string `json:"job_title" binding:"omitempty,min=2,max=150"`
	Role       string  `json:"role" binding:"required,oneof=teamleader sales_portal admin staff viewer salesman"`
	DealerCode *string `json:"dealer_code" binding:"omitempty"`
	Active     *bool   `json:"active" binding:"omitempty"`
}

type PersonUpdate struct {
	HondaId    *string `json:"honda_id" binding:"omitempty,max=20"`
	Name       *string `json:"name" binding:"omitempty,min=3,max=150"`
	JobTitle   *string `json:"job_title" binding:"omitempty,min=2,max=150"`
	Role       *string `json:"role" binding:"omitempty,oneof=teamleader sales_portal admin staff viewer salesman"`
	DealerCode *string `json:"dealer_code" binding:"omitempty"`
	Active     *bool   `json:"active" binding:"omitempty"`
}
