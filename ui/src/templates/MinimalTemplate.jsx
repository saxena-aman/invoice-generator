function MinimalTemplate({ data, currencySymbol }) {
  return (
    <div className="p-12 max-w-4xl mx-auto">
      {/* Header */}
      <div className="flex justify-between items-start mb-12">
        <div>
          <h1 className="text-4xl font-bold text-gray-900 mb-2">INVOICE</h1>
          <p className="text-gray-600">#{data.invoiceNumber || '---'}</p>
        </div>
        <div className="text-right">
          <h2 className="text-xl font-bold text-gray-900">{data.businessName || 'Your Business'}</h2>
          <p className="text-sm text-gray-600 mt-1">{data.businessEmail}</p>
          <p className="text-sm text-gray-600">{data.businessPhone}</p>
          <p className="text-sm text-gray-600 whitespace-pre-line">{data.businessAddress}</p>
        </div>
      </div>

      {/* Bill To & Dates */}
      <div className="grid grid-cols-2 gap-8 mb-12">
        <div>
          <h3 className="text-sm font-semibold text-gray-500 uppercase mb-2">Bill To</h3>
          <p className="font-semibold text-gray-900">{data.clientName || '---'}</p>
          <p className="text-sm text-gray-600">{data.clientEmail}</p>
          <p className="text-sm text-gray-600 whitespace-pre-line">{data.clientAddress}</p>
        </div>
        <div>
          <div className="flex justify-between mb-2">
            <span className="text-sm font-semibold text-gray-500">Invoice Date:</span>
            <span className="text-sm text-gray-900">{data.invoiceDate || '---'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-sm font-semibold text-gray-500">Due Date:</span>
            <span className="text-sm text-gray-900">{data.dueDate || '---'}</span>
          </div>
        </div>
      </div>

      {/* Items Table */}
      <table className="w-full mb-8">
        <thead>
          <tr className="border-b-2 border-gray-900">
            <th className="text-left py-3 text-sm font-semibold text-gray-900 uppercase">Description</th>
            <th className="text-right py-3 text-sm font-semibold text-gray-900 uppercase w-16">Qty</th>
            <th className="text-right py-3 text-sm font-semibold text-gray-900 uppercase w-24">Rate</th>
            <th className="text-right py-3 text-sm font-semibold text-gray-900 uppercase w-16">Tax%</th>
            <th className="text-right py-3 text-sm font-semibold text-gray-900 uppercase w-16">Disc%</th>
            <th className="text-right py-3 text-sm font-semibold text-gray-900 uppercase w-28">Amount</th>
          </tr>
        </thead>
        <tbody>
          {data.items.map((item, index) => (
            <tr key={index} className="border-b border-gray-200">
              <td className="py-3 text-sm text-gray-900">{item.description || '---'}</td>
              <td className="py-3 text-sm text-gray-900 text-right">{item.quantity}</td>
              <td className="py-3 text-sm text-gray-900 text-right">{currencySymbol}{item.rate.toFixed(2)}</td>
              <td className="py-3 text-sm text-gray-900 text-right">{item.taxRate || 0}%</td>
              <td className="py-3 text-sm text-gray-900 text-right">{item.discountRate || 0}%</td>
              <td className="py-3 text-sm text-gray-900 text-right">{currencySymbol}{item.amount.toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Totals */}
      <div className="flex justify-end mb-12">
        <div className="w-80">
          <div className="flex justify-between py-2">
            <span className="text-sm text-gray-600">Subtotal:</span>
            <span className="text-sm text-gray-900">{currencySymbol}{data.subtotal.toFixed(2)}</span>
          </div>
          {data.discountRate > 0 && (
            <div className="flex justify-between py-2">
              <span className="text-sm text-gray-600">Discount ({data.discountRate}%):</span>
              <span className="text-sm text-gray-900">-{currencySymbol}{data.discountAmount.toFixed(2)}</span>
            </div>
          )}
          {data.taxRate > 0 && (
            <div className="flex justify-between py-2">
              <span className="text-sm text-gray-600">Tax ({data.taxRate}%):</span>
              <span className="text-sm text-gray-900">{currencySymbol}{data.taxAmount.toFixed(2)}</span>
            </div>
          )}
          <div className="flex justify-between py-3 border-t-2 border-gray-900">
            <span className="text-lg font-bold text-gray-900">Total:</span>
            <span className="text-lg font-bold text-gray-900">{currencySymbol}{data.total.toFixed(2)}</span>
          </div>
        </div>
      </div>

      {/* Notes */}
      {data.notes && (
        <div className="border-t pt-8">
          <h3 className="text-sm font-semibold text-gray-500 uppercase mb-2">Notes</h3>
          <p className="text-sm text-gray-600 whitespace-pre-line">{data.notes}</p>
        </div>
      )}
    </div>
  );
}

export default MinimalTemplate;