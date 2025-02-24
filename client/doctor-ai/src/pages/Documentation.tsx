import React from 'react';
import { Book, Rocket, Shield, Settings, Zap } from 'lucide-react';

const DocSection: React.FC<{ 
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

const Documentation: React.FC = () => {
  const sections = [
    {
      icon: <Rocket className="w-12 h-12" />,
      title: 'Getting Started',
      description: 'Everything you need to know to begin using Doctor AI effectively.',
      items: [
        'Create and configure your account',
        'Set up your professional profile',
        'Configure notification preferences',
        'Start your first AI-assisted consultation'
      ]
    },
    {
      icon: <Book className="w-12 h-12" />,
      title: 'Platform Features',
      description: 'Comprehensive guide to Doctor AI\'s powerful capabilities.',
      items: [
        'AI-powered consultation assistance',
        'Voice transcription technology',
        'Patient history management',
        'Analytics and reporting tools'
      ]
    },
    {
      icon: <Shield className="w-12 h-12" />,
      title: 'Security & Compliance',
      description: 'Understanding our robust security measures and compliance standards.',
      items: [
        'HIPAA compliance guidelines',
        'Data encryption protocols',
        'Patient data protection',
        'Access control management'
      ]
    },
    {
      icon: <Settings className="w-12 h-12" />,
      title: 'Best Practices',
      description: 'Optimize your workflow with recommended usage patterns.',
      items: [
        'Effective consultation recording',
        'Patient data organization',
        'Report generation tips',
        'Performance optimization'
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
              Documentation • User Guide
            </span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Platform Documentation
          </h1>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            Your comprehensive guide to mastering Doctor AI's features and capabilities
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 animate-fade-in-up">
          {sections.map((section, index) => (
            <div key={index} className="transform transition-all duration-500" style={{ animationDelay: `${index * 150}ms` }}>
              <DocSection 
                icon={section.icon} 
                title={section.title} 
                description={section.description} 
                items={section.items} 
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Documentation;
