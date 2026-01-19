import { Download, Loader2 } from 'lucide-react';
import { CURRENCIES } from '../utils/constants';
import MinimalTemplate from '../templates/MinimalTemplate';
import CorporateTemplate from '../templates/CorporateTemplate';
import ModernTemplate from '../templates/ModernTemplate';
import { useState } from 'react';

import { useToast } from '../context/ToastContext';

function InvoicePreview({ invoiceData, selectedTemplate }) {
  const { addToast } = useToast();
  const [isGenerating, setIsGenerating] = useState(false);

  // Get currency symbol
  const currencySymbol = CURRENCIES.find(c => c.code === invoiceData.currency)?.symbol || '$';

  // Download PDF
  const downloadPDF = async () => {
    setIsGenerating(true);
    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL}/generate-pdf`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...invoiceData,
          selectedTemplate,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to generate PDF');
      }

      const blob = await response.blob();
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
    } catch (error) {
      console.error('Error generating PDF:', error);
      addToast('Failed to generate PDF. Please try again.', 'error');
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