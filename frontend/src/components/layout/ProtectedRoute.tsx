import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/context/AuthContext';
import { Skeleton } from '@/components/ui/Skeleton';

interface ProtectedRouteProps {
  children: JSX.Element;
  allowedRoles?: ('user' | 'admin')[];
  redirectPath?: string;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ 
  children, 
  allowedRoles = ['user'],
  redirectPath = '/login'
}) => {
  const { user, loading, checkAuth } = useAuth();
  const location = useLocation();

  // Проверка аутентификации при монтировании
  useEffect(() => {
    checkAuth();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
          <Skeleton className="h-64 w-full rounded-xl" />
        </div>
        <Footer />
      </div>
    );
  }

  if (!user) {
    return <Navigate to={redirectPath} state={{ from: location }} replace />;
  }

  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return <Navigate to="/" replace />;
  }

  return children;
};