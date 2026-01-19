import { useState, useEffect } from 'react';
import { FileText, Search, Trash2, Edit, Eye, Calendar, DollarSign } from 'lucide-react';
import { getInvoices, deleteInvoice } from '../utils/storage';
import { CURRENCIES } from '../utils/constants';
import { useDialog } from '../context/DialogContext';
import { useToast } from '../context/ToastContext';

function InvoiceList({ setCurrentView, setInvoiceData, setSelectedTemplate }) {
  const { confirm } = useDialog();
  const { addToast } = useToast();
  const [invoices, setInvoices] = useState([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredInvoices, setFilteredInvoices] = useState([]);

  // Load invoices on mount
  useEffect(() => {
    loadInvoices();
  }, []);

  // Filter invoices when search term changes
  useEffect(() => {
    if (searchTerm.trim() === '') {
      setFilteredInvoices(invoices);
    } else {
      const filtered = invoices.filter(invoice => 
        invoice.invoiceNumber?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        invoice.clientName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        invoice.businessName?.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setFilteredInvoices(filtered);
    }
  }, [searchTerm, invoices]);

  const loadInvoices = () => {
    const data = getInvoices();
    // Sort by creation date (newest first)
    const sorted = data.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));
    setInvoices(sorted);
    setFilteredInvoices(sorted);
  };

  const handleDelete = async (id) => {
    const isConfirmed = await confirm({
      title: 'Delete Invoice',
      message: 'Are you sure you want to delete this invoice? This action cannot be undone.',
      type: 'danger',
      confirmText: 'Delete',
      cancelText: 'Cancel'
    });

    if (isConfirmed) {
      deleteInvoice(id);
      loadInvoices();
      addToast('Invoice deleted successfully', 'success');
    }
  };

  const handleEdit = (invoice) => {
    setInvoiceData(invoice);
    // Restore the template if it was saved with the invoice
    if (invoice.selectedTemplate && setSelectedTemplate) {
      setSelectedTemplate(invoice.selectedTemplate);
    }
    setCurrentView('form');
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
  };

  const getCurrencySymbol = (code) => {
    return CURRENCIES.find(c => c.code === code)?.symbol || '$';
  };

  // Empty state
  if (invoices.length === 0) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-md mx-auto text-center">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gray-100 rounded-full mb-6">
            <FileText className="w-10 h-10 text-gray-400" />
          </div>
          <h2 className="text-2xl font-bold text-gray-900 mb-2">No Invoices Yet</h2>
          <p className="text-gray-600 mb-6">
            Create your first invoice to get started with managing your billing.
          </p>
          <button
            onClick={() => setCurrentView('form')}
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            Create New Invoice
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-8">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-3xl font-bold text-gray-900">My Invoices</h2>
            <p className="text-gray-600 mt-1">{invoices.length} invoice{invoices.length !== 1 ? 's' : ''} total</p>
          </div>
          <button
            onClick={() => setCurrentView('form')}
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center space-x-2"
          >
            <FileText className="w-5 h-5" />
            <span>New Invoice</span>
          </button>
        </div>

        {/* Search Bar */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Search by invoice number, client, or business name..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      {/* Invoice Grid */}
      {filteredInvoices.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500">No invoices match your search.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredInvoices.map((invoice) => (
            <div
              key={invoice.id}
              className="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow"
            >
              {/* Card Header */}
              <div className="p-6 border-b border-gray-100">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <h3 className="text-lg font-semibold text-gray-900">
                      {invoice.invoiceNumber || 'Draft'}
                    </h3>
                    <p className="text-sm text-gray-600 mt-1">
                      {invoice.clientName || 'No client'}
                    </p>
                  </div>
                  <div className="flex items-center space-x-1">
                    <button
                      onClick={() => handleEdit(invoice)}
                      className="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                      title="Edit invoice"
                    >
                      <Edit className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => handleDelete(invoice.id)}
                      className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                      title="Delete invoice"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                {/* Amount */}
                <div className="flex items-center space-x-2 text-2xl font-bold text-gray-900">
                  <DollarSign className="w-6 h-6 text-green-600" />
                  <span>
                    {getCurrencySymbol(invoice.currency)}
                    {(invoice.total || 0).toFixed(2)}
                  </span>
                </div>
              </div>

              {/* Card Body */}
              <div className="p-6 space-y-3">
                {/* Dates */}
                <div className="flex items-center space-x-2 text-sm text-gray-600">
                  <Calendar className="w-4 h-4" />
                  <span>Issued: {formatDate(invoice.invoiceDate)}</span>
                </div>
                {invoice.dueDate && (
                  <div className="flex items-center space-x-2 text-sm text-gray-600">
                    <Calendar className="w-4 h-4" />
                    <span>Due: {formatDate(invoice.dueDate)}</span>
                  </div>
                )}

                {/* Business Info */}
                {invoice.businessName && (
                  <div className="pt-3 border-t border-gray-100">
                    <p className="text-xs text-gray-500 uppercase tracking-wide mb-1">From</p>
                    <p className="text-sm font-medium text-gray-900">{invoice.businessName}</p>
                  </div>
                )}

                {/* Items Count */}
                <div className="pt-3 border-t border-gray-100">
                  <p className="text-xs text-gray-500">
                    {invoice.items?.length || 0} item{invoice.items?.length !== 1 ? 's' : ''}
                  </p>
                </div>
              </div>

              {/* Card Footer */}
              <div className="px-6 py-4 bg-gray-50 border-t border-gray-100">
                <button
                  onClick={() => handleEdit(invoice)}
                  className="w-full flex items-center justify-center space-x-2 px-4 py-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                >
                  <Eye className="w-4 h-4" />
                  <span className="font-medium">View & Edit</span>
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default InvoiceList;
