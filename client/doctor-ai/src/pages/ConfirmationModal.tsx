import React from 'react';
import { motion } from 'framer-motion';
import { X } from 'lucide-react';

interface ConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  patientName: string;
}

const ConfirmationModal: React.FC<ConfirmationModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  patientName,
}) => {
  if (!isOpen) return null;

  return (
    <div 
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      onClick={onClose}
    >
      <motion.div
        onClick={(e) => e.stopPropagation()}
        initial={{ opacity: 0, scale: 0.9, y: 50 }}
        animate={{ opacity: 1, scale: 1, y: 0 }}
        exit={{ opacity: 0, scale: 0.9, y: 50 }}
        transition={{ type: 'spring', stiffness: 300, damping: 20 }}
        className="bg-[#1C1C1C] border border-[#3a3a3a] p-6 rounded-lg w-full max-w-md mx-4 shadow-2xl"
      >
        <div className="flex justify-between items-center">
          <h2 className="text-xl font-bold text-white">Are you sure?</h2>
          <button 
            onClick={onClose} 
            className="p-2 hover:bg-[#3a3a3a] rounded-full"
          >
            <X className="w-5 h-5 text-gray-400" />
          </button>
        </div>
        <div className="mt-4">
          <p className="text-gray-300">
            Do you really want to delete <span className="font-semibold">{patientName}</span>? This action cannot be undone.
          </p>
        </div>
        <div className="mt-6 flex justify-end space-x-4">
          <button 
            onClick={onClose}
            className="px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-600 transition"
          >
            Cancel
          </button>
          <button 
            onClick={onConfirm}
            className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition"
          >
            Delete
          </button>
        </div>
      </motion.div>
    </div>
  );
};

export default ConfirmationModal;
