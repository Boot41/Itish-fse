import React from 'react';
import { FileText } from 'lucide-react';

export const EmptyState: React.FC<{
  title: string;
  description: string;
  icon?: React.ReactNode;
  action?: React.ReactNode;
}> = ({ title, description, icon, action }) => (
  <div className="flex flex-col items-center justify-center p-8 text-center">
    <div className="w-16 h-16 bg-[#2a2a2a] rounded-full flex items-center justify-center mb-4">
      {icon || <FileText className="w-8 h-8 text-[#3ECF8E]" />}
    </div>
    <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
    <p className="text-gray-400 mb-6 max-w-sm">{description}</p>
    {action}
  </div>
);