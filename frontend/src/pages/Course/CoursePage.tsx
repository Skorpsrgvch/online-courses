import { useParams } from 'react-router-dom';

const CoursePage = () => {
  const { id } = useParams<{ id: string }>();

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">Курс ID: {id}</h1>
      <div className="bg-white p-6 rounded-xl shadow-sm">
        <p className="text-gray-600">Здесь будет плеер курса, список уроков и материалы.</p>
      </div>
    </div>
  );
};

export default CoursePage;