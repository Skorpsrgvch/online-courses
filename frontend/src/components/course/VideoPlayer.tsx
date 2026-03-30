import React, { useRef, useEffect } from 'react';

interface VideoPlayerProps {
  url: string; // Ссылка на RuTube или другой источник
  onProgress?: (percent: number) => void;
}

export const VideoPlayer: React.FC<VideoPlayerProps> = ({ url, onProgress }) => {
  const iframeRef = useRef<HTMLIFrameElement>(null);

  // Простая эмуляция отслеживания просмотра для RuTube (в реальности нужен postMessage API провайдера)
  useEffect(() => {
    const checkProgress = setInterval(() => {
      // Логика проверки прогресса зависит от конкретного плеера
      // Для примера отправим 90% через 10 секунд, если видео длинное
      if (onProgress) onProgress(10); 
    }, 5000);

    return () => clearInterval(checkProgress);
  }, [onProgress]);

  return (
    <div className="aspect-video bg-black rounded-xl overflow-hidden shadow-lg relative">
      <iframe
        ref={iframeRef}
        src={url}
        title="Video Lesson"
        className="w-full h-full"
        frameBorder="0"
        allow="clipboard-write; encrypted-media; picture-in-picture;"
        allowFullScreen
      ></iframe>
    </div>
  );
};