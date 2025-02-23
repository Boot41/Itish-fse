import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '../../stores/authStore';
import { profileApi } from '../../api/profile';
import type { Doctor } from '../../types/auth';
import { 
  User as UserIcon, 
  Mail, 
  Phone, 
  Stethoscope, 
  Calendar, 
  Edit3, 
  Save, 
  X,
  Star
} from 'lucide-react';

interface ProfileFieldProps {
  icon: React.ReactNode;
  label: string;
  value: string;
  editable?: boolean;
  isEditing?: boolean;
  name?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const fastTransition = { 
  duration: 0.2,  
  type: "tween",  
  ease: "easeOut" 
};

const quickHover = { 
  scale: 1.025,
  boxShadow: '0 5px 15px rgba(62, 207, 142, 0.1)',
  transition: {
    duration: 0.15,  
    type: "tween"
  }
};

const ProfileField: React.FC<ProfileFieldProps> = ({ 
  icon, 
  label, 
  value, 
  editable = false,
  isEditing = false,
  name,
  onChange
}) => {
  const [editedValue, setEditedValue] = useState(value);

  useEffect(() => {
    setEditedValue(value);
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditedValue(e.target.value);
    onChange && onChange(e);
  };

  return (
    <motion.div 
      className="flex items-center space-x-3 p-3 bg-[#2a2a2a] rounded-xl group border border-[#3a3a3a] hover:border-[#3ECF8E] transition-all duration-150"
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={fastTransition}
      whileHover={quickHover}
    >
      <div className="bg-[#3ECF8E]/10 p-2 rounded-full">
        {icon}
      </div>
      <div className="flex-1">
        <p className="text-xs text-gray-400">{label}</p>
        {!isEditing ? (
          <p className="text-white font-medium">{value || 'Not specified'}</p>
        ) : (
          <input 
            type="text"
            name={name}
            value={editedValue}
            onChange={handleChange}
            className="bg-[#1C1C1C] border border-[#3a3a3a] rounded px-2 py-1 text-white w-full focus:border-[#3ECF8E] transition-colors duration-150"
          />
        )}
      </div>
      {editable && isEditing && (
        <div className="flex space-x-2">
          <motion.button 
            onClick={() => {
              onChange && onChange({ 
                target: { 
                  name: name || '', 
                  value: editedValue 
                } 
              } as React.ChangeEvent<HTMLInputElement>);
            }}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            whileTap={{ scale: 0.95 }}
            transition={{ duration: 0.1 }}
            className="text-[#3ECF8E] hover:text-[#4BB563] transition-colors duration-150"
          >
            <Save className="w-4 h-4" />
          </motion.button>
          <motion.button 
            onClick={() => {
              setEditedValue(value);
            }}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            whileTap={{ scale: 0.95 }}
            transition={{ duration: 0.1 }}
            className="text-red-500 hover:text-red-400 transition-colors duration-150"
          >
            <X className="w-4 h-4" />
          </motion.button>
        </div>
      )}
    </motion.div>
  );
};

const DoctorProfile: React.FC = () => {
  const [profile, setProfile] = useState<Doctor | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [editedProfile, setEditedProfile] = useState<Doctor | null>(null);

  const { getDecryptedEmail } = useAuthStore();

  const fetchProfile = async () => {
    setLoading(true);
    setError('');
    try {
      const email = await getDecryptedEmail();
      if (!email) {
        setError('No email found in auth store');
        return;
      }

      const response = await profileApi.getProfile(email);
      if (response.success && response.data) {
        setProfile(response.data);
        setEditedProfile(response.data);
      } else {
        setError(response.message || 'Failed to fetch profile');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to fetch profile');
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = async () => {
    if (!editedProfile) return;

    setLoading(true);
    setError('');
    setSuccess('');
    try {
      const email = await getDecryptedEmail();
      if (!email) {
        setError('No email found in auth store');
        return;
      }

      const response = await profileApi.updateProfile(email, editedProfile);
      if (response.success) {
        setSuccess('Profile updated successfully');
        setProfile(editedProfile);
        setIsEditing(false);
      } else {
        setError(response.message || 'Failed to update profile');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to update profile');
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
    setEditedProfile(profile);
    setError('');
    setSuccess('');
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!editedProfile) return;
    const { name, value } = e.target;
    const field = name.charAt(0).toUpperCase() + name.slice(1);
    setEditedProfile({
      ...editedProfile,
      [field]: value,
    });
  };

  useEffect(() => {
    fetchProfile();
  }, []);

  if (loading && !profile) {
    return (
      <motion.div 
        className="bg-[#1C1C1C] rounded-2xl p-6 border border-[#3a3a3a] animate-pulse"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
      >
        <div className="h-20 bg-[#2a2a2a] rounded-full mb-4"></div>
        <div className="space-y-3">
          {[1, 2, 3, 4].map((_, index) => (
            <div key={index} className="h-10 bg-[#2a2a2a] rounded"></div>
          ))}
        </div>
      </motion.div>
    );
  }

  return (
    <motion.div 
      className="bg-[#1C1C1C] rounded-2xl p-6 border border-[#3a3a3a] space-y-6"
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ 
        duration: 0.3,  
        type: "spring",
        stiffness: 150  
      }}
    >
      {/* Profile Header */}
      <motion.div 
        className="flex items-center space-x-4 pb-4 border-b border-[#3a3a3a]"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1, duration: 0.3 }}  
      >
        <div className="relative">
          <img 
            src="https://i.pravatar.cc/150?img=4" 
            alt="Dr. Profile" 
            className="w-20 h-20 rounded-full border-4 border-[#3ECF8E]/30 hover:border-[#3ECF8E] transition-all"
          />
          <motion.div 
            className="absolute bottom-0 right-0 bg-[#3ECF8E] text-[#1C1C1C] rounded-full p-1"
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ 
              type: "spring", 
              stiffness: 300  
            }}
          >
            <Star className="w-4 h-4" />
          </motion.div>
        </div>
        <div>
          <h2 className="text-2xl font-bold text-white bg-gradient-to-r from-white via-[#3ECF8E] to-[#4BB563] text-transparent bg-clip-text">
            {profile?.Name || 'Dr. Unknown'}
          </h2>
          <p className="text-sm text-[#3ECF8E] opacity-80">
            {profile?.Specialization || 'No Specialization'}
          </p>
        </div>
        <div className="ml-auto flex items-center space-x-2">
          {!isEditing ? (
            <motion.button
              onClick={() => setIsEditing(true)}
              className="group bg-[#2a2a2a] border border-[#3a3a3a] hover:border-[#3ECF8E] text-white px-4 py-2 rounded-md font-medium transition-all duration-150"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
              transition={{ duration: 0.1 }}
            >
              <Edit3 className="inline-block mr-2 w-4 h-4 text-[#3ECF8E] transform group-hover:rotate-6 transition-transform duration-150" />
              Edit Profile
            </motion.button>
          ) : (
            <div className="flex space-x-2">
              <motion.button
                onClick={handleUpdate}
                className="group bg-gradient-to-r from-[#3ECF8E] to-[#4BB563] text-[#1C1C1C] px-4 py-2 rounded-md font-medium transition-all duration-150"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                transition={{ duration: 0.1 }}
              >
                <Save className="inline-block mr-2 w-4 h-4 transform group-hover:scale-110 transition-transform duration-150" />
                Save Changes
              </motion.button>
              <motion.button
                onClick={handleCancel}
                className="group bg-[#2a2a2a] border border-[#3a3a3a] hover:border-red-500 text-white px-4 py-2 rounded-md font-medium transition-all duration-150"
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                transition={{ duration: 0.1 }}
              >
                <X className="inline-block mr-2 w-4 h-4 text-red-500 transform group-hover:rotate-180 transition-transform duration-150" />
                Cancel
              </motion.button>
            </div>
          )}
        </div>
      </motion.div>

      {/* Error/Success Messages */}
      <AnimatePresence>
        {(error || success) && (
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 20 }}
            className={`p-3 rounded-lg ${
              error 
                ? 'bg-red-500/10 text-red-500 border border-red-500/30' 
                : 'bg-[#3ECF8E]/10 text-[#3ECF8E] border border-[#3ECF8E]/30'
            }`}
          >
            {error || success}
          </motion.div>
        )}
      </AnimatePresence>

      {/* Profile Details */}
      <div className="space-y-3">
        <ProfileField 
          icon={<UserIcon className="w-5 h-5 text-[#3ECF8E]" />}
          label="Full Name"
          value={profile?.Name || ''}
          editable={isEditing}
          isEditing={isEditing}
          name="name"
          onChange={handleChange}
        />
        
        <ProfileField 
          icon={<Mail className="w-5 h-5 text-[#3ECF8E]" />}
          label="Email Address"
          value={profile?.Email || ''}
          editable={false}
        />
        
        <ProfileField 
          icon={<Phone className="w-5 h-5 text-[#3ECF8E]" />}
          label="Phone Number"
          value={profile?.Phone || ''}
          editable={isEditing}
          isEditing={isEditing}
          name="phone"
          onChange={handleChange}
        />
        
        <ProfileField 
          icon={<Stethoscope className="w-5 h-5 text-[#3ECF8E]" />}
          label="Specialization"
          value={profile?.Specialization || ''}
          editable={isEditing}
          isEditing={isEditing}
          name="specialization"
          onChange={handleChange}
        />
        
        <ProfileField 
          icon={<Calendar className="w-5 h-5 text-[#3ECF8E]" />}
          label="Member Since"
          value={profile?.CreatedAt ? new Date(profile.CreatedAt).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
          }) : ''}
          editable={false}
        />
      </div>
    </motion.div>
  );
};

export default DoctorProfile;
