import React, { useState } from 'react';
import { Modal } from '../ui/Modal';

// Заглушки изображений (превью видео)
const video1Thumb = "https://images.unsplash.com/photo-1571019614242-c5c5dee9f50b?q=80&w=2070&auto=format&fit=crop";
const video2Thumb = "https://images.unsplash.com/photo-1518611012118-696072aa579a?q=80&w=2070&auto=format&fit=crop";
const video3Thumb = "https://images.unsplash.com/photo-1599447421405-0c1a1d5f8d94?q=80&w=2070&auto=format&fit=crop";

interface VideoItem {
  id: number;
  title: string;
  duration: string;
  thumbnail: string;
  // В реальном проекте здесь будет ID видео с RuTube или ссылка на embed
  embedUrl: string; 
}

const videosData: VideoItem[] = [
  {
    id: 1,
    title: "5 упражнений для расслабления тазового дна",
    duration: "12:45",
    thumbnail: video1Thumb,
    embedUrl: "https://rutube.ru/play/embed/YOUR_VIDEO_ID_1" // Замените на реальный ID
  },
  {
    id: 2,
    title: "Как правильно дышать при болях в тазу",
    duration: "08:20",
    thumbnail: video2Thumb,
    embedUrl: "https://rutube.ru/play/embed/YOUR_VIDEO_ID_2"
  },
  {
    id: 3,
    title: "Мифы о диастазе: разбор эксперта",
    duration: "15:10",
    thumbnail: video3Thumb,
    embedUrl: "https://rutube.ru/play/embed/YOUR_VIDEO_ID_3"
  }
];

export const VideoSection: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentVideoUrl, setCurrentVideoUrl] = useState<string>("");

  const openVideo = (url: string) => {
    setCurrentVideoUrl(url);
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
    setCurrentVideoUrl("");
  };

  return (
    <section id="videos" className="py-12 md:py-16 lg:py-24 relative z-10 ">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Заголовок секции */}
        <div className="text-center mb-12 md:mb-16 max-w-3xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-serif font-bold text-gray-900 mb-4">
            Видео-материалы
          </h2>
          <p className="text-base text-gray-600 leading-relaxed">
            Короткие уроки и разборы техник для самостоятельной работы над здоровьем.
          </p>
        </div>

        {/* Сетка видео */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 lg:gap-10">
          {videosData.map((video) => (
            <div key={video.id} className="group relative cursor-pointer" onClick={() => openVideo(video.embedUrl)}>
              
              {/* Контейнер изображения */}
              <div className="relative aspect-video rounded-2xl overflow-hidden shadow-lg group-hover:shadow-2xl transition-shadow duration-300 bg-gray-900">
                <img 
                  src={video.thumbnail} 
                  alt={video.title} 
                  className="w-full h-full object-cover opacity-90 group-hover:opacity-100 group-hover:scale-105 transition-all duration-500"
                />
                
                {/* Затемнение при наведении */}
                <div className="absolute inset-0 bg-black/40 group-hover:bg-black/50 transition-colors duration-300"></div>

                {/* Кнопка Play по центру */}
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="w-16 h-16 bg-white/90 backdrop-blur-sm rounded-full flex items-center justify-center pl-1 shadow-xl group-hover:scale-110 group-hover:bg-white transition-all duration-300">
                    <svg className="w-6 h-6 text-rose-500 ml-1" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M8 5v14l11-7z" />
                    </svg>
                  </div>
                </div>

                {/* Длительность видео */}
                <div className="absolute bottom-3 right-3 bg-black/70 backdrop-blur-sm px-2 py-1 rounded-md text-xs font-medium text-white">
                  {video.duration}
                </div>
              </div>

              {/* Название под видео */}
              <div className="mt-4 text-center lg:text-left">
                <h3 className="text-lg font-serif font-bold text-gray-900 leading-snug group-hover:text-rose-600 transition-colors">
                  {video.title}
                </h3>
              </div>
            </div>
          ))}
        </div>

        {/* Кнопка "Все видео" */}
        <div className="text-center mt-12">
          <a href="/videos" className="inline-flex items-center gap-2 px-6 py-3 bg-white border border-gray-200 text-gray-700 font-medium rounded-full hover:bg-rose-50 hover:border-rose-300 hover:text-rose-600 hover:shadow-md transition-all duration-300">
            Смотреть все видео
            <svg className="w-4 h-4 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </a>
        </div>

      </div>

      {/* Модальное окно для просмотра видео */}
      <Modal isOpen={isModalOpen} onClose={closeModal} title="">
        <div className="aspect-video w-full bg-black rounded-xl overflow-hidden">
          {currentVideoUrl ? (
            <iframe
              src={currentVideoUrl}
              title="Video Player"
              className="w-full h-full"
              frameBorder="0"
              allow="clipboard-write; encrypted-media; picture-in-picture;"
              allowFullScreen
            ></iframe>
          ) : (
            <div className="flex items-center justify-center h-full text-white">Загрузка...</div>
          )}
        </div>
      </Modal>
    </section>
  );
};