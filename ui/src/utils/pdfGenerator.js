const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

/**
 * Generate a PDF invoice via the backend API.
 * @param {Object} invoiceData - The invoice form data.
 * @param {string} selectedTemplate - Template name (minimal, corporate, modern).
 * @param {string} token - JWT access token.
 * @returns {Promise<Blob>} The PDF blob.
 */
export async function generatePDF(invoiceData, selectedTemplate, token) {
    const response = await fetch(`${API_URL}/generate-pdf`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
            ...invoiceData,
            selectedTemplate,
        }),
    });

    if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.message || 'Failed to generate PDF');
    }

    return response.blob();
}

/**
 * Generate and trigger a PDF download.
 * @param {Object} invoiceData - The invoice form data.
 * @param {string} selectedTemplate - Template name.
 * @param {string} token - JWT access token.
 */
export async function downloadInvoicePDF(invoiceData, selectedTemplate, token) {
    const blob = await generatePDF(invoiceData, selectedTemplate, token);

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    const clientName = invoiceData.clientName || 'Client';
    const safeClientName = clientName.replace(/[^a-z0-9]/gi, '_').replace(/_+/g, '_');
    a.download = `Invoice_for_${safeClientName}.pdf`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
}
