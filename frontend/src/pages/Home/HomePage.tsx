import { HeroSection } from '../../components/landing/HeroSection';

import { AboutSection } from '../../components/landing/AboutSection'; 

import { ServicesSection } from '../../components/landing/ServicesSection';

import CourseCatalog  from '../../components/landing/CourseCatalog'; 

import { ArticleSection } from '../../components/landing/ArticleSection'; 

import { VideoSection } from '../../components/landing/VideoSection'; 



const HomePage = () => {
  return (
    <div className="w-full">

      <HeroSection />

      <AboutSection />

      <ServicesSection />

      <CourseCatalog />

      <ArticleSection />

      <VideoSection />

    </div>
  );
};

export default HomePage;