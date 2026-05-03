import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { coursesService } from '../../api/courses.service';

export const PaymentSuccessPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [isVerifying, setIsVerifying] = useState(true);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
  const timer = setTimeout(() => {
    // Вместо navigate используем window.location.href для полной перезагрузки страницы
    // Это гарантирует, что все интерцепторы и стейт авторизации инициализируются заново
    if (id) {
      window.location.href = `/course/${id}`; 
    } else {
      window.location.href = '/dashboard';
    }
  }, 2000); // Даем 2 секунды на показ сообщения

  return () => clearTimeout(timer);
}, [id]);

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="text-center bg-white p-8 rounded-xl shadow-lg max-w-md w-full">
        {isVerifying ? (
          <>
            <div className="w-16 h-16 border-4 border-rose-300 border-t-rose-500 rounded-full animate-spin mx-auto mb-6"></div>
            <h2 className="text-2xl font-bold text-gray-800 mb-2">Проверка платежа...</h2>
            <p className="text-gray-600">Пожалуйста, не закрывайте страницу.</p>
          </>
        ) : success ? (
          <>
            <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <svg className="w-10 h-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h2 className="text-2xl font-bold text-gray-800 mb-2">Оплата успешна!</h2>
            <p className="text-gray-600 mb-6">Доступ к курсу открыт. Перенаправляем к обучению...</p>
            <button
              onClick={() => id && navigate(`/course/${id}/learn`)}
              className="px-6 py-3 bg-rose-500 text-white rounded-xl hover:bg-rose-600 transition-colors font-medium"
            >
              Перейти к урокам
            </button>
          </>
        ) : (
          <div className="text-red-500">Ошибка проверки</div>
        )}
      </div>
    </div>
  );
};

export default PaymentSuccessPage;