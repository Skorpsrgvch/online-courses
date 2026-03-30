import React from 'react';

interface SkeletonProps {
  className?: string;
  width?: string | number;
  height?: string | number;
  borderRadius?: string;
  variant?: 'text' | 'circular' | 'rectangular';
}

export const Skeleton: React.FC<SkeletonProps> = ({
  className = '',
  width,
  height,
  borderRadius = '0.5rem', // default rounded-lg
  variant = 'rectangular',
}) => {
  const baseStyles = "bg-gray-200 animate-pulse";
  
  const variantStyles = {
    rectangular: "w-full",
    circular: "rounded-full",
    text: "h-4 w-full rounded",
  };

  const style: React.CSSProperties = {
    width,
    height,
    borderRadius: variant === 'circular' ? '50%' : borderRadius,
  };

  return (
    <div
      className={`${baseStyles} ${variantStyles[variant]} ${className}`}
      style={style}
      aria-hidden="true"
    />
  );
};