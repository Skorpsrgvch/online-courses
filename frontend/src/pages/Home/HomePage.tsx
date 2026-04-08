import React from 'react';

import { HeroSection } from '../../components/landing/HeroSection';

import { AboutSection } from '../../components/landing/AboutSection'; 

import { CourseCatalog } from '../../components/landing/CourseCatalog'; 

import { ArticleSection } from '../../components/landing/ArticleSection'; 

import { VideoSection } from '../../components/landing/VideoSection'; 

import { ReviewSlider } from '../../components/landing/ReviewSlider'; 

const HomePage = () => {
  return (
    <div className="w-full">

      <HeroSection />

      <AboutSection />

      <CourseCatalog />

      <ArticleSection />

      <VideoSection />

      <ReviewSlider />

    </div>
  );
};

export default HomePage;