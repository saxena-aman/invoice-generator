import { useState } from 'react';
import Header from './components/Header';
import InvoiceForm from './components/InvoiceForm';
import InvoicePreview from './components/InvoicePreview';
import InvoiceList from './components/InvoiceList';
import LoginPage from './components/LoginPage';
import SignupPage from './components/SignupPage';
import { saveInvoice } from './utils/storage';
import { useToast } from './context/ToastContext';
import { useAuth } from './context/AuthContext';

const INITIAL_STATE = {
  businessName: '',
  businessEmail: '',
  businessAddress: '',
  businessPhone: '',
  clientName: '',
  clientEmail: '',
  clientAddress: '',
  invoiceNumber: '',
  dueDate: '',
  items: [
    { id: 1, description: '', quantity: 1, rate: 0, taxRate: 0, discountRate: 0, amount: 0 }
  ],
  subtotal: 0,
  taxRate: 0,
  taxAmount: 0,
  discountRate: 0,
  discountAmount: 0,
  total: 0,
  currency: 'USD',
  notes: '',
  paymentTerms: ''
};

function App() {
  const { addToast } = useToast();
  const { isAuthenticated } = useAuth();
  const [authView, setAuthView] = useState('login'); // 'login' or 'signup'
  const [currentView, setCurrentView] = useState('form'); // 'form' or 'list'
  const [selectedTemplate, setSelectedTemplate] = useState('minimal');
  const [invoiceData, setInvoiceData] = useState({
    ...INITIAL_STATE,
    invoiceDate: new Date().toISOString().split('T')[0]
  });

  const handleNewInvoice = () => {
    // Save current invoice
    const invoiceToSave = {
      ...invoiceData,
      selectedTemplate
    };
    saveInvoice(invoiceToSave);
    
    // Reset to fresh state
    setInvoiceData({
      ...INITIAL_STATE,
      invoiceDate: new Date().toISOString().split('T')[0]
    });
    setSelectedTemplate('minimal');
    setCurrentView('form');
    addToast('Current invoice saved. Started a new invoice.', 'success');
  };

  // ── Auth gate ───────────────────────────────────────────────────────
  if (!isAuthenticated) {
    return authView === 'login' ? (
      <LoginPage onSwitchToSignup={() => setAuthView('signup')} />
    ) : (
      <SignupPage onSwitchToLogin={() => setAuthView('login')} />
    );
  }

  // ── Authenticated UI ───────────────────────────────────────────────
  return (
    <div className="min-h-screen bg-gray-50">
      <Header 
        currentView={currentView}
        setCurrentView={setCurrentView}
        setInvoiceData={setInvoiceData}
        setSelectedTemplate={setSelectedTemplate}
        onNewInvoice={handleNewInvoice}
      />
      
      {currentView === 'form' ? (
        <div className="container mx-auto px-4 py-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Left: Form */}
            <InvoiceForm 
              invoiceData={invoiceData}
              setInvoiceData={setInvoiceData}
              selectedTemplate={selectedTemplate}
              setSelectedTemplate={setSelectedTemplate}
            />
            
            {/* Right: Preview */}
            <InvoicePreview 
              invoiceData={invoiceData}
              selectedTemplate={selectedTemplate}
            />
          </div>
        </div>
      ) : (
        <InvoiceList 
          setCurrentView={setCurrentView}
          setInvoiceData={setInvoiceData}
          setSelectedTemplate={setSelectedTemplate}
        />
      )}
    </div>
  );
}

export default App;