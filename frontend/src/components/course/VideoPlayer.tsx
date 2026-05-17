import React, { useRef, useEffect, useState } from 'react';

interface VideoPlayerProps {
  url: string;
  onProgress?: (percent: number) => void;
}

export const VideoPlayer: React.FC<VideoPlayerProps> = ({ url, onProgress }) => {
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);

  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    document.addEventListener('webkitfullscreenchange', handleFullscreenChange);
    document.addEventListener('mozfullscreenchange', handleFullscreenChange);
    document.addEventListener('MSFullscreenChange', handleFullscreenChange);

    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
      document.removeEventListener('webkitfullscreenchange', handleFullscreenChange);
      document.removeEventListener('mozfullscreenchange', handleFullscreenChange);
      document.removeEventListener('MSFullscreenChange', handleFullscreenChange);
    };
  }, []);

  const toggleFullscreen = async () => {
    if (!iframeRef.current) return;

    try {
      if (!document.fullscreenElement) {
        // Для мобильных устройств пробуем разные методы
        if (iframeRef.current.requestFullscreen) {
          await iframeRef.current.requestFullscreen();
        } else if ((iframeRef.current as any).webkitRequestFullscreen) {
          await (iframeRef.current as any).webkitRequestFullscreen();
        } else if ((iframeRef.current as any).mozRequestFullScreen) {
          await (iframeRef.current as any).mozRequestFullScreen();
        } else if ((iframeRef.current as any).msRequestFullscreen) {
          await (iframeRef.current as any).msRequestFullscreen();
        }
      } else {
        if (document.exitFullscreen) {
          await document.exitFullscreen();
        } else if ((document as any).webkitExitFullscreen) {
          await (document as any).webkitExitFullscreen();
        } else if ((document as any).mozCancelFullScreen) {
          await (document as any).mozCancelFullScreen();
        } else if ((document as any).msExitFullscreen) {
          await (document as any).msExitFullscreen();
        }
      }
    } catch (err) {
      console.error('Ошибка переключения полноэкранного режима:', err);
    }
  };

  useEffect(() => {
    const checkProgress = setInterval(() => {
      if (onProgress) onProgress(10); 
    }, 5000);

    return () => clearInterval(checkProgress);
  }, [onProgress]);

  return (
    <div className="aspect-video bg-black rounded-xl overflow-hidden shadow-lg relative w-full group">
      <iframe
        ref={iframeRef}
        src={url}
        title="Video Lesson"
        className="w-full h-full absolute inset-0"
        frameBorder="0"
        allow="autoplay; fullscreen; picture-in-picture; encrypted-media; clipboard-write"
        allowFullScreen
        loading="eager"
      ></iframe>
      
      {/* Кнопка полноэкранного режима поверх видео */}
      <button
        onClick={toggleFullscreen}
        className="absolute bottom-4 right-4 z-20 p-3 bg-black/60 hover:bg-black/80 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity duration-300 backdrop-blur-sm"
        title={isFullscreen ? 'Выйти из полноэкранного режима' : 'На весь экран'}
      >
        {isFullscreen ? (
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25" />
          </svg>
        ) : (
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
          </svg>
        )}
      </button>
    </div>
  );
};