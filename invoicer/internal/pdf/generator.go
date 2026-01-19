package pdf

import (
	"fmt"
	"invoice-generator/invoicer/internal/models"

	"github.com/jung-kurt/gofpdf"
)

// Generator handles PDF generation for invoices
type Generator struct {
	pdf *gofpdf.Fpdf
}

// NewGenerator creates a new PDF generator
func NewGenerator() *Generator {
	pdf := gofpdf.New("P", "mm", "A4", "")
	return &Generator{pdf: pdf}
}

// GenerateInvoice creates a PDF from invoice data
func (g *Generator) GenerateInvoice(invoice *models.Invoice) ([]byte, error) {
	g.pdf.AddPage()
	g.pdf.SetFont("Arial", "", 12)

	// Select template based on selectedTemplate field
	template := invoice.SelectedTemplate
	if template == "" {
		template = "minimal" // Default to minimal
	}

	// Draw invoice based on selected template
	switch template {
	case "corporate":
		g.drawCorporateInvoice(invoice)
	case "modern":
		g.drawModernInvoice(invoice)
	default:
		g.drawMinimalInvoice(invoice)
	}

	// Get PDF as bytes
	return g.GetPDFBytes()
}

func (g *Generator) drawMinimalInvoice(invoice *models.Invoice) {
	currencySymbol := getCurrencySymbol(invoice.Currency)

	// Header - INVOICE title and Business Name (side by side)
	g.pdf.SetFont("Arial", "B", 24)
	g.pdf.SetXY(15, 15)
	g.pdf.Cell(0, 10, "INVOICE")

	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetXY(120, 15)
	g.pdf.CellFormat(75, 10, invoice.BusinessName, "", 0, "R", false, 0, "")

	// Invoice number under title
	g.pdf.SetFont("Arial", "", 10)
	g.pdf.SetTextColor(100, 100, 100)
	g.pdf.SetXY(15, 25)
	g.pdf.Cell(0, 6, fmt.Sprintf("#%s", invoice.InvoiceNumber))

	// Business contact info (right aligned)
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetXY(120, 25)
	g.pdf.CellFormat(75, 5, invoice.BusinessEmail, "", 0, "R", false, 0, "")
	g.pdf.SetXY(120, 30)
	g.pdf.CellFormat(75, 5, invoice.BusinessPhone, "", 0, "R", false, 0, "")

	// Business address (right aligned, multi-line)
	if invoice.BusinessAddress != "" {
		g.pdf.SetXY(120, 35)
		g.pdf.SetFont("Arial", "", 9)
		g.pdf.MultiCell(75, 4, invoice.BusinessAddress, "", "R", false)
	}

	g.pdf.SetTextColor(0, 0, 0)

	// Bill To & Dates section (two columns)
	y := 55.0

	// Left column - Bill To
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(120, 120, 120)
	g.pdf.SetXY(15, y)
	g.pdf.Cell(0, 5, "BILL TO")
	g.pdf.SetTextColor(0, 0, 0)

	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetXY(15, y+6)
	g.pdf.Cell(0, 5, invoice.ClientName)

	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(100, 100, 100)
	g.pdf.SetXY(15, y+11)
	g.pdf.Cell(0, 4, invoice.ClientEmail)

	if invoice.ClientAddress != "" {
		g.pdf.SetXY(15, y+15)
		g.pdf.MultiCell(80, 4, invoice.ClientAddress, "", "L", false)
	}

	// Right column - Invoice & Due Date
	g.pdf.SetTextColor(120, 120, 120)
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetXY(120, y)
	g.pdf.Cell(40, 5, "Invoice Date:")
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.CellFormat(35, 5, invoice.InvoiceDate, "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(120, 120, 120)
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetXY(120, y+6)
	g.pdf.Cell(40, 5, "Due Date:")
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.CellFormat(35, 5, invoice.DueDate, "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)

	// Items table
	tableY := 90.0
	g.pdf.SetY(tableY)

	// Table header with thick bottom border
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetX(15)

	headers := []struct {
		text  string
		width float64
		align string
	}{
		{"DESCRIPTION", 60, "L"},
		{"QTY", 18, "R"},
		{"RATE", 25, "R"},
		{"TAX%", 18, "R"},
		{"DISC%", 18, "R"},
		{"AMOUNT", 28, "R"},
	}

	for _, h := range headers {
		g.pdf.CellFormat(h.width, 6, h.text, "", 0, h.align, false, 0, "")
	}
	g.pdf.Ln(6)

	// Thick line under header
	g.pdf.SetLineWidth(0.5)
	g.pdf.Line(15, g.pdf.GetY(), 195, g.pdf.GetY())
	g.pdf.Ln(2)

	// Table rows
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetLineWidth(0.1)

	for _, item := range invoice.Items {
		g.pdf.SetX(15)

		// Description
		g.pdf.CellFormat(60, 6, truncateString(item.Description, 40), "", 0, "L", false, 0, "")

		// Quantity
		g.pdf.CellFormat(18, 6, fmt.Sprintf("%.0f", item.Quantity), "", 0, "R", false, 0, "")

		// Rate
		g.pdf.CellFormat(25, 6, fmt.Sprintf("%s%.2f", currencySymbol, item.Rate), "", 0, "R", false, 0, "")

		// Tax%
		g.pdf.CellFormat(18, 6, fmt.Sprintf("%.0f%%", item.TaxRate), "", 0, "R", false, 0, "")

		// Disc%
		g.pdf.CellFormat(18, 6, fmt.Sprintf("%.0f%%", item.DiscountRate), "", 0, "R", false, 0, "")

		// Amount
		g.pdf.CellFormat(28, 6, fmt.Sprintf("%s%.2f", currencySymbol, item.Amount), "", 0, "R", false, 0, "")

		g.pdf.Ln(6)

		// Thin line under each row
		g.pdf.SetDrawColor(220, 220, 220)
		g.pdf.Line(15, g.pdf.GetY(), 195, g.pdf.GetY())
		g.pdf.Ln(1)
	}

	g.pdf.SetDrawColor(0, 0, 0)

	// Totals section (right aligned)
	totalsY := g.pdf.GetY() + 8
	totalsX := 125.0

	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(100, 100, 100)

	// Subtotal
	g.pdf.SetXY(totalsX, totalsY)
	g.pdf.Cell(35, 5, "Subtotal:")
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(35, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Subtotal), "", 0, "R", false, 0, "")
	totalsY += 5

	// Discount (if applicable)
	if invoice.DiscountRate > 0 {
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetXY(totalsX, totalsY)
		g.pdf.Cell(35, 5, fmt.Sprintf("Discount (%.0f%%):", invoice.DiscountRate))
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(35, 5, fmt.Sprintf("-%s%.2f", currencySymbol, invoice.DiscountAmount), "", 0, "R", false, 0, "")
		totalsY += 5
	}

	// Tax (if applicable)
	if invoice.TaxRate > 0 {
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetXY(totalsX, totalsY)
		g.pdf.Cell(35, 5, fmt.Sprintf("Tax (%.0f%%):", invoice.TaxRate))
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(35, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.TaxAmount), "", 0, "R", false, 0, "")
		totalsY += 5
	}

	// Total with thick top border
	totalsY += 3
	g.pdf.SetLineWidth(0.5)
	g.pdf.Line(totalsX, totalsY, 195, totalsY)
	totalsY += 2

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.SetXY(totalsX, totalsY)
	g.pdf.Cell(35, 6, "Total:")
	g.pdf.CellFormat(35, 6, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Total), "", 0, "R", false, 0, "")

	g.pdf.SetLineWidth(0.1)

	// Notes section (if present)
	if invoice.Notes != "" {
		notesY := totalsY + 20
		g.pdf.SetXY(15, notesY)
		g.pdf.SetFont("Arial", "B", 9)
		g.pdf.SetTextColor(120, 120, 120)
		g.pdf.Cell(0, 5, "NOTES")
		g.pdf.Ln(5)

		g.pdf.SetFont("Arial", "", 9)
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetX(15)
		g.pdf.MultiCell(180, 4, invoice.Notes, "", "L", false)
	}

	g.pdf.SetTextColor(0, 0, 0)
}

