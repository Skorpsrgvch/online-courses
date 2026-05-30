import { createBrowserRouter, RouterProvider, useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { AuthProvider } from './context/AuthContext';
import { ProtectedRoute } from './components/layout/ProtectedRoute';
import MainLayout from './components/layout/MainLayout';


import HomePage from './pages/Home/HomePage';
import LoginPage from './pages/Auth/LoginPage';
import RegisterPage from './pages/Auth/RegisterPage';
import PasswordRecoveryPage from './pages/Auth/PasswordRecoveryPage';
import DashboardPage from './pages/Dashboard/DashboardPage';
import CoursePage from './pages/Course/CoursePage';
import LearnPage from './pages/Course/LearnPage';
import AllCoursesPage from './pages/Course/AllCoursesPage';

import AdminPage from './pages/Admin/AdminPage';
import CourseFormPage from './pages/Admin/CourseFormPage';
import PaymentSuccessPage from './pages/Course/PaymentSuccessPage';
import ServicesPage from './pages/Service/ServicesPage';
import ServicesFormPage from './pages/Admin/ServiceFormPage';
import StudentsListPage from './pages/Admin/StudentsListPage';
import StudentDetailPage from './pages/Admin/StudentDetailPage';

import PrivacyPolicyPage from './pages/Policy/PrivacyPolicyPage';
import ConsentPage from './pages/Policy/ConsentPage';
import OfferPage from './pages/Policy/OfferPage';


// Компонент для перенаправления на главную (должен быть объявлен ДО использования)
const NavigateToHome = () => {
  const navigate = useNavigate();

  useEffect(() => {
    navigate('/', { replace: true });
  }, [navigate]);

  return null;
};

const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />, // Обертка для всех страниц
    children: [
      { index: true, element: <HomePage /> },
      { path: 'login', element: <LoginPage /> },
      { path: 'register', element: <RegisterPage /> },
      { path: 'password-recovery', element: <PasswordRecoveryPage /> },
      { path: 'courses', element: <AllCoursesPage /> },
      { path: 'services', element: <ServicesPage /> },
      { path: "privacy-policy", element: <PrivacyPolicyPage /> },
      { path: "privacy-policy", element: <ConsentPage /> },
      { path: "privacy-policy", element: <OfferPage /> },
      
      {
        path: 'dashboard',
        element: (
          <ProtectedRoute>
            <DashboardPage />
          </ProtectedRoute>
        ),
      },
      {
        path: 'course/:id',
        element: (
            <CoursePage />
        ),
      },
      {
        path: 'course/:id/learn',
        element: <ProtectedRoute>
          <LearnPage />
          </ProtectedRoute>,
      },
      {
        path: 'admin',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
            <AdminPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/admin/courses/new',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
            <CourseFormPage mode="create" />
          </ProtectedRoute>
        ),
      },
      {
        path: '/admin/courses/:id/edit',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
            <CourseFormPage mode="edit" />
          </ProtectedRoute>
        ),
      },
      {
        path: '/admin/services/new',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
            <ServicesFormPage mode="create" />
          </ProtectedRoute>
        ),
      },
      {
        path: '/admin/services/:id/edit',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
            <ServicesFormPage mode="edit" />
          </ProtectedRoute>
        ),
      },
      {
        path: '/admin/students',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
            <StudentsListPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/admin/students/:id',
        element: (
          <ProtectedRoute allowedRoles={['admin']}>
           <StudentDetailPage />
          </ProtectedRoute>
        ),
      },
      {
        path: '/course/:id/payment-success',
        element: (
          <ProtectedRoute>
            <PaymentSuccessPage />
          </ProtectedRoute>
        ),
      },

      
      { path: '*', element: <NavigateToHome /> },
    ],
  },
]);

function App() {
  return (
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  );
}

export default App;

