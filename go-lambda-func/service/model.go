package src

import "time"

type EmployeeRole string

type CreateBusinessRequest struct {
	Name       *string `json:"name"`
	LegalName  *string `json:"legal_name"`
	Email      *string `json:"email"`
	Address1   *string `json:"address1"`
	Address2   *string `json:"address2"`
	City       *string `json:"city"`
	State      *string `json:"state"`
	PostalCode *string `json:"postal_code"`
	Country    *string `json:"country"`
	Currency   *string `json:"currency"`
	Phone      *string `json:"phone"`
	Locale     *string `json:"locale"`
	Gst        *string `json:"gst"`
	Pan        *string `json:"pan"`
	CreatedBy  *string `json:"created_by"`
}

type CreateBusinessResponse struct {
	BusinessId *int       `json:"business_id"`
	Name       *string    `json:"name"`
	LegalName  *string    `json:"legal_name"`
	Email      *string    `json:"email"`
	Address1   *string    `json:"address1"`
	Address2   *string    `json:"address2"`
	City       *string    `json:"city"`
	State      *string    `json:"state"`
	PostalCode *string    `json:"postal_code"`
	Country    *string    `json:"country"`
	Currency   *string    `json:"currency"`
	Phone      *string    `json:"phone"`
	Locale     *string    `json:"locale"`
	Gst        *string    `json:"gst"`
	Pan        *string    `json:"pan"`
	CreatedBy  *string    `json:"created_by"`
	CreatedAt  *time.Time `json:"created_at"`
}

type BusinessDao struct {
	PK *string `json:"PK"`
	SK *string `json:"SK"`
	*Business
}

type Business struct {
	Type            *string
	BusinessId      *int                    `json:"business_id"`
	Name            *string                 `json:"name"`
	LegalName       *string                 `json:"legal_name"`
	Email           *string                 `json:"email"`
	Address1        *string                 `json:"address1"`
	Address2        *string                 `json:"address2"`
	City            *string                 `json:"city"`
	State           *string                 `json:"state"`
	PostalCode      *string                 `json:"postal_code"`
	Country         *string                 `json:"country"`
	Currency        *string                 `json:"currency"`
	Phone           *string                 `json:"phone"`
	Locale          *string                 `json:"locale"`
	CreatedBy       *string                 `json:"created_by"`
	UpdatedBy       *string                 `json:"updated_by"`
	CreatedAt       *time.Time              `json:"created_at"`
	UpdatedAt       *time.Time              `json:"updated_at"`
	CustomAttribute *map[string]interface{} `json:"custom_attribute"`
	Logo            *Image                  `json:"logo"`
}

type EmployeeDao struct {
	PK *string `json:"PK"`
	SK *string `json:"SK"`
	*Employee
}

type Employee struct {
	Type          *string
	EmployeeId    *string    `json:"employee_id"`
	FirstName     *string    `json:"first_name"`
	MiddleName    *string    `json:"middle_name"`
	LastName      *string    `json:"last_name"`
	Locale        *string    `json:"locale"`
	Email         *string    `json:"email"`
	Phone         *string    `json:"phone"`
	Gender        *string    `json:"gender"`
	Picture       *string    `json:"picture"`
	EmailVerified *bool      `json:"email_verified"`
	PhoneVerified *bool      `json:"phone_verified"`
	Dob           *time.Time `json:"dob"`
	JoinedAt      *time.Time `json:"joined_at"`
	CreatedBy     *string    `json:"created_by"`
	CreatedAt     *time.Time `json:"created_at"`
}

type XPOSApiError struct {
	error
	ErrorMessage string `json:"error_message"`
}

type Sequence struct {
	SequenceType  *string `json:"sequenceType"`
	SequenceValue *int    `json:"sequenceValue"`
}

type CreateStoreEmployeeRequest struct {
	BusinessId    *string `json:"business_id"`
	Phone         *string `json:"phone"`
	FirstName     *string `json:"first_name"`
	MiddleName    *string `json:"middle_name"`
	LastName      *string `json:"last_name"`
	Locale        *string `json:"locale"`
	CreatedUserId *string `json:"created_user_id"`
}

type EmployeeUpdateRequest struct {
	FirstName  *string    `json:"first_name"`
	MiddleName *string    `json:"middle_name"`
	LastName   *string    `json:"last_name"`
	Locale     *string    `json:"locale"`
	Email      *string    `json:"email"`
	Gender     *string    `json:"gender"`
	Picture    *string    `json:"picture"`
	Dob        *time.Time `json:"dob"`
}

type Image struct {
	Large  *string `json:"large"`
	Medium *string `json:"medium"`
	Small  *string `json:"small"`
}

type StoreEmployeeRoleDao struct {
	PK   *string `json:"PK"`
	SK   *string `json:"SK"`
	GPK1 *string `json:"GPK1"`
	GSK1 *string `json:"GSK1"`
	*StoreEmployeeRole
}

type StoreEmployeeRole struct {
	EmployeeId   *string    `json:"employee_id"`
	StoreId      *string    `json:"store_id"`
	Roles        *string    `json:"roles"`
	Locale       *string    `json:"locale"`
	StoreName    *string    `json:"store_name"`
	EmployeeName *string    `json:"employee_name"`
	JoinedAt     *time.Time `json:"joined_at"`
	CreatedBy    *string    `json:"created_by"`
	CreatedAt    *time.Time `json:"created_at"`
}

type StoreEmployeeResponse struct {
	Employee  *Employee          `json:"employee"`
	StoreData *StoreEmployeeRole `json:"store_data"`
}
