import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { authService } from '../api';
import { User } from '../api/types';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string, agreeToTerms: boolean) => Promise<void>;
  logout: () => Promise<void>;
  error: string | null;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Проверка сессии при загрузке приложения
  useEffect(() => {
    const checkAuth = async () => {
      try {
        const currentUser = await authService.getCurrentUser();
        setUser(currentUser);
      } catch (err) {
        setUser(null);
        // Ошибку 401 игнорируем, это нормальная ситуация для неавторизованного пользователя
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setError(null);
      const userData = await authService.login({ email, password });
      setUser(userData);
    } catch (err: any) {
      setError(err.message || 'Ошибка входа');
      throw err;
    }
  };

  const register = async (name: string, email: string, password: string, agreeToTerms: boolean) => {
    try {
      setError(null);
      const userData = await authService.register({ name, email, password, agreeToTerms });
      setUser(userData);
    } catch (err: any) {
      setError(err.message || 'Ошибка регистрации');
      throw err;
    }
  };

  const logout = async () => {
    try {
      await authService.logout();
    } catch (err) {
      console.error('Ошибка при выходе:', err);
    } finally {
      setUser(null);
    }
  };

  const clearError = () => setError(null);

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
        error,
        clearError,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth должен использоваться внутри AuthProvider');
  }
  return context;
};