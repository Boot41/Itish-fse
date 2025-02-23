import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { UserPlus, Calendar, Link, Loader2, X } from 'lucide-react';
import { Toaster, toast } from 'sonner';
import { Patient } from '../../types/Patient';
import { useAuthStore } from '../../stores/authStore';
import { useRefreshStore } from '../../stores/refreshStore';

interface CreatePatientModalProps {
  isOpen: boolean;
  onClose: () => void;
  onPatientCreated: (patient: Patient) => void;
}

const CreatePatientModal: React.FC<CreatePatientModalProps> = ({ 
  isOpen, 
  onClose, 
  onPatientCreated 
}) => {
  const [formData, setFormData] = useState({
    name: '',
    age: '',
    gender: '',
    url: '',
  });
  const [errors, setErrors] = useState<{ [key: string]: string }>({});
  const [isLoading, setIsLoading] = useState(false);

  const { getDecryptedEmail } = useAuthStore();
  const [doctorEmail, setDoctorEmail] = useState<string | null>(null);
  const { triggerPatientRefresh, triggerTranscriptRefresh } = useRefreshStore();

  useEffect(() => {
    const decryptEmail = async () => {
      try {
        const decrypted = await getDecryptedEmail();
        setDoctorEmail(decrypted);
      } catch (error) {
        console.error('Failed to decrypt email:', error);
        toast.error('Failed to retrieve doctor email');
      }
    };
    decryptEmail();
  }, [getDecryptedEmail]);

  const validateForm = () => {
    const newErrors: { [key: string]: string } = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Patient name is required';
    }
    if (!formData.age.trim()) {
      newErrors.age = 'Age is required';
    } else if (isNaN(parseInt(formData.age)) || parseInt(formData.age) <= 0) {
      newErrors.age = 'Invalid age';
    }
    if (!formData.gender.trim()) {
      newErrors.gender = 'Gender is required';
    }
    if (!formData.url.trim()) {
      newErrors.url = 'Audio URL is required';
    } else if (!formData.url.startsWith('http')) {
      newErrors.url = 'Audio URL must start with http or https';
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;
    if (!doctorEmail) {
      toast.error('Doctor email not found. Please log in again.');
      return;
    }
    if (isLoading) return;

    setIsLoading(true);
    try {
      const response = await fetch('http://localhost:8080/patients/', {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
          'Origin': import.meta.env.VITE_FRONTEND_URL,
        },
        body: JSON.stringify({
          email: doctorEmail,
          patient: {
            name: formData.name,
            age: parseInt(formData.age),
            gender: formData.gender,
            url: formData.url,
          },
          audio: formData.url,
        }),
      });
      if (!response.ok) {
        const errorText = await response.text();
        let errorData;
        try {
          errorData = JSON.parse(errorText);
          throw new Error(errorData.message || 'Failed to create patient');
        } catch {
          throw new Error(errorText || `HTTP error! status: ${response.status}`);
        }
      }
      // Construct a minimal Patient object
      const patient: Patient = {
        id: '',
        name: formData.name,
        age: parseInt(formData.age),
        gender: formData.gender,
      };
      // Refresh the global patient/transcript lists
      triggerPatientRefresh();
      triggerTranscriptRefresh();
      toast.success('Patient created successfully!');
      onPatientCreated(patient);
      onClose();
    } catch (error: any) {
      console.error('Patient creation error:', error);
      toast.error(error.message || 'Failed to create patient');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  return (
    <>
      <Toaster position="top-right" richColors theme="dark" />
      <AnimatePresence>
        {isOpen && (
          <div 
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
            onClick={onClose}
          >
            <motion.div
              variants={{
                hidden: { opacity: 0, scale: 0.9, y: 50 },
                visible: {
                  opacity: 1,
                  scale: 1,
                  y: 0,
                  transition: { type: 'spring', stiffness: 300, damping: 20 },
                },
                exit: { opacity: 0, scale: 0.9, y: 50, transition: { duration: 0.2 } },
              }}
              initial="hidden"
              animate="visible"
              exit="exit"
              className="relative bg-[#2a2a2a] rounded-2xl p-8 w-full max-w-md mx-4 border border-[#3a3a3a] shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            >
              <button
                onClick={onClose}
                className="absolute top-4 right-4 text-gray-400 hover:text-white transition-colors"
                disabled={isLoading}
              >
                <X className="w-6 h-6" />
              </button>
              <div className="text-center mb-6">
                <div className="mx-auto w-16 h-16 bg-[#3ECF8E]/10 rounded-full flex items-center justify-center mb-4">
                  <UserPlus className="w-8 h-8 text-[#3ECF8E]" />
                </div>
                <h2 className="text-2xl font-bold text-white">Create New Patient</h2>
                <p className="text-gray-400 mt-2">Enter patient details for transcription</p>
              </div>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label htmlFor="name" className="block text-sm font-medium text-gray-300 mb-1">
                    Patient Name
                  </label>
                  <div className="flex items-center border border-[#3a3a3a] rounded-lg">
                    <UserPlus className="w-5 h-5 ml-3 text-gray-400" />
                    <input
                      type="text"
                      id="name"
                      name="name"
                      value={formData.name}
                      onChange={handleChange}
                      placeholder="Enter patient name"
                      className="w-full bg-transparent p-2 text-white placeholder-gray-500 focus:outline-none"
                      disabled={isLoading}
                    />
                  </div>
                  {errors.name && <p className="text-red-500 text-xs mt-1">{errors.name}</p>}
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label htmlFor="age" className="block text-sm font-medium text-gray-300 mb-1">
                      Age
                    </label>
                    <div className="flex items-center border border-[#3a3a3a] rounded-lg">
                      <Calendar className="w-5 h-5 ml-3 text-gray-400" />
                      <input
                        type="number"
                        id="age"
                        name="age"
                        value={formData.age}
                        onChange={handleChange}
                        placeholder="Patient age"
                        className="w-full bg-transparent p-2 text-white placeholder-gray-500 focus:outline-none"
                        disabled={isLoading}
                      />
                    </div>
                    {errors.age && <p className="text-red-500 text-xs mt-1">{errors.age}</p>}
                  </div>
                  <div>
                    <label htmlFor="gender" className="block text-sm font-medium text-gray-300 mb-1">
                      Gender
                    </label>
                    <select
                      id="gender"
                      name="gender"
                      value={formData.gender}
                      onChange={handleChange}
                      className="w-full bg-[#1C1C1C] p-2 rounded-lg border border-[#3a3a3a] text-white"
                      disabled={isLoading}
                    >
                      <option value="">Select Gender</option>
                      <option value="male">Male</option>
                      <option value="female">Female</option>
                      <option value="other">Other</option>
                    </select>
                    {errors.gender && <p className="text-red-500 text-xs mt-1">{errors.gender}</p>}
                  </div>
                </div>
                <div>
                  <label htmlFor="url" className="block text-sm font-medium text-gray-300 mb-1">
                    Audio URL
                  </label>
                  <div className="flex items-center border border-[#3a3a3a] rounded-lg">
                    <Link className="w-5 h-5 ml-3 text-gray-400" />
                    <input
                      type="text"
                      id="url"
                      name="url"
                      value={formData.url}
                      onChange={handleChange}
                      placeholder="Enter audio transcription URL"
                      className="w-full bg-transparent p-2 text-white placeholder-gray-500 focus:outline-none"
                      disabled={isLoading}
                    />
                  </div>
                  {errors.url && <p className="text-red-500 text-xs mt-1">{errors.url}</p>}
                </div>
                <button
                  type="submit"
                  className={`w-full flex items-center justify-center p-3 rounded-lg transition-all duration-300 ${
                    isLoading
                      ? 'bg-[#3ECF8E]/50 cursor-not-allowed'
                      : 'bg-[#3ECF8E] hover:bg-[#45E6A0] active:scale-95'
                  } text-black font-semibold`}
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <div className="flex items-center">
                      <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                      Creating Patient...
                    </div>
                  ) : (
                    'Create Patient'
                  )}
                </button>
              </form>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </>
  );
};

export default CreatePatientModal;
