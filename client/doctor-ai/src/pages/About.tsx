import React from 'react';
import { CloudLightning, ShieldCheck, FileText, BarChart2 } from 'lucide-react';

const FeatureCard: React.FC<{ icon: React.ReactNode; title: string; description: string }> = ({ icon, title, description }) => (
  <div className="group glass-card p-8 hover:border-[#3ECF8E] transition-all duration-300 ease-in-out transform hover:scale-[1.01] hover:shadow-lg hover:shadow-[#3ECF8E]/20 overflow-hidden will-change-transform transform-gpu">
    <div className="text-[#3ECF8E] mb-6 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
    <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 group-hover:text-gray-300 transition-colors duration-300 mb-6">{description}</p>
    <div className="mt-4 h-1 w-0 group-hover:w-full bg-gradient-to-r from-[#3ECF8E] to-[#4BB563] transition-all duration-300"></div>
  </div>
);

const About: React.FC = () => {
  const features = [
    {
      icon: <CloudLightning className="w-12 h-12" />,
      title: 'AI-Powered Voice Transcription',
      description: 'Secure real-time transcription of doctor-patient conversations using advanced NLP and cloud-based LLM APIs.'
    },
    {
      icon: <FileText className="w-12 h-12" />,
      title: 'Automated Report Generation',
      description: 'Transform raw transcriptions into comprehensive, structured medical reports with intelligent insights.'
    },
    {
      icon: <ShieldCheck className="w-12 h-12" />,
      title: 'Secure Authentication',
      description: 'Robust identity verification and data protection ensuring compliance with HIPAA/GDPR regulations.'
    },
    {
      icon: <BarChart2 className="w-12 h-12" />,
      title: 'Doctor Dashboard',
      description: 'Track consultation trends, access past reports, and gain insights into your medical practice.'
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
              About Us • Innovation in Healthcare
            </span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            AI-Powered Doctor Assistant
          </h1>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            Revolutionizing healthcare documentation through intelligent AI-driven solutions
          </p>
        </div>

        <section className="mb-24">
          <h2 className="text-4xl md:text-5xl font-bold text-center mb-16 bg-gradient-to-r from-[#3ECF8E] via-white to-[#4BB563] text-transparent bg-clip-text opacity-90 transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">Key Features</h2>
          <div className="grid md:grid-cols-2 gap-8 animate-fade-in-up">
            {features.map((feature, index) => (
              <div key={index} className="transform transition-all duration-500" style={{ animationDelay: `${index * 150}ms` }}>
                <FeatureCard 
                  icon={feature.icon} 
                  title={feature.title} 
                  description={feature.description} 
                />
              </div>
            ))}
          </div>
        </section>

        <section className="mission-section mt-32 text-center bg-[#2a2a2a] bg-opacity-30 border border-[#3a3a3a] rounded-xl p-12 transform hover:scale-[1.01] hover:shadow-2xl hover:shadow-[#3ECF8E]/20 transition-all duration-500 will-change-transform transform-gpu">
          <h2 className="text-4xl md:text-5xl font-bold mb-8 bg-gradient-to-r from-[#3ECF8E] via-[#4BB563] to-white text-transparent bg-clip-text animate-gradient-x transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Our Mission
          </h2>
          <p className="text-xl text-gray-400 max-w-3xl mx-auto leading-relaxed mb-12 transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            To empower healthcare professionals by providing cutting-edge AI tools that streamline documentation, 
            enhance patient interactions, and ultimately improve the quality of medical care.
          </p>
          <div className="inline-flex items-center bg-[#2a2a2a] border border-[#3a3a3a] rounded-full px-8 py-3 text-white transform hover:scale-[1.01] transition-all duration-300 cursor-pointer will-change-transform transform-gpu">
            <div className="absolute inset-0 bg-gradient-to-r from-[#3ECF8E]/20 to-[#4BB563]/20 transform scale-x-0 group-hover/button:scale-x-100 transition-transform duration-500 origin-left"></div>
            <span className="bg-gradient-to-r from-[#3ECF8E] via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text font-semibold relative">
              Learn More About Our Vision
            </span>
          </div>
        </section>
      </div>
    </div>
  );
};

export default About;
