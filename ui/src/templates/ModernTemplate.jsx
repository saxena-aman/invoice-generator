function ModernTemplate({ data, currencySymbol }) {
  return (
    <div className="p-4 max-w-4xl mx-auto bg-purple-50">
      {/* Header */}
      <div className="mb-8">
        <div className="flex justify-between items-start mb-8">
          <div>
            <div className="inline-block bg-purple-600 text-white px-6 py-3 rounded-lg mb-2">
              <h1 className="text-3xl font-bold">INVOICE</h1>
            </div>
            <p className="text-gray-700 font-semibold mt-2">#{data.invoiceNumber || '---'}</p>
          </div>
          <div className="text-right">
            <h2 className="text-2xl font-bold text-purple-600">
              {data.businessName || 'Your Business'}
            </h2>
            <p className="text-sm text-gray-600 mt-2">{data.businessEmail}</p>
            <p className="text-sm text-gray-600">{data.businessPhone}</p>
          </div>
        </div>

        {/* Bill To & Dates Cards */}
        <div className="grid grid-cols-2 gap-6 mb-8">
          <div className="bg-white rounded-r-lg p-6 shadow-sm border border-purple-100 relative">
            <div className="absolute left-0 top-0 bottom-0 w-1 bg-purple-600"></div>
            <div className="flex items-center mb-3">
              <h3 className="text-sm font-bold text-gray-900 uppercase">Bill To</h3>
            </div>
            <p className="font-bold text-gray-900 mb-1">{data.clientName || '---'}</p>
            <p className="text-sm text-gray-600">{data.clientEmail}</p>
            <p className="text-sm text-gray-600 whitespace-pre-line">{data.clientAddress}</p>
          </div>

          <div className="bg-white rounded-lg p-6 shadow-sm border border-purple-100 space-y-3">
            <div className="flex justify-between">
              <span className="text-sm font-semibold text-gray-600">Invoice Date</span>
              <span className="text-sm text-gray-900 font-medium">{data.invoiceDate || '---'}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm font-semibold text-gray-600">Due Date</span>
              <span className="text-sm text-gray-900 font-medium">{data.dueDate || '---'}</span>
            </div>
            <div className="pt-3 border-t border-gray-200">
              <div className="flex justify-between items-center">
                <span className="text-sm font-bold text-gray-900">Amount Due</span>
                <span className="text-2xl font-bold text-purple-600">
                  {currencySymbol}{data.total.toFixed(2)}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Items Table */}
      <div className="bg-white rounded-b-lg shadow-sm border border-gray-100 overflow-hidden mb-6">
        <table className="w-full">
          <thead>
            <tr className="bg-purple-600 text-white">
              <th className="text-left py-4 px-6 text-xs font-bold uppercase">Description</th>
              <th className="text-center py-4 px-6 text-xs font-bold uppercase w-16">Qty</th>
              <th className="text-right py-4 px-6 text-xs font-bold uppercase w-24">Rate</th>
              <th className="text-right py-4 px-6 text-xs font-bold uppercase w-16">Tax%</th>
              <th className="text-right py-4 px-6 text-xs font-bold uppercase w-16">Disc%</th>
              <th className="text-right py-4 px-6 text-xs font-bold uppercase w-28">Amount</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((item, index) => (
              <tr key={index} className="border-b border-gray-100 hover:bg-purple-50 transition-colors">
                <td className="py-4 px-6 text-sm text-gray-900">{item.description || '---'}</td>
                <td className="py-4 px-6 text-sm text-gray-900 text-center">{item.quantity}</td>
                <td className="py-4 px-6 text-sm text-gray-900 text-right">{currencySymbol}{item.rate.toFixed(2)}</td>
                <td className="py-4 px-6 text-sm text-gray-900 text-right">{item.taxRate || 0}%</td>
                <td className="py-4 px-6 text-sm text-gray-900 text-right">{item.discountRate || 0}%</td>
                <td className="py-4 px-6 text-sm font-semibold text-gray-900 text-right">{currencySymbol}{item.amount.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Totals */}
      <div className="flex justify-end mb-8">
        <div className="w-96 bg-white rounded-lg shadow-sm border border-gray-100 overflow-hidden">
          <div className="p-6 space-y-3">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600">Subtotal</span>
              <span className="text-sm text-gray-900 font-medium">{currencySymbol}{data.subtotal.toFixed(2)}</span>
            </div>
            {data.discountRate > 0 && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Discount ({data.discountRate}%)</span>
                <span className="text-sm text-gray-900 font-medium">-{currencySymbol}{data.discountAmount.toFixed(2)}</span>
              </div>
            )}
            {data.taxRate > 0 && (
              <div className="flex justify-between">
                <span className="text-sm text-gray-600">Tax ({data.taxRate}%)</span>
                <span className="text-sm text-gray-900 font-medium">{currencySymbol}{data.taxAmount.toFixed(2)}</span>
              </div>
            )}
          </div>
          <div className="bg-purple-600 text-white px-6 py-4 flex justify-between items-center">
            <span className="font-bold text-lg">Total</span>
            <span className="text-2xl font-bold">{currencySymbol}{data.total.toFixed(2)}</span>
          </div>
        </div>
      </div>

      {/* Notes */}
      {data.notes && (
        <div className="bg-white rounded-r-lg p-6 shadow-sm border border-gray-100 relative">
          <div className="absolute left-0 top-0 bottom-0 w-1 bg-purple-600"></div>
          <div className="flex items-center mb-3">
            <h3 className="text-xs font-bold text-gray-900 uppercase">Payment Notes</h3>
          </div>
          <p className="text-sm text-gray-700 whitespace-pre-line">{data.notes}</p>
        </div>
      )}
    </div>
  );
}

export default ModernTemplate;