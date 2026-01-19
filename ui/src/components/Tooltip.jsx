import { useState } from 'react';

export default function Tooltip({ content, children }) {
  const [isVisible, setIsVisible] = useState(false);

  return (
    <div 
      className="relative flex items-center"
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      {children}
      {isVisible && (
        <div className="absolute top-full mt-2 left-1/2 transform -translate-x-1/2 z-50 px-3 py-2 bg-white text-gray-700 text-xs font-semibold rounded-lg shadow-xl border border-gray-100 whitespace-nowrap animate-fade-in-down">
          {content}
          {/* Arrow */}
          <div className="absolute -top-1 left-1/2 transform -translate-x-1/2 w-2 h-2 bg-white rotate-45 border-t border-l border-gray-100"></div>
        </div>
      )}
    </div>
  );
}
