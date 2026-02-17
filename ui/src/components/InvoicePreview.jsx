import { Download, Loader2 } from 'lucide-react';
import { CURRENCIES } from '../utils/constants';
import MinimalTemplate from '../templates/MinimalTemplate';
import CorporateTemplate from '../templates/CorporateTemplate';
import ModernTemplate from '../templates/ModernTemplate';
import { useState } from 'react';

import { useToast } from '../context/ToastContext';
import { useAuth } from '../context/AuthContext';
import { downloadInvoicePDF } from '../utils/pdfGenerator';

function InvoicePreview({ invoiceData, selectedTemplate }) {
  const { addToast } = useToast();
  const { token } = useAuth();
  const [isGenerating, setIsGenerating] = useState(false);

  // Get currency symbol
  const currencySymbol = CURRENCIES.find(c => c.code === invoiceData.currency)?.symbol || '$';

  // Download PDF
  const downloadPDF = async () => {
    setIsGenerating(true);
    try {
      await downloadInvoicePDF(invoiceData, selectedTemplate, token);
    } catch (error) {
      console.error('Error generating PDF:', error);
      addToast(error.message || 'Failed to generate PDF. Please try again.', 'error');
    } finally {
      setIsGenerating(false);
    }
  };

  // Select template component
  const TemplateComponent = () => {
    switch (selectedTemplate) {
      case 'corporate':
        return <CorporateTemplate data={invoiceData} currencySymbol={currencySymbol} />;
      case 'modern':
        return <ModernTemplate data={invoiceData} currencySymbol={currencySymbol} />;
      default:
        return <MinimalTemplate data={invoiceData} currencySymbol={currencySymbol} />;
    }
  };

  return (
    <div className="space-y-4">
      {/* Download Button */}
      <div className="flex justify-end">
        <button
          onClick={downloadPDF}
          disabled={isGenerating}
          className="flex items-center space-x-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isGenerating ? (
            <Loader2 className="w-4 h-4 animate-spin" />
          ) : (
            <Download className="w-4 h-4" />
          )}
          <span>{isGenerating ? 'Generating PDF...' : 'Download PDF'}</span>
        </button>
      </div>

      {/* Preview Container */}
      <div className="bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="bg-white">
          <TemplateComponent />
        </div>
      </div>
    </div>
  );
}

export default InvoicePreview;