func (g *Generator) drawCorporateInvoice(invoice *models.Invoice) {
	currencySymbol := getCurrencySymbol(invoice.Currency)

	// Blue header background (RGB: 30, 58, 138 = blue-900)
	g.pdf.SetFillColor(30, 58, 138)
	g.pdf.Rect(0, 0, 210, 35, "F")

	// Header content - INVOICE title and Business Name
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.SetFont("Arial", "B", 20)
	g.pdf.SetXY(15, 10)
	g.pdf.Cell(0, 8, "INVOICE")

	// Invoice number
	g.pdf.SetFont("Arial", "", 10)
	g.pdf.SetTextColor(191, 219, 254) // blue-200
	g.pdf.SetXY(15, 20)
	g.pdf.Cell(0, 5, fmt.Sprintf("#%s", invoice.InvoiceNumber))

	// Business name (right aligned)
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.SetXY(110, 10)
	g.pdf.CellFormat(85, 8, invoice.BusinessName, "", 0, "R", false, 0, "")

	// Business contact (right aligned)
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(191, 219, 254)
	g.pdf.SetXY(110, 20)
	g.pdf.CellFormat(85, 4, invoice.BusinessEmail, "", 0, "R", false, 0, "")
	g.pdf.SetXY(110, 24)
	g.pdf.CellFormat(85, 4, invoice.BusinessPhone, "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)

	// Bill To & Invoice Info section
	y := 45.0

	// Left column - Bill To (in gray box)
	g.pdf.SetFillColor(249, 250, 251) // gray-50
	g.pdf.Rect(15, y, 90, 35, "F")

	g.pdf.SetFont("Arial", "B", 8)
	g.pdf.SetTextColor(120, 120, 120)
	g.pdf.SetXY(18, y+3)
	g.pdf.Cell(0, 4, "BILL TO")

	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.SetXY(18, y+10)
	g.pdf.Cell(0, 5, invoice.ClientName)

	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(80, 80, 80)
	g.pdf.SetXY(18, y+16)
	g.pdf.Cell(0, 4, invoice.ClientEmail)

	if invoice.ClientAddress != "" {
		g.pdf.SetXY(18, y+21)
		g.pdf.MultiCell(80, 4, invoice.ClientAddress, "", "L", false)
	}

	// Right column - Invoice details (in gray boxes)
	g.pdf.SetFillColor(249, 250, 251)
	g.pdf.Rect(110, y, 85, 10, "F")
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(80, 80, 80)
	g.pdf.SetXY(113, y+3)
	g.pdf.Cell(40, 5, "Invoice Date")
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(39, 5, invoice.InvoiceDate, "", 0, "R", false, 0, "")

	g.pdf.SetFillColor(249, 250, 251)
	g.pdf.Rect(110, y+12, 85, 10, "F")
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(80, 80, 80)
	g.pdf.SetXY(113, y+15)
	g.pdf.Cell(40, 5, "Due Date")
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(39, 5, invoice.DueDate, "", 0, "R", false, 0, "")

	// Amount Due box (highlighted in blue)
	g.pdf.SetFillColor(219, 234, 254) // blue-50
	g.pdf.SetDrawColor(147, 197, 253) // blue-200
	g.pdf.Rect(110, y+24, 85, 11, "FD")
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(30, 58, 138) // blue-900
	g.pdf.SetXY(113, y+27)
	g.pdf.Cell(40, 5, "Amount Due")
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.CellFormat(39, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Total), "", 0, "R", false, 0, "")

	g.pdf.SetDrawColor(0, 0, 0)
	g.pdf.SetTextColor(0, 0, 0)

	// Items table
	tableY := 95.0
	g.pdf.SetY(tableY)

	// Table header with gray background
	g.pdf.SetFillColor(229, 231, 235) // gray-200
	g.pdf.SetFont("Arial", "B", 8)
	g.pdf.SetX(15)

	headers := []struct {
		text  string
		width float64
		align string
	}{
		{"DESCRIPTION", 55, "L"},
		{"QTY", 18, "C"},
		{"RATE", 25, "R"},
		{"TAX%", 18, "R"},
		{"DISC%", 18, "R"},
		{"AMOUNT", 28, "R"},
	}

	for _, h := range headers {
		g.pdf.CellFormat(h.width, 7, h.text, "1", 0, h.align, true, 0, "")
	}
	g.pdf.Ln(7)

	// Table rows
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetFillColor(255, 255, 255)

	for _, item := range invoice.Items {
		g.pdf.SetX(15)
		g.pdf.CellFormat(55, 7, truncateString(item.Description, 35), "1", 0, "L", false, 0, "")
		g.pdf.CellFormat(18, 7, fmt.Sprintf("%.0f", item.Quantity), "1", 0, "C", false, 0, "")
		g.pdf.CellFormat(25, 7, fmt.Sprintf("%s%.2f", currencySymbol, item.Rate), "1", 0, "R", false, 0, "")
		g.pdf.CellFormat(18, 7, fmt.Sprintf("%.0f%%", item.TaxRate), "1", 0, "R", false, 0, "")
		g.pdf.CellFormat(18, 7, fmt.Sprintf("%.0f%%", item.DiscountRate), "1", 0, "R", false, 0, "")
		g.pdf.SetFont("Arial", "B", 9)
		g.pdf.CellFormat(28, 7, fmt.Sprintf("%s%.2f", currencySymbol, item.Amount), "1", 0, "R", false, 0, "")
		g.pdf.SetFont("Arial", "", 9)
		g.pdf.Ln(7)
	}

	// Totals section
	totalsY := g.pdf.GetY() + 8
	totalsX := 125.0

	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(100, 100, 100)

	// Subtotal
	g.pdf.SetXY(totalsX, totalsY)
	g.pdf.Cell(35, 5, "Subtotal")
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(35, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Subtotal), "", 0, "R", false, 0, "")
	totalsY += 5

	// Discount
	if invoice.DiscountRate > 0 {
		g.pdf.SetFillColor(249, 250, 251)
		g.pdf.Rect(totalsX, totalsY, 70, 5, "F")
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetXY(totalsX, totalsY)
		g.pdf.Cell(35, 5, fmt.Sprintf("Discount (%.0f%%)", invoice.DiscountRate))
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(35, 5, fmt.Sprintf("-%s%.2f", currencySymbol, invoice.DiscountAmount), "", 0, "R", false, 0, "")
		totalsY += 5
	}

	// Tax
	if invoice.TaxRate > 0 {
		g.pdf.SetFillColor(249, 250, 251)
		g.pdf.Rect(totalsX, totalsY, 70, 5, "F")
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetXY(totalsX, totalsY)
		g.pdf.Cell(35, 5, fmt.Sprintf("Tax (%.0f%%)", invoice.TaxRate))
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(35, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.TaxAmount), "", 0, "R", false, 0, "")
		totalsY += 5
	}

	// Total (blue background)
	totalsY += 2
	g.pdf.SetFillColor(30, 58, 138) // blue-900
	g.pdf.Rect(totalsX, totalsY, 70, 9, "F")
	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.SetXY(totalsX+2, totalsY+2)
	g.pdf.Cell(33, 5, "Total Due")
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.CellFormat(33, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Total), "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)

	// Notes
	if invoice.Notes != "" {
		notesY := totalsY + 18
		g.pdf.SetFillColor(249, 250, 251)
		g.pdf.Rect(15, notesY, 180, 4, "F")

		g.pdf.SetFont("Arial", "B", 8)
		g.pdf.SetTextColor(120, 120, 120)
		g.pdf.SetXY(15, notesY+5)
		g.pdf.Cell(0, 4, "PAYMENT NOTES")
		g.pdf.Ln(5)

		g.pdf.SetFont("Arial", "", 9)
		g.pdf.SetTextColor(80, 80, 80)
		g.pdf.SetX(15)
		g.pdf.MultiCell(180, 4, invoice.Notes, "", "L", false)
	}

	g.pdf.SetTextColor(0, 0, 0)
}

