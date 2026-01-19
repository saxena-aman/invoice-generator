import { useState } from 'react';
import { Plus, Trash2, Save } from 'lucide-react';
import { CURRENCIES, TEMPLATES } from '../utils/constants';
import { useToast } from '../context/ToastContext';
import { saveInvoice } from '../utils/storage';

function InvoiceForm({ invoiceData, setInvoiceData, selectedTemplate, setSelectedTemplate }) {
  const { addToast } = useToast();
  const [isSaving, setIsSaving] = useState(false);

  // Update field handler
  const updateField = (field, value) => {
    setInvoiceData(prev => ({ ...prev, [field]: value }));
  };

  // Add new line item
  const addLineItem = () => {
    const newItem = {
      id: Date.now(),
      description: '',
      quantity: 1,
      rate: 0,
      taxRate: 0,
      discountRate: 0,
      amount: 0
    };
    setInvoiceData(prev => ({
      ...prev,
      items: [...prev.items, newItem]
    }));
  };

  // Update line item
  const updateLineItem = (id, field, value) => {
    setInvoiceData(prev => {
      const updatedItems = prev.items.map(item => {
        if (item.id === id) {
          const updated = { ...item, [field]: value };
          // Recalculate amount with item-level tax and discount
          if (field === 'quantity' || field === 'rate' || field === 'taxRate' || field === 'discountRate') {
            const baseAmount = updated.quantity * updated.rate;
            const discountAmount = (baseAmount * (updated.discountRate || 0)) / 100;
            const amountAfterDiscount = baseAmount - discountAmount;
            const taxAmount = (amountAfterDiscount * (updated.taxRate || 0)) / 100;
            updated.amount = amountAfterDiscount + taxAmount;
          }
          return updated;
        }
        return item;
      });

      // Recalculate totals
      const subtotal = updatedItems.reduce((sum, item) => sum + item.amount, 0);
      const discountAmount = (subtotal * prev.discountRate) / 100;
      const amountAfterDiscount = subtotal - discountAmount;
      const taxAmount = (amountAfterDiscount * prev.taxRate) / 100;
      const total = amountAfterDiscount + taxAmount;

      return {
        ...prev,
        items: updatedItems,
        subtotal,
        discountAmount,
        taxAmount,
        total
      };
    });
  };

  // Remove line item
  const removeLineItem = (id) => {
    setInvoiceData(prev => {
      const updatedItems = prev.items.filter(item => item.id !== id);
      const subtotal = updatedItems.reduce((sum, item) => sum + item.amount, 0);
      const discountAmount = (subtotal * prev.discountRate) / 100;
      const amountAfterDiscount = subtotal - discountAmount;
      const taxAmount = (amountAfterDiscount * prev.taxRate) / 100;
      const total = amountAfterDiscount + taxAmount;

      return {
        ...prev,
        items: updatedItems,
        subtotal,
        discountAmount,
        taxAmount,
        total
      };
    });
  };

  // Update tax rate
  const updateTaxRate = (rate) => {
    const taxRate = parseFloat(rate) || 0;
    const discountAmount = (invoiceData.subtotal * invoiceData.discountRate) / 100;
    const amountAfterDiscount = invoiceData.subtotal - discountAmount;
    const taxAmount = (amountAfterDiscount * taxRate) / 100;
    const total = amountAfterDiscount + taxAmount;

    setInvoiceData(prev => ({
      ...prev,
      taxRate,
      taxAmount,
      total
    }));
  };

  // Update discount rate
  const updateDiscountRate = (rate) => {
    const discountRate = parseFloat(rate) || 0;
    const discountAmount = (invoiceData.subtotal * discountRate) / 100;
    const amountAfterDiscount = invoiceData.subtotal - discountAmount;
    const taxAmount = (amountAfterDiscount * invoiceData.taxRate) / 100;
    const total = amountAfterDiscount + taxAmount;

    setInvoiceData(prev => ({
      ...prev,
      discountRate,
      discountAmount,
      taxAmount,
      total
    }));
  };

  // Save invoice
  const handleSave = () => {
    setIsSaving(true);
    try {
      const invoiceToSave = {
        ...invoiceData,
        selectedTemplate // Save the template with the invoice
      };
      saveInvoice(invoiceToSave);
      saveInvoice(invoiceToSave);
      addToast('Invoice saved successfully!', 'success');
    } catch (error) {
      addToast('Error saving invoice', 'error');
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6 space-y-6">
      {/* Template Selector */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Template Style
        </label>
        <div className="grid grid-cols-3 gap-3">
          {TEMPLATES.map(template => (
            <button
              key={template.id}
              onClick={() => setSelectedTemplate(template.id)}
              className={`p-3 rounded-lg border-2 text-center transition-all ${
                selectedTemplate === template.id
                  ? 'border-blue-600 bg-blue-50'
                  : 'border-gray-200 hover:border-gray-300'
              }`}
            >
              <div className="font-medium text-sm">{template.name}</div>
              <div className="text-xs text-gray-500">{template.description}</div>
            </button>
          ))}
        </div>
      </div>

      {/* Business Details */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Your Business</h3>
        <div className="space-y-3">
          <input
            type="text"
            placeholder="Business Name"
            value={invoiceData.businessName}
            onChange={(e) => updateField('businessName', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <input
            type="email"
            placeholder="Email"
            value={invoiceData.businessEmail}
            onChange={(e) => updateField('businessEmail', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <input
            type="tel"
            placeholder="Phone"
            value={invoiceData.businessPhone}
            onChange={(e) => updateField('businessPhone', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <textarea
            placeholder="Business Address"
            value={invoiceData.businessAddress}
            onChange={(e) => updateField('businessAddress', e.target.value)}
            rows="2"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      {/* Client Details */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Bill To</h3>
        <div className="space-y-3">
          <input
            type="text"
            placeholder="Client Name"
            value={invoiceData.clientName}
            onChange={(e) => updateField('clientName', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <input
            type="email"
            placeholder="Client Email"
            value={invoiceData.clientEmail}
            onChange={(e) => updateField('clientEmail', e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <textarea
            placeholder="Client Address"
            value={invoiceData.clientAddress}
            onChange={(e) => updateField('clientAddress', e.target.value)}
            rows="2"
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>

      {/* Invoice Details */}
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Invoice Details</h3>
        <div className="grid grid-cols-2 gap-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Invoice Number
            </label>
            <input
              type="text"
              placeholder="INV-001"
              value={invoiceData.invoiceNumber}
              onChange={(e) => updateField('invoiceNumber', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Currency
            </label>
            <select
              value={invoiceData.currency}
              onChange={(e) => updateField('currency', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              {CURRENCIES.map(curr => (
                <option key={curr.code} value={curr.code}>
                  {curr.symbol} {curr.code}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Invoice Date
            </label>
            <input
              type="date"
              value={invoiceData.invoiceDate}
              onChange={(e) => updateField('invoiceDate', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Due Date
            </label>
            <input
              type="date"
              value={invoiceData.dueDate}
              onChange={(e) => updateField('dueDate', e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
      </div>

      {/* Line Items */}
      <div>
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-lg font-semibold text-gray-900">Items</h3>
          <button
            onClick={addLineItem}
            className="flex items-center space-x-1 px-3 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
          >
            <Plus className="w-4 h-4" />
            <span>Add Item</span>
          </button>
        </div>

        <div className="space-y-3">
          {invoiceData.items.map((item, index) => (
            <div key={item.id} className="border border-gray-200 rounded-lg p-3">
              <div className="flex items-start space-x-2">
                <div className="flex-1 space-y-2">
                  {/* Description Field */}
                  <div>
                    <label className="block text-xs font-medium text-gray-700 mb-1">
                      Description
                    </label>
                    <input
                      type="text"
                      placeholder="Item description"
                      value={item.description}
                      onChange={(e) => updateLineItem(item.id, 'description', e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>
                  
                  
                  {/* Qty, Rate, Tax, Discount, Amount Fields */}
                  <div className="grid grid-cols-5 gap-2">
                    {/* Quantity */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 mb-1">
                        Qty
                      </label>
                      <input
                        type="text"
                        inputMode="numeric"
                        placeholder="0"
                        value={item.quantity}
                        onChange={(e) => updateLineItem(item.id, 'quantity', parseFloat(e.target.value) || 0)}
                        onFocus={(e) => e.target.select()}
                        className="w-full px-2 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                      />
                    </div>
                    
                    {/* Rate */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 mb-1">
                        Rate
                      </label>
                      <input
                        type="text"
                        inputMode="decimal"
                        placeholder="0.00"
                        value={item.rate}
                        onChange={(e) => updateLineItem(item.id, 'rate', parseFloat(e.target.value) || 0)}
                        onFocus={(e) => e.target.select()}
                        className="w-full px-2 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                      />
                    </div>

                    {/* Tax % */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 mb-1">
                        Tax %
                      </label>
                      <input
                        type="text"
                        inputMode="decimal"
                        placeholder="0"
                        value={item.taxRate || 0}
                        onChange={(e) => updateLineItem(item.id, 'taxRate', parseFloat(e.target.value) || 0)}
                        onFocus={(e) => e.target.select()}
                        className="w-full px-2 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                      />
                    </div>

                    {/* Discount % */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 mb-1">
                        Disc %
                      </label>
                      <input
                        type="text"
                        inputMode="decimal"
                        placeholder="0"
                        value={item.discountRate || 0}
                        onChange={(e) => updateLineItem(item.id, 'discountRate', parseFloat(e.target.value) || 0)}
                        onFocus={(e) => e.target.select()}
                        className="w-full px-2 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                      />
                    </div>
                    
                    {/* Amount */}
                    <div>
                      <label className="block text-xs font-medium text-gray-700 mb-1">
                        Amount
                      </label>
                      <input
                        type="text"
                        placeholder="0.00"
                        value={item.amount.toFixed(2)}
                        readOnly
                        className="w-full px-2 py-2 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 text-sm"
                      />
                    </div>
                  </div>
                </div>
                {invoiceData.items.length > 1 && (
                  <button
                    onClick={() => removeLineItem(item.id)}
                    className="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Totals */}
      <div className="border-t pt-4">
        <div className="space-y-2 max-w-xs ml-auto">
          <div className="flex justify-between text-sm">
            <span className="text-gray-600">Subtotal:</span>
            <span className="font-medium">{CURRENCIES.find(c => c.code === invoiceData.currency)?.symbol}{invoiceData.subtotal.toFixed(2)}</span>
          </div>
          
          {/* Discount */}
          <div>
            <div className="flex justify-between items-center text-sm mb-1">
              <label className="text-gray-700 font-medium">Discount %:</label>
              <input
                type="text"
                inputMode="decimal"
                value={invoiceData.discountRate}
                onChange={(e) => updateDiscountRate(e.target.value)}
                onFocus={(e) => e.target.select()}
                className="w-20 px-2 py-1 border border-gray-300 rounded text-right text-sm"
                placeholder="0"
              />
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Discount Amount:</span>
              <span className="font-medium">{CURRENCIES.find(c => c.code === invoiceData.currency)?.symbol}{invoiceData.discountAmount.toFixed(2)}</span>
            </div>
          </div>

          {/* Tax */}
          <div>
            <div className="flex justify-between items-center text-sm mb-1">
              <label className="text-gray-700 font-medium">Tax %:</label>
              <input
                type="text"
                inputMode="decimal"
                value={invoiceData.taxRate}
                onChange={(e) => updateTaxRate(e.target.value)}
                onFocus={(e) => e.target.select()}
                className="w-20 px-2 py-1 border border-gray-300 rounded text-right text-sm"
                placeholder="0"
              />
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-600">Tax Amount:</span>
              <span className="font-medium">{CURRENCIES.find(c => c.code === invoiceData.currency)?.symbol}{invoiceData.taxAmount.toFixed(2)}</span>
            </div>
          </div>

          <div className="flex justify-between text-lg font-bold pt-2 border-t">
            <span>Total:</span>
            <span>{CURRENCIES.find(c => c.code === invoiceData.currency)?.symbol}{invoiceData.total.toFixed(2)}</span>
          </div>
        </div>
      </div>

      {/* Notes */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Notes
        </label>
        <textarea
          placeholder="Additional notes or payment instructions..."
          value={invoiceData.notes}
          onChange={(e) => updateField('notes', e.target.value)}
          rows="3"
          className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>

      {/* Save Button */}
      <button
        onClick={handleSave}
        disabled={isSaving}
        className="w-full flex items-center justify-center space-x-2 px-4 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <Save className="w-5 h-5" />
        <span>{isSaving ? 'Saving...' : 'Save Invoice'}</span>
      </button>
    </div>
  );
}

export default InvoiceForm;