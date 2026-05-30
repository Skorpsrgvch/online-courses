import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { adminService } from '../../api/admin.service';
import type { StudentDetailsResponse, StudentCourseDetail } from '../../api/types';
import { useAuth } from '../../context/AuthContext';

const StudentDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  
  const [data, setData] = useState<StudentDetailsResponse | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user || user.role !== 'admin') {
      navigate('/');
      return;
    }
    if (!id) return;

    const loadData = async () => {
      try {
        const result = await adminService.getStudentDetails(Number(id));
        setData(result);
      } catch (err) {
        console.error('Ошибка загрузки деталей:', err);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [id, navigate, user]);

  const formatDate = (dateStr?: string | null) => {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('ru-RU', {
      day: 'numeric', month: 'long', year: 'numeric'
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-rose-500"></div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-xl font-bold text-gray-800 mb-2">Студент не найден</h2>
          <button onClick={() => navigate('/admin/students')} className="text-rose-600 hover:underline">
            Вернуться к списку
          </button>
        </div>
      </div>
    );
  }

  const { student, courses } = data;

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-5xl mx-auto">
        <header className="mb-8">
          <button 
            onClick={() => navigate('/admin/students')}
            className="mb-4 inline-flex items-center gap-2 text-sm text-gray-500 hover:text-rose-600 font-medium"
          >
            ← Назад к списку
          </button>
          <div className="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
            <h1 className="text-2xl font-serif font-bold text-gray-900">{student.name}</h1>
            <p className="text-gray-500 mt-1">{student.email}</p>
            <div className="mt-4 flex flex-wrap gap-4 text-sm">
              <div className="px-3 py-1 bg-gray-100 rounded-full text-gray-700">
                Регистрация: <span className="font-semibold">{formatDate(student.registered_at)}</span>
              </div>
              <div className="px-3 py-1 bg-rose-50 text-rose-700 rounded-full font-bold">
                Общий прогресс: {student.progress_percent}%
              </div>
            </div>
          </div>
        </header>

        <h2 className="text-xl font-bold text-gray-800 mb-4">Купленные курсы</h2>
        
        <div className="space-y-4">
          {courses.map((course) => (
            <div key={course.course_id} className="bg-white p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow">
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div className="flex-1">
                  <h3 className="text-lg font-bold text-gray-900">{course.course_title}</h3>
                  <p className="text-sm text-gray-500 mt-1">
                    Куплен: {formatDate(course.purchased_at)}
                  </p>
                  {course.last_lesson_title && (
                    <p className="text-xs text-gray-400 mt-2">
                      Последний урок: <span className="text-gray-600">{course.last_lesson_title}</span>
                    </p>
                  )}
                </div>

                <div className="w-full md:w-64">
                  <div className="flex justify-between text-xs mb-1">
                    <span className="font-medium text-gray-600">
                      {course.completed_lessons} / {course.total_lessons} уроков
                    </span>
                    <span className="font-bold text-rose-600">{course.progress_percent}%</span>
                  </div>
                  <div className="w-full h-3 bg-gray-100 rounded-full overflow-hidden">
                    <div 
                      className={`h-full rounded-full transition-all duration-500 ${
                        course.progress_percent === 100 ? 'bg-green-500' : 'bg-rose-500'
                      }`}
                      style={{ width: `${course.progress_percent}%` }}
                    />
                  </div>
                </div>
              </div>
            </div>
          ))}
          
          {courses.length === 0 && (
            <div className="text-center py-12 bg-white rounded-xl border border-dashed border-gray-300">
              <p className="text-gray-500">У этого студента пока нет купленных курсов.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default StudentDetailPage;