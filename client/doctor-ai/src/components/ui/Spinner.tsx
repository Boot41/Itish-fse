import React from 'react';

export const Spinner: React.FC<{ className?: string }> = ({ className = '' }) => (
  <div className={`relative ${className}`}>
    <div className="w-10 h-10 border-4 border-[#3ECF8E]/20 border-t-[#3ECF8E] rounded-full animate-spin"></div>
  </div>
);