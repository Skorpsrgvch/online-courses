import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ProtectedRoute } from './components/layout/ProtectedRoute';

// Pages
import HomePage from './pages/Home';
import LoginPage from './pages/Auth/Login';
import RegisterPage from './pages/Auth/Register';
import DashboardPage from './pages/Dashboard';
import CoursePage from './pages/Course';
import AdminPage from './pages/Admin';
import PasswordRecoveryPage from './pages/Auth/PasswordRecovery';

const router = createBrowserRouter([
  {
    path: '/',
    element: <HomePage />,
  },
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/password-recovery',
    element: <PasswordRecoveryPage />,
  },
  {
    path: '/dashboard',
    element: (
      <ProtectedRoute>
        <DashboardPage />
      </ProtectedRoute>
    ),
  },
  {
    path: '/course/:id',
    element: (
      <ProtectedRoute>
        <CoursePage />
      </ProtectedRoute>
    ),
  },
  {
    path: '/admin',
    element: (
      <ProtectedRoute allowedRoles={['admin']}>
        <AdminPage />
      </ProtectedRoute>
    ),
  },
  // Обработка несуществующих маршрутов
  {
    path: '*',
    element: <NavigateToHome />,
  }
]);

function App() {
  return (
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  );
}

// Компонент для перенаправления на главную
const NavigateToHome = () => {
  const navigate = useNavigate();
  
  useEffect(() => {
    navigate('/', { replace: true });
  }, [navigate]);
  
  return null;
};

export default App;