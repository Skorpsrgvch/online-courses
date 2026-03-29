import React, { useState, useEffect } from 'react';
import DOMPurify from '@/lib/dompurify';

interface VideoPlayerProps {
  videoUrl: string;
  onProgressUpdate: (progress: number) => void;
  onComplete: () => void;
}

const VideoPlayer: React.FC<VideoPlayerProps> = ({ 
  videoUrl, 
  onProgressUpdate,
  onComplete 
}) => {
  const [isRutube, setIsRutube] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Проверка домена RuTube для безопасности
    const cleanUrl = DOMPurify.sanitize(videoUrl, { 
      ALLOWED_URI_REGEXP: /^(https?:)?\/\/(www\.)?rutube\.ru/
    });
    
    setIsRutube(cleanUrl.includes('rutube.ru'));
    
    if (!cleanUrl.includes('rutube.ru')) {
      setError('Поддерживаются только видео с RuTube');
      return;
    }
    
    // Эмуляция прогресса просмотра
    const timer = setInterval(() => {
      setProgress(prev => {
        if (prev >= 90) {
          clearInterval(timer);
          onComplete();
          return 90;
        }
        return prev + 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [videoUrl, onComplete]);

  useEffect(() => {
    onProgressUpdate(progress);
  }, [progress, onProgressUpdate]);

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-xl p-4 text-red-700">
        {error}
      </div>
    );
  }

  if (!isRutube) {
    return null;
  }

  return (
    <div className="w-full bg-black rounded-xl overflow-hidden">
      <div className="relative" style={{ paddingTop: '56.25%' }}>
        <iframe
          src={DOMPurify.sanitize(videoUrl, {
            ADD_ATTR: ['allow', 'frameborder'],
            FORBID_ATTR: ['on*'],
            ADD_TAGS: ['iframe']
          })}
          className="absolute top-0 left-0 w-full h-full"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowFullScreen
          title="Видео урок"
        />
      </div>
      <div className="px-4 py-3 bg-gray-800 text-white">
        <div className="flex justify-between items-center">
          <span className="text-sm">Прогресс: {progress}%</span>
          <button 
            onClick={onComplete}
            disabled={progress < 90}
            className={`text-sm font-medium px-3 py-1 rounded ${
              progress >= 90 
                ? 'bg-pink-500 hover:bg-pink-600' 
                : 'bg-gray-600 cursor-not-allowed'
            }`}
          >
            Отметить как пройденное
          </button>
        </div>
        <div className="h-1 bg-gray-700 rounded mt-2 overflow-hidden">
          <div 
            className="h-full bg-pink-500 transition-all duration-300" 
            style={{ width: `${progress}%` }}
          ></div>
        </div>
      </div>
    </div>
  );
};

export default VideoPlayer;