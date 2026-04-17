import { useState } from 'react';
import { Link } from 'react-router-dom';
import { apiClient } from '../../api';

type Step = 'email' | 'token';

const PasswordRecoveryPage = () => {
  const [step, setStep] = useState<Step>('email');
  const [email, setEmail] = useState('');
  const [token, setToken] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSendReset = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.post('/auth/forgot-password', { email });
      setSuccessMessage(response.data.message);
      // Если токен вернулся (тестовый режим) — показываем шаг сброса
      if (response.data.token) {
        setToken(response.data.token);
        setStep('token');
      } else {
        setStep('token');
      }
    } catch (err: any) {
      setError(err.message || 'Ошибка отправки запроса');
    } finally {
      setIsLoading(false);
    }
  };

  const handleResetPassword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }
    if (newPassword.length < 6) {
      setError('Минимум 6 символов');
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      await apiClient.post('/auth/reset-password', {
        token,
        new_password: newPassword,
      });
      setSuccessMessage('Пароль успешно сброшен! Теперь вы можете войти.');
      setStep('email');
      setToken('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (err: any) {
      setError(err.message || 'Ошибка сброса пароля');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="max-w-md w-full">
        <Link to="/login" className="text-sm text-gray-500 hover:text-rose-500 transition-colors mb-4 inline-block">
          ← Назад ко входу
        </Link>

        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <h2 className="text-2xl font-serif font-bold text-center text-gray-800 mb-2">
            Восстановление пароля
          </h2>

          {error && (
            <div className="mb-4 p-3 bg-red-50 text-red-700 text-sm rounded-lg border border-red-100">
              {error}
            </div>
          )}

          {successMessage && step === 'token' && !error && (
            <div className="mb-4 p-3 bg-green-50 text-green-700 text-sm rounded-lg border border-green-100">
              {successMessage}
            </div>
          )}

          {successMessage && step === 'email' && (
            <div className="mb-4 p-3 bg-green-50 text-green-700 text-sm rounded-lg border border-green-100">
              {successMessage}
            </div>
          )}

          {/* Шаг 1: Ввод email */}
          {step === 'email' && (
            <form onSubmit={handleSendReset} className="space-y-4">
              <p className="text-sm text-gray-600 text-center mb-4">
                Введите email, привязанный к вашему аккаунту.
              </p>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400"
                  placeholder="example@mail.ru"
                  required
                />
              </div>
              <button
                type="submit"
                disabled={isLoading || !email}
                className="w-full py-2.5 px-4 bg-rose-500 text-white rounded-lg hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
              >
                {isLoading ? 'Отправка...' : 'Отправить инструкцию'}
              </button>
            </form>
          )}

          {/* Шаг 2: Ввод токена и нового пароля */}
          {step === 'token' && (
            <form onSubmit={handleResetPassword} className="space-y-4 mt-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Токен из письма
                </label>
                <input
                  type="text"
                  value={token}
                  onChange={(e) => setToken(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400 font-mono text-sm"
                  placeholder="Вставьте токен..."
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Новый пароль</label>
                <input
                  type="password"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400"
                  placeholder="Минимум 6 символов"
                  minLength={6}
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Подтвердите пароль</label>
                <input
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-rose-300 focus:border-rose-400"
                  placeholder="Повторите пароль"
                  required
                />
              </div>
              <button
                type="submit"
                disabled={isLoading || !token || !newPassword || !confirmPassword}
                className="w-full py-2.5 px-4 bg-rose-500 text-white rounded-lg hover:bg-rose-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium"
              >
                {isLoading ? 'Сброс...' : 'Сбросить пароль'}
              </button>
            </form>
          )}

          <p className="text-center text-sm text-gray-500 mt-6">
            Вспомнили пароль?{' '}
            <Link to="/login" className="text-rose-500 hover:underline font-medium">
              Войти
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};

export default PasswordRecoveryPage;
