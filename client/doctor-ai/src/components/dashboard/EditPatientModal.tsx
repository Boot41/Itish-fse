import React, { useState, useEffect, FormEvent, ChangeEvent } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

interface EditPatientModalProps {
  show: boolean;
  handleClose: () => void;
  patientId: string;
  doctorEmail: string;
  onRefresh: () => void;
  initialData?: { name: string; age: string; gender: string } | null;
}

interface PatientData {
  name: string;
  age: string;
  gender: string;
}

const EditPatientModal: React.FC<EditPatientModalProps> = ({
  show,
  handleClose,
  patientId,
  doctorEmail,
  onRefresh,
  initialData = null,
}) => {
  const [patientData, setPatientData] = useState<PatientData>(
    initialData || { name: '', age: '', gender: '' }
  );
  const [isDataReady, setIsDataReady] = useState<boolean>(!!initialData);

  // Update patientData when initialData changes
  useEffect(() => {
    if (initialData) {
      setPatientData(initialData);
      setIsDataReady(true);
    }
  }, [initialData]);

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    fetch(`http://localhost:8080/patients/${patientId}/update`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
      },
      body: JSON.stringify({
        email: doctorEmail,
        patient_id: patientId,
        update_data: patientData,
      }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error('Failed to update patient');
        }
        onRefresh();
        handleClose();
      })
      .catch((error) => console.error('Error updating patient:', error));
  };

  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setPatientData((prev) => ({ ...prev, [name]: value }));
  };

  return (
    <AnimatePresence>
      {show && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
          onClick={handleClose}
        >
          <motion.div
            onClick={(e) => e.stopPropagation()}
            initial={{ opacity: 0, scale: 0.95, y: 20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 20 }}
            transition={{ type: 'spring', stiffness: 350, damping: 25 }}
            className="bg-[#2a2a2a] rounded-2xl p-8 w-full max-w-xl mx-4 border border-gray-700/50 shadow-2xl backdrop-blur-sm"
          >
            <div className="flex justify-between items-center mb-6">
              <div>
                <h2 className="text-2xl font-bold text-white">Edit Patient</h2>
                <p className="text-gray-400 text-sm mt-1">Update the patient information below</p>
              </div>
              <button 
                onClick={handleClose} 
                className="p-2 rounded-full hover:bg-gray-700/50 transition-colors duration-200"
              >
                <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            {!isDataReady ? (
              <div className="mt-4 flex items-center justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[#3ECF8E]"></div>
                <span className="ml-3 text-gray-400">Preparing form data...</span>
              </div>
            ) : (
              <form onSubmit={handleSubmit} className="space-y-6">
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-300 mb-2">Full Name</label>
                    <div className="relative">
                      <input
                        type="text"
                        name="name"
                        value={patientData.name}
                        onChange={handleChange}
                        placeholder="Enter patient's full name"
                        className="w-full px-4 py-3 bg-[#1C1C1C] border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[#3ECF8E] focus:border-transparent transition-all duration-200"
                      />
                      <div className="absolute inset-0 bg-gradient-to-b from-[#3ECF8E]/5 to-transparent rounded-xl pointer-events-none"></div>
                    </div>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-300 mb-2">Age</label>
                      <div className="relative">
                        <input
                          type="number"
                          name="age"
                          value={patientData.age}
                          onChange={handleChange}
                          placeholder="Enter age"
                          min="0"
                          max="150"
                          className="w-full px-4 py-3 bg-[#1C1C1C] border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[#3ECF8E] focus:border-transparent transition-all duration-200"
                        />
                        <div className="absolute inset-0 bg-gradient-to-b from-[#3ECF8E]/5 to-transparent rounded-xl pointer-events-none"></div>
                      </div>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-300 mb-2">Gender</label>
                      <div className="relative">
                        <select
                          name="gender"
                          value={patientData.gender}
                          onChange={handleChange}
                          className="w-full px-4 py-3 bg-[#1C1C1C] border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-[#3ECF8E] focus:border-transparent transition-all duration-200 appearance-none"
                        >
                          <option value="" className="bg-[#1C1C1C]">Select gender</option>
                          <option value="Male" className="bg-[#1C1C1C]">Male</option>
                          <option value="Female" className="bg-[#1C1C1C]">Female</option>
                          <option value="Other" className="bg-[#1C1C1C]">Other</option>
                        </select>
                        <div className="absolute inset-0 bg-gradient-to-b from-[#3ECF8E]/5 to-transparent rounded-xl pointer-events-none"></div>
                        <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                          <svg className="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                          </svg>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <div className="flex justify-end space-x-3 pt-4">
                  <button
                    type="button"
                    onClick={handleClose}
                    className="px-4 py-2.5 border border-gray-600 rounded-xl text-gray-300 hover:bg-gray-700/50 focus:outline-none focus:ring-2 focus:ring-gray-500 transition-all duration-200"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="px-6 py-2.5 bg-[#3ECF8E] text-black rounded-xl hover:bg-[#45E6A0] focus:outline-none focus:ring-2 focus:ring-[#3ECF8E] focus:ring-offset-2 focus:ring-offset-[#2a2a2a] transition-all duration-200 font-medium"
                  >
                    Save Changes
                  </button>
                </div>
              </form>
            )}
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
};

export default EditPatientModal;
