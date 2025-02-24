import React from 'react';

interface DocumentProps {
  title: string;
  children: React.ReactNode;
  sections?: Array<{ title: string; content: React.ReactNode }>;
}

const Document: React.FC<DocumentProps> = ({ title, children, sections }) => {
  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">{title}</h1>
      
      <div className="text-gray-700 leading-relaxed mb-8">
        {children}
      </div>

      {sections && sections.map((section, index) => (
        <div key={index} className="mb-8">
          <h2 className="text-2xl font-semibold text-gray-800 mb-4">
            {section.title}
          </h2>
          <div className="text-gray-700 leading-relaxed">
            {section.content}
          </div>
        </div>
      ))}
    </div>
  );
};

export default Document;