func (g *Generator) drawModernInvoice(invoice *models.Invoice) {
	currencySymbol := getCurrencySymbol(invoice.Currency)

	// Set purple gradient background (solid purple for PDF)
	g.pdf.SetFillColor(243, 232, 255) // purple-100
	g.pdf.Rect(0, 0, 210, 297, "F")

	// Header - INVOICE in gradient box
	g.pdf.SetFillColor(147, 51, 234) // purple-600
	g.pdf.RoundedRect(15, 12, 60, 12, 3, "1234", "F")
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.SetFont("Arial", "B", 20)
	g.pdf.SetXY(15, 14.5)
	g.pdf.CellFormat(60, 8, "INVOICE", "", 0, "C", false, 0, "")

	// Invoice number
	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetTextColor(80, 80, 80)
	g.pdf.SetXY(15, 27)
	g.pdf.Cell(0, 5, fmt.Sprintf("#%s", invoice.InvoiceNumber))

	// Business name (right aligned with gradient color effect)
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetTextColor(147, 51, 234) // purple-600
	g.pdf.SetXY(110, 15)
	g.pdf.CellFormat(85, 8, invoice.BusinessName, "", 0, "R", false, 0, "")

	// Business contact (right aligned)
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(100, 100, 100)
	g.pdf.SetXY(110, 25)
	g.pdf.CellFormat(85, 4, invoice.BusinessEmail, "", 0, "R", false, 0, "")
	g.pdf.SetXY(110, 29)
	g.pdf.CellFormat(85, 4, invoice.BusinessPhone, "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)

	// Bill To & Dates cards (white rounded boxes)
	y := 45.0

	// Bill To card
	g.pdf.SetFillColor(255, 255, 255)
	g.pdf.RoundedRect(15, y, 85, 30, 3, "23", "F")

	// Purple accent bar
	g.pdf.SetFillColor(147, 51, 234)
	g.pdf.Rect(15, y, 2, 30, "F")

	g.pdf.SetFont("Arial", "B", 8)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.SetXY(20, y+3)
	g.pdf.Cell(0, 4, "BILL TO")

	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetXY(20, y+10)
	g.pdf.Cell(0, 5, invoice.ClientName)

	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(100, 100, 100)
	g.pdf.SetXY(20, y+16)
	g.pdf.Cell(0, 4, invoice.ClientEmail)

	if invoice.ClientAddress != "" {
		g.pdf.SetXY(20, y+21)
		g.pdf.MultiCell(75, 4, invoice.ClientAddress, "", "L", false)
	}

	// Invoice details card
	g.pdf.SetFillColor(255, 255, 255)
	g.pdf.RoundedRect(110, y, 85, 30, 3, "1234", "F")

	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(100, 100, 100)
	g.pdf.SetXY(113, y+5)
	g.pdf.Cell(40, 4, "Invoice Date")
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(39, 4, invoice.InvoiceDate, "", 0, "R", false, 0, "")

	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(100, 100, 100)
	g.pdf.SetXY(113, y+12)
	g.pdf.Cell(40, 4, "Due Date")
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(39, 4, invoice.DueDate, "", 0, "R", false, 0, "")

	// Divider
	g.pdf.SetDrawColor(229, 231, 235)
	g.pdf.Line(113, y+19, 192, y+19)
	g.pdf.SetDrawColor(0, 0, 0)

	// Amount Due
	g.pdf.SetFont("Arial", "B", 9)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.SetXY(113, y+22)
	g.pdf.Cell(40, 5, "Amount Due")
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetTextColor(147, 51, 234)
	g.pdf.CellFormat(39, 5, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Total), "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)

	// Items table (white rounded box)
	tableY := 85.0
	g.pdf.SetFillColor(255, 255, 255)
	tableHeight := 15.0 + float64(len(invoice.Items))*7.0
	g.pdf.RoundedRect(15, tableY, 180, tableHeight, 3, "34", "F")

	// Table header with gradient
	g.pdf.SetFillColor(147, 51, 234) // purple-600
	g.pdf.Rect(15, tableY, 180, 8, "F")

	g.pdf.SetFont("Arial", "B", 8)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.SetXY(20, tableY+2)

	headers := []struct {
		text  string
		width float64
		align string
	}{
		{"DESCRIPTION", 55, "L"},
		{"QTY", 18, "C"},
		{"RATE", 25, "R"},
		{"TAX%", 18, "R"},
		{"DISC%", 18, "R"},
		{"AMOUNT", 26, "R"},
	}

	x := 20.0
	for _, h := range headers {
		g.pdf.SetXY(x, tableY+2)
		g.pdf.CellFormat(h.width, 4, h.text, "", 0, h.align, false, 0, "")
		x += h.width
	}

	// Table rows
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)
	rowY := tableY + 10

	for _, item := range invoice.Items {
		g.pdf.SetXY(20, rowY)
		g.pdf.CellFormat(55, 5, truncateString(item.Description, 35), "", 0, "L", false, 0, "")
		g.pdf.CellFormat(18, 5, fmt.Sprintf("%.0f", item.Quantity), "", 0, "C", false, 0, "")
		g.pdf.CellFormat(25, 5, fmt.Sprintf("%s%.2f", currencySymbol, item.Rate), "", 0, "R", false, 0, "")
		g.pdf.CellFormat(18, 5, fmt.Sprintf("%.0f%%", item.TaxRate), "", 0, "R", false, 0, "")
		g.pdf.CellFormat(18, 5, fmt.Sprintf("%.0f%%", item.DiscountRate), "", 0, "R", false, 0, "")
		g.pdf.SetFont("Arial", "B", 9)
		g.pdf.CellFormat(26, 5, fmt.Sprintf("%s%.2f", currencySymbol, item.Amount), "", 0, "R", false, 0, "")
		g.pdf.SetFont("Arial", "", 9)
		rowY += 7

		// Light separator
		g.pdf.SetDrawColor(243, 244, 246)
		g.pdf.Line(20, rowY-1, 190, rowY-1)
		g.pdf.SetDrawColor(0, 0, 0)
	}

	// Totals card (white rounded box)
	totalsY := tableY + tableHeight + 8
	totalsHeight := 25.0
	if invoice.DiscountRate > 0 {
		totalsHeight += 5
	}
	if invoice.TaxRate > 0 {
		totalsHeight += 5
	}

	g.pdf.SetFillColor(255, 255, 255)
	g.pdf.RoundedRect(110, totalsY, 85, totalsHeight, 3, "1234", "F")

	ty := totalsY + 5
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(100, 100, 100)

	// Subtotal
	g.pdf.SetXY(113, ty)
	g.pdf.Cell(40, 4, "Subtotal")
	g.pdf.SetFont("Arial", "", 9)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.CellFormat(39, 4, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Subtotal), "", 0, "R", false, 0, "")
	ty += 5

	// Discount
	if invoice.DiscountRate > 0 {
		g.pdf.SetFont("Arial", "", 9)
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetXY(113, ty)
		g.pdf.Cell(40, 4, fmt.Sprintf("Discount (%.0f%%)", invoice.DiscountRate))
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(39, 4, fmt.Sprintf("-%s%.2f", currencySymbol, invoice.DiscountAmount), "", 0, "R", false, 0, "")
		ty += 5
	}

	// Tax
	if invoice.TaxRate > 0 {
		g.pdf.SetFont("Arial", "", 9)
		g.pdf.SetTextColor(100, 100, 100)
		g.pdf.SetXY(113, ty)
		g.pdf.Cell(40, 4, fmt.Sprintf("Tax (%.0f%%)", invoice.TaxRate))
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.CellFormat(39, 4, fmt.Sprintf("%s%.2f", currencySymbol, invoice.TaxAmount), "", 0, "R", false, 0, "")
		ty += 5
	}

	// Total with gradient background
	ty += 2
	g.pdf.SetFillColor(147, 51, 234) // purple-600
	g.pdf.Rect(110, ty, 85, 8, "F")
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.SetTextColor(255, 255, 255)
	g.pdf.SetXY(113, ty+2)
	g.pdf.Cell(40, 4, "Total")
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.CellFormat(39, 4, fmt.Sprintf("%s%.2f", currencySymbol, invoice.Total), "", 0, "R", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)

	// Notes
	if invoice.Notes != "" {
		notesY := totalsY + totalsHeight + 10
		g.pdf.SetFillColor(255, 255, 255)
		g.pdf.RoundedRect(15, notesY, 180, 30, 3, "23", "F")

		// Purple accent bar
		g.pdf.SetFillColor(147, 51, 234)
		g.pdf.Rect(15, notesY, 2, 30, "F")

		g.pdf.SetFont("Arial", "B", 8)
		g.pdf.SetTextColor(0, 0, 0)
		g.pdf.SetXY(20, notesY+3)
		g.pdf.Cell(0, 4, "PAYMENT NOTES")
		g.pdf.Ln(5)

		g.pdf.SetFont("Arial", "", 9)
		g.pdf.SetTextColor(80, 80, 80)
		g.pdf.SetX(20)
		g.pdf.MultiCell(170, 4, invoice.Notes, "", "L", false)
	}

	g.pdf.SetTextColor(0, 0, 0)
}

// Helper functions

func getCurrencySymbol(currency string) string {
	symbols := map[string]string{
		"USD": "$",
		"EUR": "€",
		"GBP": "£",
		"JPY": "¥",
		"AUD": "A$",
		"CAD": "C$",
		"INR": "₹",
	}

	if symbol, ok := symbols[currency]; ok {
		return symbol
	}
	return currency + " "
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// GetPDFBytes returns the PDF as a byte slice
func (g *Generator) GetPDFBytes() ([]byte, error) {
	var buf []byte
	writer := &bytesWriter{data: &buf}
	err := g.pdf.Output(writer)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// bytesWriter implements io.Writer for in-memory PDF generation
type bytesWriter struct {
	data *[]byte
}

func (w *bytesWriter) Write(p []byte) (n int, err error) {
	*w.data = append(*w.data, p...)
	return len(p), nil
}
