import React from 'react';

interface CourseProgressProps {
  completedLessons: number;
  totalLessons: number;
  courseTitle: string;
}

/**
 * Компонент отображения прогресса прохождения курса.
 */
export const CourseProgress: React.FC<CourseProgressProps> = ({
  completedLessons,
  totalLessons,
  courseTitle,
}) => {
  const percentage = totalLessons > 0 ? Math.round((completedLessons / totalLessons) * 100) : 0;
  const isComplete = completedLessons === totalLessons;

  return (
    <div className="bg-white p-5 rounded-xl shadow-sm border border-gray-100">
      <div className="flex items-center justify-between mb-3">
        <h3 className="font-semibold text-gray-800 truncate">{courseTitle}</h3>
        {isComplete && (
          <span className="text-xs px-2 py-1 bg-green-100 text-green-700 rounded-full font-medium">
            ✓ Завершён
          </span>
        )}
      </div>

      <div className="flex items-center gap-3">
        {/* Прогресс-бар */}
        <div className="flex-1 bg-gray-200 rounded-full h-3 overflow-hidden">
          <div
            className={`h-full rounded-full transition-all duration-500 ${
              isComplete ? 'bg-green-500' : 'bg-rose-500'
            }`}
            style={{ width: `${percentage}%` }}
          />
        </div>
        <span className="text-sm font-medium text-gray-600 min-w-[48px] text-right">
          {percentage}%
        </span>
      </div>

      <p className="text-xs text-gray-500 mt-2">
        {completedLessons} из {totalLessons} уроков
      </p>
    </div>
  );
};
