import { FileText, List, Download, Upload } from 'lucide-react';
import { exportData, importData } from '../utils/storage';
import { useRef } from 'react';
import { useToast } from '../context/ToastContext';
import Tooltip from './Tooltip';

function Header({ currentView, setCurrentView, setInvoiceData, setSelectedTemplate, onNewInvoice }) {
  const { addToast } = useToast();
  const fileInputRef = useRef(null);

  const handleExport = () => {
    try {
      exportData();
    } catch (error) {
      addToast('Error exporting data: ' + error.message, 'error');
    }
  };

  const handleImport = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (event) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      const data = await importData(file);
      
      if (data.invoices && Array.isArray(data.invoices)) {
        // It's a backup file
        addToast('Data imported successfully!', 'success');
        
        if (data.invoices.length === 1) {
          // If it's just one invoice in the backup, load it immediately
          const invoice = data.invoices[0];
          setInvoiceData(prev => ({
            ...prev,
            ...invoice
          }));
          if (invoice.selectedTemplate) {
            setSelectedTemplate(invoice.selectedTemplate);
          }
          setCurrentView('form');
        } else {
          // Multiple invoices, go to list
          if (currentView === 'list') {
             setTimeout(() => window.location.reload(), 1000);
          } else {
             setCurrentView('list');
          }
        }
      } else {
        // It's likely a single invoice file (not backup format)
        setInvoiceData(prev => ({
          ...prev,
          ...data
        }));
        if (data.selectedTemplate) {
          setSelectedTemplate(data.selectedTemplate);
        }
        setCurrentView('form');
        addToast('Invoice loaded into form!', 'success');
      }
    } catch (error) {
      addToast('Error importing data: ' + error.message, 'error');
    }

    // Reset input
    event.target.value = '';
  };

  return (
    <header className="bg-white shadow-sm border-b">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <FileText className="w-8 h-8 text-blue-600" />
            <h1 className="text-2xl font-bold text-gray-900">InvoiceGen</h1>
          </div>
          
          <nav className="flex items-center space-x-4">
            <button
              onClick={onNewInvoice}
              className={`px-4 py-2 rounded-lg flex items-center space-x-2 ${
                currentView === 'form' 
                  ? 'bg-blue-600 text-white' 
                  : 'text-gray-600 hover:bg-gray-100'
              }`}
            >
              <FileText className="w-4 h-4" />
              <span>New Invoice</span>
            </button>
            
            <button
              onClick={() => setCurrentView('list')}
              className={`px-4 py-2 rounded-lg flex items-center space-x-2 ${
                currentView === 'list' 
                  ? 'bg-blue-600 text-white' 
                  : 'text-gray-600 hover:bg-gray-100'
              }`}
            >
              <List className="w-4 h-4" />
              <span>My Invoices</span>
            </button>
            
            <Tooltip content="Download Invoice Template">
              <button 
                onClick={handleExport}
                className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg"
              >
                <Download className="w-5 h-5" />
              </button>
            </Tooltip>
            
            <Tooltip content="Upload Existing Invoice Template">
              <button 
                onClick={handleImport}
                className="p-2 text-gray-600 hover:bg-gray-100 rounded-lg"
              >
                <Upload className="w-5 h-5" />
              </button>
            </Tooltip>

            {/* Hidden file input for import */}
            <input
              ref={fileInputRef}
              type="file"
              accept=".json"
              onChange={handleFileChange}
              className="hidden"
            />
          </nav>
        </div>
      </div>
    </header>
  );
}

export default Header;