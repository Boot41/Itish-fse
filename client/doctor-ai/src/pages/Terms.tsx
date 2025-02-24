import React from 'react';
import { Scale, FileText, AlertCircle, UserCheck, Zap } from 'lucide-react';

const TermsSection: React.FC<{ 
  icon: React.ReactNode; 
  title: string; 
  description: string; 
  items: string[] 
}> = ({ icon, title, description, items }) => (
  <div className="group glass-card p-8 hover:border-[#3ECF8E] transition-all duration-300 ease-in-out transform hover:scale-[1.01] hover:shadow-lg hover:shadow-[#3ECF8E]/20 overflow-hidden will-change-transform transform-gpu">
    <div className="text-[#3ECF8E] mb-6 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
    <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 mb-6 group-hover:text-gray-300 transition-colors duration-300">{description}</p>
    <ul className="space-y-3">
      {items.map((item, index) => (
        <li key={index} className="flex items-center space-x-2 group-hover:translate-x-1 transition-transform duration-300">
          <Zap className="w-5 h-5 text-[#3ECF8E]" />
          <span className="text-gray-300">{item}</span>
        </li>
      ))}
    </ul>
  </div>
);

const Terms: React.FC = () => {
  const sections = [
    {
      icon: <Scale className="w-12 h-12" />,
      title: 'Terms of Use',
      description: 'General terms and conditions for using Doctor AI.',
      items: [
        'Acceptance of terms',
        'User registration requirements',
        'Account responsibilities',
        'Platform usage guidelines'
      ]
    },
    {
      icon: <FileText className="w-12 h-12" />,
      title: 'Service Agreement',
      description: 'Details about our service delivery and commitments.',
      items: [
        'Service availability',
        'Quality standards',
        'Support services',
        'Updates and maintenance'
      ]
    },
    {
      icon: <AlertCircle className="w-12 h-12" />,
      title: 'Limitations',
      description: 'Understanding service limitations and user restrictions.',
      items: [
        'Usage restrictions',
        'Liability limitations',
        'Service modifications',
        'Termination conditions'
      ]
    },
    {
      icon: <UserCheck className="w-12 h-12" />,
      title: 'User Obligations',
      description: 'Your responsibilities while using our platform.',
      items: [
        'Account security',
        'Content guidelines',
        'Professional conduct',
        'Compliance requirements'
      ]
    }
  ];

  return (
    <div className="min-h-screen bg-[#1C1C1C]">
      <div className="container mx-auto px-4 pt-24 pb-20">
        <div className="text-center max-w-4xl mx-auto mb-16">
          <div 
            className="inline-flex items-center bg-[#2a2a2a] border border-[#3a3a3a] rounded-full px-6 py-2 text-sm mb-8 text-white transform hover:scale-[1.01] transition-all duration-300 cursor-pointer will-change-transform transform-gpu"
          >
            <span className="bg-gradient-to-r from-[#3ECF8E] via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text font-semibold">
              Terms • Conditions
            </span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Terms of Service
          </h1>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            Understanding your rights and responsibilities while using Doctor AI
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 animate-fade-in-up">
          {sections.map((section, index) => (
            <div key={index} className="transform transition-all duration-500" style={{ animationDelay: `${index * 150}ms` }}>
              <TermsSection 
                icon={section.icon} 
                title={section.title} 
                description={section.description} 
                items={section.items} 
              />
            </div>
          ))}
        </div>

        <div className="mt-16 text-center">
          <p className="text-gray-400 max-w-3xl mx-auto leading-relaxed">
            Last updated: {new Date().toLocaleDateString('en-US', { 
              year: 'numeric', 
              month: 'long', 
              day: 'numeric' 
            })}
          </p>
          <p className="text-gray-500 mt-4 text-sm">
            By using Doctor AI, you agree to these terms and conditions.
          </p>
        </div>
      </div>
    </div>
  );
};

export default Terms;
