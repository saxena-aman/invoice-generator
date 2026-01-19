package handlers

import (
	"encoding/json"
	"fmt"
	"invoice-generator/invoicer/internal/models"
	"invoice-generator/invoicer/internal/pdf"
	"net/http"
)

// InvoiceHandler handles invoice-related HTTP requests
type InvoiceHandler struct{}

// NewInvoiceHandler creates a new invoice handler
func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{}
}

// GeneratePDF handles POST /api/generate-pdf requests
func (h *InvoiceHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse JSON request body
	var invoice models.Invoice
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&invoice); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate invoice data
	if err := validateInvoice(&invoice); err != nil {
		http.Error(w, fmt.Sprintf("Invalid invoice data: %v", err), http.StatusBadRequest)
		return
	}

	// Generate PDF
	generator := pdf.NewGenerator()
	pdfData, err := generator.GenerateInvoice(&invoice)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating PDF: %v", err), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice-%s.pdf", invoice.InvoiceNumber))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfData)))

	// Write PDF data to response
	w.WriteHeader(http.StatusOK)
	w.Write(pdfData)
}

// HealthCheck handles GET /health requests
func (h *InvoiceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "invoice-generator-api",
	})
}

// validateInvoice performs basic validation on invoice data
func validateInvoice(invoice *models.Invoice) error {
	if invoice.InvoiceNumber == "" {
		return fmt.Errorf("invoice number is required")
	}
	if invoice.BusinessName == "" {
		return fmt.Errorf("business name is required")
	}
	if invoice.ClientName == "" {
		return fmt.Errorf("client name is required")
	}
	if len(invoice.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}
	if invoice.Total <= 0 {
		return fmt.Errorf("total must be greater than zero")
	}
	return nil
}
