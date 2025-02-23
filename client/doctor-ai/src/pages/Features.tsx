import React from 'react';
import { 
  Mic, 
  FileText, 
  ShieldCheck, 
  BarChart2, 
  Zap 
} from 'lucide-react';

const FeatureHighlight: React.FC<{ 
  icon: React.ReactNode; 
  title: string; 
  description: string; 
  details: string[] 
}> = ({ icon, title, description, details }) => (
  <div className="group glass-card p-8 hover:border-[#3ECF8E] transition-all duration-300 ease-in-out transform hover:scale-[1.01] hover:shadow-lg hover:shadow-[#3ECF8E]/20 overflow-hidden will-change-transform transform-gpu">
    <div className="text-[#3ECF8E] mb-6 transform group-hover:scale-110 transition-transform duration-300">{icon}</div>
    <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-[#3ECF8E] transition-colors duration-300">{title}</h3>
    <p className="text-gray-400 mb-6 group-hover:text-gray-300 transition-colors duration-300">{description}</p>
    <ul className="space-y-3">
      {details.map((detail, index) => (
        <li key={index} className="flex items-center space-x-2 group-hover:translate-x-1 transition-transform duration-300">
          <Zap className="w-5 h-5 text-[#3ECF8E]" />
          <span className="text-gray-300">{detail}</span>
        </li>
      ))}
    </ul>
  </div>
);

const Features: React.FC = () => {
  const featuresList = [
    {
      icon: <Mic className="w-12 h-12" />,
      title: 'AI-Powered Voice Transcription',
      description: 'Advanced real-time conversation transcription with natural language processing.',
      details: [
        'Accurate medical terminology recognition',
        'Real-time audio-to-text conversion',
        'Noise reduction and clarity enhancement',
        'Support for multiple medical specialties'
      ]
    },
    {
      icon: <FileText className="w-12 h-12" />,
      title: 'Intelligent Report Generation',
      description: 'Automated medical report creation with AI-driven insights and structuring.',
      details: [
        'Comprehensive consultation summary',
        'Automatic medical terminology standardization',
        'Customizable report templates',
        'Easy editing and refinement tools'
      ]
    },
    {
      icon: <ShieldCheck className="w-12 h-12" />,
      title: 'Secure Authentication & Compliance',
      description: 'Robust security measures ensuring patient data protection and regulatory compliance.',
      details: [
        'HIPAA and GDPR compliant',
        'Multi-factor authentication',
        'End-to-end encryption',
        'Secure cloud storage'
      ]
    },
    {
      icon: <BarChart2 className="w-12 h-12" />,
      title: 'Advanced Analytics Dashboard',
      description: 'Comprehensive insights into medical practice performance and patient interactions.',
      details: [
        'Consultation trend analysis',
        'Performance metrics tracking',
        'Customizable reporting',
        'Export and sharing capabilities'
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
              Features • AI-Powered Tools
            </span>
          </div>
          <h1 className="text-4xl md:text-6xl font-bold mb-6 bg-gradient-to-br from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Transformative Features
          </h1>
          <p className="text-xl text-gray-400 mb-12 max-w-2xl mx-auto leading-relaxed transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            Empowering doctors with cutting-edge AI technology for seamless medical documentation
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 animate-fade-in-up">
          {featuresList.map((feature, index) => (
            <FeatureHighlight 
              key={index}
              icon={feature.icon}
              title={feature.title}
              description={feature.description}
              details={feature.details}
            />
          ))}
        </div>

        <section className="future-section mt-32 text-center bg-[#2a2a2a] bg-opacity-30 border border-[#3a3a3a] rounded-xl p-12 transform hover:scale-[1.01] hover:shadow-2xl hover:shadow-[#3ECF8E]/20 transition-all duration-500 will-change-transform transform-gpu">
          <h2 className="text-4xl md:text-5xl font-bold mb-8 bg-gradient-to-r from-[#3ECF8E] via-[#4BB563] to-white text-transparent bg-clip-text animate-gradient-x transform hover:scale-[1.01] transition-all duration-500 will-change-transform transform-gpu">
            Future-Ready Healthcare Technology
          </h2>
          <p className="text-xl text-gray-400 max-w-3xl mx-auto leading-relaxed mb-12 transform hover:scale-[1.01] transition-all duration-300 will-change-transform transform-gpu">
            Our AI-driven solution is continuously evolving, integrating the latest advancements 
            in natural language processing and machine learning to provide unparalleled support 
            for medical professionals.
          </p>
          <div className="inline-flex items-center bg-[#2a2a2a] border border-[#3a3a3a] rounded-full px-8 py-3 text-white transform hover:scale-[1.01] transition-all duration-300 cursor-pointer will-change-transform transform-gpu">
            <span className="bg-gradient-to-r from-[#3ECF8E] via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text font-semibold">
              Explore Future Innovations
            </span>
          </div>
        </section>
      </div>
    </div>
  );
};

export default Features;
