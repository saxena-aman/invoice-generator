# Invoice Generator

A complete full-stack application for generating professional PDF invoices, featuring a React frontend and a robust Go backend.

## 📂 Project Structure

- **`ui/`**: React Frontend (Vite) - Handles the user interface, form, and preview.
- **`invoicer/`**: Go Backend - API service that generates high-quality PDFs.

## 🚀 Getting Started

### Prerequisites

- **Node.js** (for Frontend)
- **Go** (for Backend)

### 1. Start the Backend (API)

The backend handles PDF generation and runs on port `8080`.

```bash
cd invoicer
go run main.go
# Server will start at http://localhost:8080
```

### 2. Start the Frontend (UI)

The frontend enables you to create and manage invoices and runs on port `5173`.

```bash
cd ui
npm install
npm run dev
# Open http://localhost:5173
```

## ✨ Features

- **Professional PDF Generation:** Server-side PDF generation using Go for perfect layout control.
- **Multiple Templates:**
  - **Minimal:** Clean and simple.
  - **Corporate:** Professional layout.
  - **Modern:** Stylish design with accent colors.
- **Data Management:**
  - **Smart Save:** Auto-saves drafts.
  - **Backup/Restore:** Export your invoices to JSON and restore them anytime.
- **UX Enhancements:**
  - Beautiful Toast Notifications.
  - Custom Confirmation Dialogs.
  - Real-time Preview.

## 🛠️ Tech Stack

- **Frontend:** React, Vite, TailwindCSS, Lucide Icons
- **Backend:** Go (Golang), gofpdf, Gorilla Mux
