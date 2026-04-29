import React from 'react';
import { Outlet } from 'react-router-dom';
import { Header } from './Header';
import { Footer } from './Footer';
import { CookieBanner } from '../ui/CookieBanner';


const MainLayout: React.FC = () => {
  return (
    <div className="min-h-screen flex flex-col bg-white font-sans text-gray-900 selection:bg-rose-200 selection:text-rose-900">
      {/* Хедер фиксированный */}
      <Header />
      
      {/* Основной контент растягивается, прижимая футер к низу */}
      <main className="flex-grow pt-16 bg-linear-to-br from-rose-50 via-white to-lavender-50 relative z-0"> 
       
        
        <Outlet />
      </main>

      {/* Футер */}
      <Footer />
      
      {/* Баннер cookie */}
      <CookieBanner />
    </div>
  );
};

export default MainLayout;