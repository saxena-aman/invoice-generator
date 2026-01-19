// Save invoice to localStorage
export const saveInvoice = (invoice) => {
    const invoices = getInvoices();

    if (invoice.id) {
        const index = invoices.findIndex(inv => inv.id === invoice.id);
        if (index !== -1) {
            // Update existing
            invoices[index] = { ...invoices[index], ...invoice };
            localStorage.setItem('invoices', JSON.stringify(invoices));
            return invoices[index];
        }
    }

    // Create new
    const newInvoice = {
        ...invoice,
        id: invoice.id || Date.now(),
        createdAt: invoice.createdAt || new Date().toISOString()
    };

    invoices.push(newInvoice);
    localStorage.setItem('invoices', JSON.stringify(invoices));
    return newInvoice;
};

// Get all invoices
export const getInvoices = () => {
    const data = localStorage.getItem('invoices');
    return data ? JSON.parse(data) : [];
};

// Get single invoice
export const getInvoice = (id) => {
    const invoices = getInvoices();
    return invoices.find(inv => inv.id === id);
};

// Update invoice
export const updateInvoice = (id, updatedData) => {
    const invoices = getInvoices();
    const index = invoices.findIndex(inv => inv.id === id);

    if (index !== -1) {
        invoices[index] = { ...invoices[index], ...updatedData };
        localStorage.setItem('invoices', JSON.stringify(invoices));
        return invoices[index];
    }
    return null;
};

// Delete invoice
export const deleteInvoice = (id) => {
    const invoices = getInvoices();
    const filtered = invoices.filter(inv => inv.id !== id);
    localStorage.setItem('invoices', JSON.stringify(filtered));
};

// Export all data
export const exportData = () => {
    const data = {
        invoices: getInvoices(),
        businessInfo: localStorage.getItem('businessInfo'),
        exportedAt: new Date().toISOString()
    };

    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `invoices-backup-${Date.now()}.json`;
    a.click();
    URL.revokeObjectURL(url);
};

// Import data
export const importData = (file) => {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                const data = JSON.parse(e.target.result);
                if (data.invoices) {
                    localStorage.setItem('invoices', JSON.stringify(data.invoices));
                }
                if (data.businessInfo) {
                    localStorage.setItem('businessInfo', data.businessInfo);
                }
                resolve(data);
            } catch (error) {
                reject(error);
            }
        };
        reader.onerror = reject;
        reader.readAsText(file);
    });
};