import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { apiClient } from '../../api';

type Step = 'email' | 'code' | 'success';

const PasswordRecoveryPage = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState<Step>('email');
  
  const [email, setEmail] = useState('');
  const [code, setCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  
  const [timeLeft, setTimeLeft] = useState(0);


  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  useEffect(() => {
    let timer: number;
    if (timeLeft > 0) {
      timer = window.setInterval(() => {
        setTimeLeft((prev) => prev - 1);
      }, 1000);
    }
    return () => clearInterval(timer);
  }, [timeLeft]);

  const handleSendCode = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    
    try {
      await apiClient.post('/auth/forgot-password', { email });
      
      setSuccessMessage(`Код отправлен на ${email}`);
      setStep('code');
      setTimeLeft(120);
      setCode('');
    } catch (err: any) {
      setError(err.message || 'Ошибка соединения. Попробуйте позже.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleResendCode = async () => {
    if (timeLeft > 0) return;
    
    setIsLoading(true);
    setError(null);
    try {
      await apiClient.post('/auth/forgot-password', { email });
      setTimeLeft(120);
      setSuccessMessage('Новый код отправлен!');
    } catch (err: any) {
      setError(err.message || 'Ошибка отправки');
    } finally {
      setIsLoading(false);
    }
  };

  // Шаг 2: Сброс пароля по коду
  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (newPassword !== confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }
    if (newPassword.length < 6) {
      setError('Пароль должен быть не менее 6 символов');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await apiClient.post('/auth/reset-password', {
        code,
        new_password: newPassword,
      });
      
      setStep('success');
      setSuccessMessage('Пароль успешно изменен!');
    } catch (err: any) {
      setError(err.message || 'Неверный код или ошибка сервера');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 sm:p-8">
          <h2 className="text-2xl font-serif font-bold text-center text-gray-800 mb-2">
            Восстановление пароля
          </h2>
          <p className="text-center text-gray-500 text-sm mb-6">
            Следуйте инструкциям для смены пароля
          </p>

          {error && (
            <div className="mb-4 p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100 animate-pulse">
              {error}
            </div>
          )}

          {successMessage && step !== 'success' && (
            <div className="mb-4 p-3 bg-green-50 text-green-700 text-sm rounded-lg border border-green-100">
              {successMessage}
            </div>
          )}

          {step === 'email' && (
            <form onSubmit={handleSendCode} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none transition-all"
                  placeholder="example@mail.ru"
                  required
                />
              </div>
              <button
                type="submit"
                disabled={isLoading || !email}
                className="w-full py-3 px-4 bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium shadow-md shadow-rose-200"
              >
                {isLoading ? 'Отправка...' : 'Получить код'}
              </button>
            </form>
          )}

          {step === 'code' && (
            <form onSubmit={handleResetPassword} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Код из письма
                </label>
                <input
                  type="text"
                  inputMode="numeric"
                  pattern="[0-9]*"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 font-mono text-lg tracking-widest text-center outline-none"
                  placeholder="000000"
                  required
                  autoFocus
                />
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Новый пароль</label>
                  <input
                    type="password"
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none"
                    placeholder="••••••"
                    minLength={6}
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Подтверждение</label>
                  <input
                    type="password"
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 outline-none"
                    placeholder="••••••"
                    required
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={isLoading || code.length !== 6 || !newPassword || !confirmPassword}
                className="w-full py-3 px-4 bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium shadow-md shadow-rose-200"
              >
                {isLoading ? 'Сохранение...' : 'Сменить пароль'}
              </button>

              {/* Блок повторной отправки */}
              <div className="pt-4 border-t border-gray-100 text-center">
                {timeLeft > 0 ? (
                  <p className="text-xs text-gray-400">
                    Отправить код повторно через <span className="font-mono font-medium text-gray-600">{timeLeft}</span> сек.
                  </p>
                ) : (
                  <button
                    type="button"
                    onClick={handleResendCode}
                    disabled={isLoading}
                    className="text-xs text-rose-500 hover:text-rose-600 font-medium hover:underline disabled:opacity-50"
                  >
                    Отправить код повторно
                  </button>
                )}
              </div>
            </form>
          )}

          {step === 'success' && (
            <div className="text-center py-6">
              <div className="w-16 h-16 bg-green-100 text-green-500 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h3 className="text-lg font-bold text-gray-800 mb-2">Готово!</h3>
              <p className="text-gray-600 mb-6">Ваш пароль был успешно изменен.</p>
              <Link
                to="/login"
                className="inline-block w-full sm:w-auto px-6 py-3 bg-rose-500 text-white rounded-2xl! hover:bg-rose-600 transition-colors font-medium shadow-md shadow-rose-200"
                style={{ textDecoration: 'none' }}
              >
                Войти с новым паролем
              </Link>
            </div>
          )}

          <p className="text-center text-sm text-gray-500 mt-8 pt-6 border-t border-gray-100">
            Вспомнили пароль?{' '}
            <button
              type="button"
              onClick={() => navigate('/login')}
              className="text-rose-500 hover:underline font-medium"
            >
              Войти
            </button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default PasswordRecoveryPage;