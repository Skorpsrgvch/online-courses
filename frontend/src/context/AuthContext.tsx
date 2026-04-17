import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { authService, getAccessToken } from '../api/auth.service';
import type { AuthResponse, User } from '../api/types';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (fullName: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  error: string | null;
  clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

function authResponseToUser(data: AuthResponse): User {
  return {
    id: data.user_id,
    email: data.email,
    name: data.name,
    role: data.role as 'user' | 'admin',
  };
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // При загрузке проверяем, есть ли у нас валидный токен
  useEffect(() => {
    const checkAuth = async () => {
      const token = getAccessToken();
      if (!token) {
        // Нет токена — пользователь не авторизован
        setIsLoading(false);
        return;
      }
      try {
        const me = await authService.getMe();
        setUser(authResponseToUser(me));
      } catch {
        // Токен невалиден
        authService.clearTokens();
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

  const login = async (email: string, password: string) => {
    try {
      setError(null);
      const data = await authService.login({ email, password });
      setUser(authResponseToUser(data));
    } catch (err: any) {
      setError(err.message || 'Ошибка входа');
      throw err;
    }
  };

  const register = async (fullName: string, email: string, password: string) => {
    try {
      setError(null);
      await authService.register({ full_name: fullName, email, password });
    } catch (err: any) {
      setError(err.message || 'Ошибка регистрации');
      throw err;
    }
  };

  const logout = () => {
    authService.logout();
    setUser(null);
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
