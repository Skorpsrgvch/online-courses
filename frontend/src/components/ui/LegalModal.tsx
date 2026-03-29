import React from 'react';
import DOMPurify from '@/lib/dompurify';

interface LegalModalProps {
  title: string;
  content: string;
  isOpen: boolean;
  onClose: () => void;
}

const LegalModal: React.FC<LegalModalProps> = ({ 
  title, 
  content, 
  isOpen, 
  onClose 
}) => {
  if (!isOpen) return null;

  // Очистка HTML контента от XSS угроз
  const cleanContent = DOMPurify.sanitize(content, {
    ALLOWED_TAGS: [
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'p', 'a', 'ul', 'ol', 'li', 'strong', 'em', 
      'br', 'hr', 'blockquote', 'table', 'thead', 'tbody', 'tr', 'th', 'td'
    ],
    ALLOWED_ATTR: ['href', 'title', 'target'],
    FORBID_ATTR: ['style', 'on*']
  });

  return (
    <div 
      className="fixed inset-0 z-50 overflow-y-auto" 
      aria-labelledby="modal-title" 
      role="dialog" 
      aria-modal="true"
    >
      <div className="flex items-end justify-center min-h-screen px-4 pb-20 pt-4 text-center sm:block sm:p-0">
        <div 
          className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
          aria-hidden="true"
          onClick={onClose}
        ></div>
        
        <div className="inline-block align-bottom bg-white rounded-2xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
          <div className="bg-white px-6 pt-6 pb-8">
            <div className="flex justify-between items-start">
              <h2 id="modal-title" className="text-2xl font-bold text-gray-800">{title}</h2>
              <button 
                onClick={onClose}
                className="text-gray-400 hover:text-gray-500"
                aria-label="Закрыть"
              >
                <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div className="mt-6 prose max-w-none text-gray-600">
              <div dangerouslySetInnerHTML={{ __html: cleanContent }} />
            </div>
            
            <div className="mt-8 flex justify-end">
              <button 
                onClick={onClose}
                className="px-6 py-2 bg-pink-500 text-white rounded-lg hover:bg-pink-600 transition-colors"
              >
                Понятно
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LegalModal;