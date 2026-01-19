package models

// LineItem represents a single line item in the invoice
type LineItem struct {
	Description  string  `json:"description"`
	Quantity     float64 `json:"quantity"`
	Rate         float64 `json:"rate"`
	TaxRate      float64 `json:"taxRate"`
	DiscountRate float64 `json:"discountRate"`
	Amount       float64 `json:"amount"`
}

// Invoice represents the complete invoice data
type Invoice struct {
	// Invoice details
	InvoiceNumber string `json:"invoiceNumber"`
	InvoiceDate   string `json:"invoiceDate"`
	DueDate       string `json:"dueDate"`

	// Business information
	BusinessName    string `json:"businessName"`
	BusinessEmail   string `json:"businessEmail"`
	BusinessPhone   string `json:"businessPhone"`
	BusinessAddress string `json:"businessAddress"`

	// Client information
	ClientName    string `json:"clientName"`
	ClientEmail   string `json:"clientEmail"`
	ClientAddress string `json:"clientAddress"`

	// Line items
	Items []LineItem `json:"items"`

	// Totals
	Subtotal       float64 `json:"subtotal"`
	DiscountRate   float64 `json:"discountRate"`
	DiscountAmount float64 `json:"discountAmount"`
	TaxRate        float64 `json:"taxRate"`
	TaxAmount      float64 `json:"taxAmount"`
	Total          float64 `json:"total"`

	// Additional
	Currency         string `json:"currency"`
	Notes            string `json:"notes"`
	SelectedTemplate string `json:"selectedTemplate"` // "minimal", "corporate", or "modern"
}
