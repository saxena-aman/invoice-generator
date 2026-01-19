import React, { createContext, useContext, useState, useRef } from 'react';
import { AlertTriangle, HelpCircle, X, Check } from 'lucide-react';

const DialogContext = createContext(null);

export const DialogProvider = ({ children }) => {
  const [dialog, setDialog] = useState({
    isOpen: false,
    title: '',
    message: '',
    type: 'confirm', // confirm, danger, info
    confirmText: 'Confirm',
    cancelText: 'Cancel',
    onConfirm: () => {},
    onCancel: () => {},
  });

  const resolver = useRef(null);

  const confirm = (options) => {
    return new Promise((resolve) => {
      resolver.current = resolve;
      setDialog({
        isOpen: true,
        title: options.title || 'Are you sure?',
        message: options.message || '',
        type: options.type || 'confirm',
        confirmText: options.confirmText || 'Confirm',
        cancelText: options.cancelText || 'Cancel',
        onConfirm: () => {
          setDialog(prev => ({ ...prev, isOpen: false }));
          resolve(true);
        },
        onCancel: () => {
          setDialog(prev => ({ ...prev, isOpen: false }));
          resolve(false);
        },
      });
    });
  };

  const closeDialog = () => {
    setDialog(prev => ({ ...prev, isOpen: false }));
    if (resolver.current) {
      resolver.current(false);
    }
  };

  return (
    <DialogContext.Provider value={{ confirm }}>
      {children}
      
      {dialog.isOpen && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
          {/* Backdrop */}
          <div 
            className="absolute inset-0 bg-black/50 backdrop-blur-sm transition-opacity"
            onClick={closeDialog}
          />
          
          {/* Dialog Panel */}
          <div className="relative bg-white rounded-xl shadow-2xl max-w-md w-full overflow-hidden transform transition-all animate-scale-in">
            {/* Header */}
            <div className="px-6 py-4 border-b border-gray-100 flex items-center justify-between bg-gray-50/50">
              <h3 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
                {dialog.type === 'danger' && <AlertTriangle className="w-5 h-5 text-red-500" />}
                {dialog.type === 'confirm' && <HelpCircle className="w-5 h-5 text-blue-500" />}
                {dialog.title}
              </h3>
              <button 
                onClick={closeDialog}
                className="text-gray-400 hover:text-gray-500 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            {/* Body */}
            <div className="px-6 py-6">
              <p className="text-gray-600 leading-relaxed">
                {dialog.message}
              </p>
            </div>
            
            {/* Footer */}
            <div className="px-6 py-4 bg-gray-50 flex items-center justify-end space-x-3">
              <button
                onClick={dialog.onCancel}
                className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-200 transition-colors font-medium"
              >
                {dialog.cancelText}
              </button>
              <button
                onClick={dialog.onConfirm}
                className={`flex items-center space-x-2 px-4 py-2 text-white rounded-lg focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors font-medium ${
                  dialog.type === 'danger' 
                    ? 'bg-red-600 hover:bg-red-700 focus:ring-red-500' 
                    : 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500'
                }`}
              >
                <span>{dialog.confirmText}</span>
              </button>
            </div>
          </div>
        </div>
      )}
    </DialogContext.Provider>
  );
};

export const useDialog = () => {
  const context = useContext(DialogContext);
  if (!context) {
    throw new Error('useDialog must be used within a DialogProvider');
  }
  return context;
};
