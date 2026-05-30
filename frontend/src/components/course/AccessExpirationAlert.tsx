import React from 'react';

interface AccessExpirationAlertProps {
  isExpired: boolean;
  daysRemaining?: number;
  expiresAt?: string;
  onRenew?: () => void;
}

export const AccessExpirationAlert: React.FC<AccessExpirationAlertProps> = ({
  isExpired,
  daysRemaining,
  expiresAt,
  onRenew,
}) => {
  if (isExpired) {
    return (
      <div className="mb-6 p-4 bg-red-50 border-l-4 border-red-500 rounded-r-lg shadow-sm">
        <div className="flex items-start">
          <div className="flex-shrink-0">
            <svg className="h-5 w-5 text-red-500" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="ml-3 w-full">
            <h3 className="text-sm font-bold text-red-800">Доступ к курсу истёк</h3>
            <p className="mt-1 text-sm text-red-700">
              Срок доступа к материалам курса завершился {expiresAt ? new Date(expiresAt).toLocaleDateString('ru-RU') : ''}.
            </p>
            {onRenew && (
              <div className="mt-3">
                <button
                  onClick={onRenew}
                  className="px-4 py-2 bg-red-600 text-white text-xs font-medium rounded-2xl! hover:bg-red-700 transition-colors shadow-sm"
                >
                  Продлить доступ
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  if (daysRemaining !== undefined && daysRemaining <= 14) {
    let colorClass = 'bg-blue-50 border-blue-500 text-blue-800';
    if (daysRemaining <= 3) {
      colorClass = 'bg-red-50 border-red-500 text-red-800';
    } else if (daysRemaining <= 7) {
      colorClass = 'bg-orange-50 border-orange-500 text-orange-800';
    }

    const dayWord = daysRemaining === 1 ? 'день' : (daysRemaining >= 2 && daysRemaining <= 4 ? 'дня' : 'дней');

    return (
      <div className={`mb-6 p-4 border-l-4 rounded-r-lg shadow-sm ${colorClass.replace('text-', 'text-').split(' ')[2] || 'text-gray-800'}`} 
           style={{ backgroundColor: colorClass.includes('red') ? '#fef2f2' : colorClass.includes('orange') ? '#fff7ed' : '#eff6ff', borderColor: colorClass.includes('red') ? '#ef4444' : colorClass.includes('orange') ? '#f97316' : '#3b82f6' }}>
        <div className="flex items-start">
          <div className="flex-shrink-0">
            <svg className={`h-5 w-5 ${colorClass.includes('red') ? 'text-red-500' : colorClass.includes('orange') ? 'text-orange-500' : 'text-blue-500'}`} viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="ml-3">
            <p className={`text-sm font-medium ${colorClass.includes('red') ? 'text-red-800' : colorClass.includes('orange') ? 'text-orange-800' : 'text-blue-800'}`}>
              До окончания доступа осталось <span className="font-bold">{daysRemaining} {dayWord}</span>.
              {expiresAt && <span className="block text-xs mt-1 opacity-90">Истекает: {new Date(expiresAt).toLocaleDateString('ru-RU')}</span>}
            </p>
          </div>
        </div>
      </div>
    );
  }

  return null;
};