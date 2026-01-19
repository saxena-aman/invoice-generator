function CorporateTemplate({ data, currencySymbol }) {
  return (
    <div className="p-12 max-w-4xl mx-auto">
      {/* Header with Blue Background */}
      <div className="bg-blue-900 text-white p-8 -m-12 mb-8">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold mb-1">INVOICE</h1>
            <p className="text-blue-200">#{data.invoiceNumber || '---'}</p>
          </div>
          <div className="text-right">
            <h2 className="text-2xl font-bold">{data.businessName || 'Your Business'}</h2>
            <p className="text-sm text-blue-200 mt-2">{data.businessEmail}</p>
            <p className="text-sm text-blue-200">{data.businessPhone}</p>
          </div>
        </div>
      </div>

      <div className="px-4">
        {/* Bill To & Invoice Info */}
        <div className="grid grid-cols-2 gap-8 mb-8">
          <div>
            <div className="bg-gray-50 p-2">
              <h3 className="text-xs font-semibold text-gray-500 uppercase mb-3">Bill To</h3>
              <p className="font-bold text-gray-900 mb-1">{data.clientName || '---'}</p>
              <p className="text-sm text-gray-600">{data.clientEmail}</p>
              <p className="text-sm text-gray-600 whitespace-pre-line">{data.clientAddress}</p>
            </div>
          </div>
          <div className="space-y-3">
            <div className="flex justify-between items-center bg-gray-50 p-3">
              <span className="text-sm font-semibold text-gray-600">Invoice Date</span>
              <span className="text-sm text-gray-900 font-medium">{data.invoiceDate || '---'}</span>
            </div>
            <div className="flex justify-between items-center bg-gray-50 p-3">
              <span className="text-sm font-semibold text-gray-600">Due Date</span>
              <span className="text-sm text-gray-900 font-medium">{data.dueDate || '---'}</span>
            </div>
            <div className="flex justify-between items-center bg-blue-50 p-3 border border-blue-200">
              <span className="text-sm font-bold text-blue-900">Amount Due</span>
              <span className="text-lg font-bold text-blue-900">{currencySymbol}{data.total.toFixed(2)}</span>
            </div>
          </div>
        </div>

        {/* Items Table */}
        <table className="w-full mb-8">
          <thead>
            <tr className="bg-gray-100">
              <th className="text-left py-3 px-4 text-xs font-bold text-gray-700 uppercase">Description</th>
              <th className="text-center py-3 px-4 text-xs font-bold text-gray-700 uppercase w-16">Qty</th>
              <th className="text-right py-3 px-4 text-xs font-bold text-gray-700 uppercase w-24">Rate</th>
              <th className="text-right py-3 px-4 text-xs font-bold text-gray-700 uppercase w-16">Tax%</th>
              <th className="text-right py-3 px-4 text-xs font-bold text-gray-700 uppercase w-16">Disc%</th>
              <th className="text-right py-3 px-4 text-xs font-bold text-gray-700 uppercase w-28">Amount</th>
            </tr>
          </thead>
          <tbody>
            {data.items.map((item, index) => (
              <tr key={index} className="border-b border-gray-200">
                <td className="py-4 px-4 text-sm text-gray-900">{item.description || '---'}</td>
                <td className="py-4 px-4 text-sm text-gray-900 text-center">{item.quantity}</td>
                <td className="py-4 px-4 text-sm text-gray-900 text-right">{currencySymbol}{item.rate.toFixed(2)}</td>
                <td className="py-4 px-4 text-sm text-gray-900 text-right">{item.taxRate || 0}%</td>
                <td className="py-4 px-4 text-sm text-gray-900 text-right">{item.discountRate || 0}%</td>
                <td className="py-4 px-4 text-sm font-semibold text-gray-900 text-right">{currencySymbol}{item.amount.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Totals */}
        <div className="flex justify-end mb-8">
          <div className="w-80">
            <div className="flex justify-between py-2 px-4">
              <span className="text-sm text-gray-600">Subtotal</span>
              <span className="text-sm text-gray-900">{currencySymbol}{data.subtotal.toFixed(2)}</span>
            </div>
            {data.discountRate > 0 && (
              <div className="flex justify-between py-2 px-4 bg-gray-50">
                <span className="text-sm text-gray-600">Discount ({data.discountRate}%)</span>
                <span className="text-sm text-gray-900">-{currencySymbol}{data.discountAmount.toFixed(2)}</span>
              </div>
            )}
            {data.taxRate > 0 && (
              <div className="flex justify-between py-2 px-4 bg-gray-50">
                <span className="text-sm text-gray-600">Tax ({data.taxRate}%)</span>
                <span className="text-sm text-gray-900">{currencySymbol}{data.taxAmount.toFixed(2)}</span>
              </div>
            )}
            <div className="flex justify-between py-3 px-4 bg-blue-900 text-white">
              <span className="font-bold">Total Due</span>
              <span className="text-xl font-bold">{currencySymbol}{data.total.toFixed(2)}</span>
            </div>
          </div>
        </div>

        {/* Notes */}
        {data.notes && (
          <div className="pt-4">
            <div className="bg-gray-50 h-1 w-full mb-3"></div>
            <h3 className="text-xs font-bold text-gray-500 uppercase mb-3">Payment Notes</h3>
            <p className="text-sm text-gray-700 whitespace-pre-line">{data.notes}</p>
          </div>
        )}
      </div>
    </div>
  );
}

export default CorporateTemplate;