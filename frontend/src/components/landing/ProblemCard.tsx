import React from 'react';
import { motion } from 'framer-motion';

interface ProblemCardProps {
  title: string;
  description: string;
  icon: React.ReactNode;
  index: number;
}

const ProblemCard: React.FC<ProblemCardProps> = ({ 
  title, 
  description, 
  icon, 
  index 
}) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ 
        duration: 0.4, 
        delay: index * 0.1,
        ease: "easeOut"
      }}
      className="bg-white p-6 rounded-2xl shadow-md hover:shadow-lg transition-shadow border border-gray-100"
    >
      <div className="flex items-start mb-4">
        <div className="bg-pink-50 p-3 rounded-xl mr-4">
          {icon}
        </div>
        <h3 className="text-xl font-semibold text-gray-800 flex-1">{title}</h3>
      </div>
      <p className="text-gray-600 leading-relaxed">{description}</p>
    </motion.div>
  );
};

export default ProblemCard;