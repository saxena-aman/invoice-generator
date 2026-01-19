# Invoice Generator - Go Backend

A Go-based REST API for generating professional PDF invoices.

## Features

- ✅ RESTful API for PDF generation
- ✅ Support for item-level tax and discount
- ✅ Support for bill-level tax and discount  
- ✅ Professional PDF layout
- ✅ CORS enabled for React frontend
- ✅ JSON request/response
- ✅ Input validation

## Project Structure

```
invoicer/
├── main.go                  # HTTP server and routing
├── internal/
│   ├── models/
│   │   └── invoice.go      # Data models
│   ├── handlers/
│   │   └── invoice.go      # HTTP handlers
│   └── pdf/
│       └── generator.go    # PDF generation logic
├── go.mod
└── go.sum
```

## API Endpoints

### Generate PDF
**POST** `/api/generate-pdf`

**Request Body:**
```json
{
  "selectedTemplate": "minimal",  // Options: "minimal", "corporate", "modern"
  "invoiceNumber": "INV-001",
  "invoiceDate": "2026-01-20",
  "dueDate": "2026-02-20",
  "businessName": "My Business",
  "businessEmail": "business@example.com",
  "businessPhone": "+1234567890",
  "businessAddress": "123 Business St",
  "clientName": "Client Name",
  "clientEmail": "client@example.com",
  "clientAddress": "456 Client Ave",
  "currency": "USD",
  "items": [
    {
      "description": "Service",
      "quantity": 10,
      "rate": 100,
      "taxRate": 10,
      "discountRate": 5,
      "amount": 1045
    }
  ],
  "subtotal": 1045,
  "discountRate": 5,
  "discountAmount": 52.25,
  "taxRate": 10,
  "taxAmount": 99.28,
  "total": 1092.03,
  "notes": "Payment terms: Net 30"
}
```

**Template Options:**
- `minimal` - Clean, simple design (default if not specified)
- `corporate` - Professional blue theme with header background
- `modern` - Stylish purple gradient design with cards

**Response:**
- Content-Type: `application/pdf`
- Returns PDF file as attachment

### Health Check
**GET** `/health`

**Response:**
```json
{
  "status": "healthy",
  "service": "invoice-generator-api"
}
```

## Running the Server

### Development
```bash
# From the invoicer directory
go run main.go
```

### Production Build
```bash
# Build binary
go build -o bin/invoicer

# Run binary
./bin/invoicer
```

The server will start on port `8080`.

## Dependencies

- `github.com/jung-kurt/gofpdf` - PDF generation
- `github.com/gorilla/mux` - HTTP routing
- `github.com/rs/cors` - CORS middleware

## Testing

Use curl to test the API:

```bash
curl -X POST http://localhost:8080/api/generate-pdf \
  -H "Content-Type: application/json" \
  -d @test-invoice.json \
  --output invoice.pdf
```

## Integration with React Frontend

The React app (running on `localhost:5173`) can call this API to generate PDFs server-side instead of client-side PDF generation.

Update the `handleDownloadPDF` function in `InvoiceForm.jsx` to:

```javascript
const handleDownloadPDF = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/generate-pdf', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(invoiceData),
    });
    
    if (!response.ok) throw new Error('PDF generation failed');
    
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `invoice-${invoiceData.invoiceNumber}.pdf`;
    a.click();
  } catch (error) {
    console.error('Error generating PDF:', error);
  }
};
```

## License

MIT
