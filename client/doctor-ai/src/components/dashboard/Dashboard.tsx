import React, { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  BarChart, 
  Bar, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  ResponsiveContainer, 
  PieChart, 
  Pie, 
  Cell,
  Sector
} from 'recharts';
import { 
  Users, 
  MessageCircle, 
  ArrowUpRight,
  Stethoscope,
  Clock,
  Bot
  // ClipboardList
} from 'lucide-react';
import DoctorProfile from '../profile/DoctorProfile';
import ChatWindow from '../chat/ChatWindow';
import { dashboardApi } from '../../api/dashboard';

// Define types for our graph data
interface DailyStatData {
  date: string;
  count: number;
}

interface ApiStatistic {
  date: string;
  count: number;
}

interface PatientData {
  name: string;
  value: number;
  description?: string;
}

interface Patient {
  id: string;
  name: string;
  created_at: string;
}

interface ApiPatientResponse {
  patients: Patient[];
  totalPatients: number;
  totalTranscriptions: number;
  patientTypes: PatientData[];
}

interface DashboardData {
  dailyStats: DailyStatData[];
  monthlyStats: DailyStatData[];
  busiestDays: DailyStatData[];
  patientStats: PatientData[];
  recentPatients: Patient[];
  totalPatients: number;
  totalTranscriptions: number;
}

const OverviewCard: React.FC<{ 
  title: string; 
  value: number; 
  trend: 'up' | 'down';
  icon: React.ReactNode;
}> = ({ title, value, trend, icon }) => (
  <div className="bg-[#2a2a2a] rounded-2xl p-6 border border-[#3a3a3a] flex items-center justify-between">
    <div>
      <p className="text-sm text-gray-400 mb-2">{title}</p>
      <h3 className="text-2xl font-bold text-white">{value}</h3>
    </div>
    <div className="flex items-center">
      {icon}
      {trend === 'up' ? (
        <ArrowUpRight className="h-5 w-5 text-green-500 ml-2" />
      ) : (
        <ArrowUpRight className="h-5 w-5 text-red-500 ml-2 rotate-180" />
      )}
    </div>
  </div>
);

const RecentPatients: React.FC<{ patients: Patient[] }> = ({ patients }) => {
  return (
    <motion.div 
      className="bg-[#2a2a2a] rounded-2xl p-6 border border-[#3a3a3a]"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, delay: 0.2 }}
    >
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-semibold text-white flex items-center gap-3">
          <Stethoscope className="w-5 h-5 text-[#3ECF8E]" />
          Recent Patients
        </h3>
        <Clock className="w-4 h-4 text-gray-400" />
      </div>
      <div className="max-h-[300px] overflow-y-auto scrollbar-thin scrollbar-thumb-[#3a3a3a] scrollbar-track-[#1C1C1C]">
        {patients.length === 0 ? (
          <div className="text-center text-gray-400 py-4">
            No recent patients found
          </div>
        ) : (
          <div className="space-y-3">
            {patients.map((patient, index) => (
              <motion.div 
                key={patient.id} 
                className="flex items-center p-3 bg-[#1C1C1C] rounded-lg hover:bg-[#2a2a2a] transition-colors group"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ 
                  duration: 0.3, 
                  delay: index * 0.1 
                }}
                whileHover={{ 
                  scale: 1.025,
                  boxShadow: '0 5px 15px rgba(62, 207, 142, 0.1)'
                }}
              >
                <div className="flex-1">
                  <p className="text-white font-medium group-hover:text-[#3ECF8E] transition-colors">
                    {patient.name}
                  </p>
                  <p className="text-xs text-gray-400">
                    {new Date(patient.created_at).toLocaleDateString()}
                  </p>
                </div>
              </motion.div>
            ))}
          </div>
        )}
      </div>
    </motion.div>
  );
};

const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isChatOpen, setIsChatOpen] = useState(true);
  const [dashboardData, setDashboardData] = useState<DashboardData>({
    monthlyStats: [],
    busiestDays: [],
    dailyStats: [],
    patientStats: [],
    recentPatients: [],
    totalPatients: 0,
    totalTranscriptions: 0
  });

  useEffect(() => {
    // Reset scroll position on mount and after data load
    window.scrollTo(0, 0);
    const timer = setTimeout(() => window.scrollTo(0, 0), 100);
    return () => clearTimeout(timer);
  }, [loading]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        setError(null);

        // Fetch all data in parallel
        const [
          dailyResponse,
          monthlyResponse,
          busiestResponse,
          patientsResponse,
          transcriptsResponse
        ] = await Promise.all([
          dashboardApi.getDailyStatistics(30),
          dashboardApi.getMonthlyStatistics(6),
          dashboardApi.getBusiestDays(30),
          dashboardApi.getDashboardPatients(30),
          dashboardApi.getDashboardTranscripts(30)
        ]);

        // Transform daily stats
        const dailyStats = dailyResponse?.map((stat: ApiStatistic) => ({
          date: new Date(stat.date).toLocaleDateString('en-US', { 
            month: 'short', 
            day: 'numeric' 
          }),
          count: stat.count
        })) || [];

        // Transform monthly stats
        const monthlyStats = monthlyResponse?.map((stat: ApiStatistic) => ({
          date: new Date(stat.date).toLocaleDateString('en-US', { 
            year: 'numeric',
            month: 'short'
          }),
          count: stat.count
        })) || [];

        // Transform busiest days
        const busiestDays = busiestResponse?.map((stat: ApiStatistic) => ({
          date: new Date(stat.date).toLocaleDateString('en-US', { 
            month: 'short', 
            day: 'numeric' 
          }),
          count: stat.count
        })) || [];

        // Get patient types from response
        const patientData = patientsResponse as ApiPatientResponse;
        const patients = patientData?.patients || [];
        const totalPatientsCount = patientData?.totalPatients || 0;

        setDashboardData({
          dailyStats,
          monthlyStats,
          busiestDays,
          patientStats: patientData.patientTypes || [
            { name: 'New Patients', value: 0, description: 'First-time consultations' },
            { name: 'Regular Patients', value: 0, description: 'Consistent follow-up visits' },
            { name: 'Returning Patients', value: 0, description: 'Returning after 6+ months' },
            { name: 'Emergency Cases', value: 0, description: 'Urgent care visits' }
          ],
          recentPatients: patients,
          totalPatients: totalPatientsCount,
          totalTranscriptions: transcriptsResponse?.totalTranscriptions || 0
        });
      } catch (error) {
        console.error('Failed to fetch dashboard data', error);
        setError(error instanceof Error ? error.message : 'Failed to fetch dashboard data');
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  // Custom color gradients with modern colors
  const COLOR_GRADIENTS = [
    '#3ECF8E',   // Primary green
    '#6366F1',   // Indigo
    '#F59E0B',   // Amber
    '#EC4899',   // Pink
    '#8B5CF6',   // Purple
    '#10B981'    // Emerald
  ];

  // Add subtle animations
  useEffect(() => {
    const style = document.createElement('style');
    style.textContent = `
      @keyframes floatUp {
        0%, 100% { transform: translateY(0); }
        50% { transform: translateY(-4px); }
      }
      @keyframes pulse {
        0%, 100% { opacity: 1; transform: scale(1); }
        50% { opacity: 0.9; transform: scale(0.99); }
      }
      @keyframes glow {
        0%, 100% { filter: brightness(1); }
        50% { filter: brightness(1.1); }
      }
      @keyframes slideIn {
        from { transform: translateX(-10px); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
      }
      .legend-item {
        animation: slideIn 0.3s ease-out forwards;
        transition: all 0.2s ease;
      }
      .legend-item:hover {
        transform: translateX(4px);
        background: rgba(62, 207, 142, 0.1) !important;
      }
    `;
    document.head.appendChild(style);
    return () => {
      document.head.removeChild(style);
      // Return void to satisfy TypeScript
      return;
    };
  }, []);

  const [activeIndex, setActiveIndex] = useState<number>(-1);

  const renderActiveShape = (props: any) => {
    const RADIAN = Math.PI / 180;
    const { 
      cx, cy, midAngle, innerRadius, outerRadius, 
      startAngle, endAngle, payload, percent 
    } = props;

    const sin = Math.sin(-midAngle * RADIAN);
    const cos = Math.cos(-midAngle * RADIAN);
    const sx = cx + (outerRadius + 10) * cos;
    const sy = cy + (outerRadius + 10) * sin;
    const mx = cx + (outerRadius + 30) * cos;
    const my = cy + (outerRadius + 30) * sin;
    const ex = mx + (cos >= 0 ? 1 : -1) * 15;
    const ey = my;

    const gradientColor = COLOR_GRADIENTS[dashboardData.patientStats.findIndex(stat => stat.name === payload.name)];

    return (
      <g style={{ pointerEvents: 'none' }}>
        <Sector
          cx={cx}
          cy={cy}
          innerRadius={innerRadius}
          outerRadius={outerRadius}
          startAngle={startAngle}
          endAngle={endAngle}
          fill={gradientColor}
        />
        <Sector
          cx={cx}
          cy={cy}
          startAngle={startAngle}
          endAngle={endAngle}
          innerRadius={outerRadius + 6}
          outerRadius={outerRadius + 10}
          fill={gradientColor}
        />
        <path
          d={`M${sx},${sy}L${mx},${my}L${ex},${ey}`}
          stroke={gradientColor}
          strokeWidth={2}
          fill="none"
        />
        <text
          x={ex + (cos >= 0 ? 1 : -1) * 15}
          y={ey}
          textAnchor={cos >= 0 ? 'start' : 'end'}
          fill="white"
          className="text-xs font-medium"
          style={{ zIndex: 10, position: 'relative', textShadow: '0 1px 3px rgba(0,0,0,0.3)' }}
        >
          {`${payload.name} (${(percent * 100).toFixed(0)}%)`}
        </text>
      </g>
    );
  };

  if (loading) {
    return (
      <div className="fixed inset-0 bg-[#1C1C1C] flex items-center justify-center z-[100]">
        <div className="text-white text-lg">Loading dashboard data...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="fixed inset-0 bg-[#1C1C1C] flex items-center justify-center z-[100]">
        <div className="text-red-500 text-lg">{error}</div>
      </div>
    );
  }

  return (
    <motion.div 
      className="min-h-screen bg-[#1C1C1C]"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 0.3 }}
    >
      <div className="container mx-auto px-4 py-4 pt-16 relative">
        <motion.div 
          className="grid grid-cols-1 lg:grid-cols-3 gap-8 relative"
          initial="hidden"
          animate="visible"
        >
          {/* Doctor Profile */}
          <motion.div 
            className="lg:col-span-1"
            variants={{
              hidden: { opacity: 0, y: 20 },
              visible: { opacity: 1, y: 0 }
            }}
          >
            <div className="sticky top-20 space-y-6 z-[2]">
              <div className="relative bg-[#1C1C1C]">
                <DoctorProfile />
              </div>
              <AnimatePresence mode="wait">
                {isChatOpen ? (
                  <motion.div 
                    key="chat"
                    className="relative"
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: 10 }}
                    transition={{ type: 'spring', stiffness: 400, damping: 35 }}
                  >
                    <ChatWindow 
                      className="h-[300px]" 
                      onClose={() => setIsChatOpen(false)}
                    />
                  </motion.div>
                ) : (
                  <motion.button
                    key="button"
                    onClick={() => setIsChatOpen(true)}
                    className="relative w-full flex items-center justify-center gap-2 p-4 bg-[#2a2a2a] rounded-2xl border border-[#3a3a3a] hover:bg-[#333] transition-colors duration-200"
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: 10 }}
                    transition={{ type: 'spring', stiffness: 400, damping: 35 }}
                  >
                    <Bot className="w-5 h-5 text-emerald-500" />
                    <span className="text-white font-medium">Open AI Assistant</span>
                  </motion.button>
                )}
              </AnimatePresence>
            </div>
          </motion.div>

          {/* Navigation */}
          <motion.div 
            className="lg:col-span-2 relative z-[1]"
            variants={{
              hidden: { opacity: 0, y: 20 },
              visible: { opacity: 1, y: 0 }
            }}
          >
            {/* <div className="flex items-center space-x-6 mb-8 border-b border-[#3a3a3a] pb-4">
              <a
                href="http://localhost:8080/transcription"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg transition-all duration-300 text-gray-400 hover:text-white hover:bg-[#2a2a2a]"
              >
                <ClipboardList className="w-5 h-5" />
                <span>Logs</span>
              </a>
              <a
                href="http://localhost:8080/patients/"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center space-x-2 px-4 py-2 rounded-lg transition-all duration-300 text-gray-400 hover:text-white hover:bg-[#2a2a2a]"
              >
                <Users className="w-5 h-5" />
                <span>Patients</span>
              </a>
            </div> */}

            {/* Dashboard Content */}
            {
              <motion.div 
                className="space-y-6"
                variants={{
                  hidden: { opacity: 0, y: 20 },
                  visible: { opacity: 1, y: 0 }
                }}
          >
            {/* Overview Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <OverviewCard 
                title="Total Patients" 
                value={dashboardData.totalPatients} 
                trend="up"
                icon={<Users className="h-6 w-6 text-primary" />}
              />
              <OverviewCard 
                title="Consultations" 
                value={dashboardData.totalTranscriptions} 
                trend="down"
                icon={<MessageCircle className="h-6 w-6 text-primary" />}
              />
            </div>

            {/* Recent Patients */}
            <RecentPatients patients={dashboardData.recentPatients} />

            {/* Daily Transcriptions Bar Chart */}
            <motion.div 
              className="bg-[#2a2a2a] rounded-2xl p-6 border border-[#3a3a3a] relative"
              whileHover={{ scale: 1.02 }}
            >
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                <Clock className="w-5 h-5 mr-2 text-[#3ECF8E]" />
                Daily Transcriptions
              </h3>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={dashboardData.dailyStats}>
                  <CartesianGrid 
                    strokeDasharray="3 3" 
                    stroke="#3a3a3a" 
                    opacity={0.5}
                  />
                  <XAxis 
                    dataKey="date" 
                    tick={{ fill: 'white' }} 
                    interval={2}
                    angle={-45}
                    textAnchor="end"
                    height={60}
                  />
                  <YAxis 
                    tick={{ fill: 'white' }} 
                    width={40}
                  />
                  <Bar 
                    dataKey="count" 
                    fill="#60efbc" 
                    barSize={20} 
                    radius={[4, 4, 0, 0]}
                    animationDuration={500}
                    animationEasing="ease"
                  />
                </BarChart>
              </ResponsiveContainer>
            </motion.div>

            {/* Monthly Statistics Line Chart */}
            <motion.div 
              className="bg-[#2a2a2a] rounded-2xl p-6 border border-[#3a3a3a] relative"
              whileHover={{ scale: 1.02 }}
            >
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                <Clock className="w-5 h-5 mr-2 text-[#3ECF8E]" />
                Monthly Trends
              </h3>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={dashboardData.monthlyStats}>
                  <CartesianGrid 
                    strokeDasharray="3 3" 
                    stroke="#3a3a3a" 
                    opacity={0.5}
                  />
                  <XAxis 
                    dataKey="date" 
                    tick={{ fill: 'white' }} 
                    interval={0}
                    angle={-45}
                    textAnchor="end"
                    height={60}
                  />
                  <YAxis 
                    tick={{ fill: 'white' }} 
                    width={40}
                  />
                  <Bar 
                    dataKey="count" 
                    fill="#949aff" 
                    barSize={30} 
                    radius={[4, 4, 0, 0]}
                    animationDuration={500}
                    animationEasing="ease"
                  />
                </BarChart>
              </ResponsiveContainer>
            </motion.div>

            {/* Busiest Days Stats */}
            <motion.div 
              className="bg-[#2a2a2a] rounded-2xl p-6 border border-[#3a3a3a]"
              whileHover={{ scale: 1.02 }}
            >
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                <Clock className="w-5 h-5 mr-2 text-[#F59E0B]" />
                Busiest Days
              </h3>
              <div className="space-y-4">
                {dashboardData.busiestDays.map((day, index) => (
                  <div key={day.date} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div className="w-8 h-8 rounded-full bg-[#1C1C1C] flex items-center justify-center text-[#F59E0B] font-semibold">
                        {index + 1}
                      </div>
                      <span className="text-white">{day.date}</span>
                    </div>
                    <div className="flex items-center space-x-2">
                      <span className="text-[#F59E0B] font-semibold">{day.count}</span>
                      <span className="text-gray-400">consultations</span>
                    </div>
                  </div>
                ))}
              </div>
            </motion.div>

            {/* Patient Types 3D Pie Chart */}
            <motion.div 
              className="bg-[#2a2a2a] rounded-2xl p-8 border border-[#3a3a3a] relative z-10"
              whileHover={{ scale: 1.02 }}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.4 }}
            >
              <div className="flex items-center justify-between mb-8">
                <div className="flex items-center space-x-3">
                  <div className="p-2 bg-[#EF4444]/10 rounded-lg">
                    <Stethoscope className="w-5 h-5 text-[#EF4444]" />
                  </div>
                  <h3 className="text-lg font-semibold text-white">Patient Types</h3>
                </div>
              </div>
              <div className="grid grid-cols-12 gap-8 items-center justify-center" style={{ position: 'relative' }}>
                <div className="col-span-3 flex flex-col space-y-2 mt-8" style={{ marginTop: '-1.5rem', marginLeft: '2.5rem' }}>
                  {dashboardData.patientStats.map((_, index) => (
                    <div 
                      key={index} 
                      className="legend-item flex items-center space-x-3 p-3 rounded-lg transition-all duration-300 hover:bg-[#3a3a3a]/50"
                    >
                      <div 
                        className={`h-4 w-6 rounded-full`}
                        style={{ 
                          background: COLOR_GRADIENTS[index],
                          boxShadow: `0 0 10px ${COLOR_GRADIENTS[index]}40`
                        }}
                      />
                      <div className="flex-1">
                        <p className="font-medium text-sm text-white/90">{dashboardData.patientStats[index].name}</p>
                        <p className="text-xs text-white/50">{dashboardData.patientStats[index].description}</p>
                      </div>
                    </div>
                  ))}
                </div>
                <div className="col-span-9 overflow-visible flex items-center" style={{ zIndex: 20, position: 'relative', marginTop: '-2rem' }}>
                  <ResponsiveContainer width="100%" height={360}>
                    <PieChart style={{ overflow: 'visible' }} margin={{ top: 10, right: 60, bottom: 10, left: 60 }}>
                      <Pie
                        data={dashboardData.patientStats}
                        cx="52%"
                        cy="50%"
                        labelLine={false}
                        innerRadius={60}
                        outerRadius={130}
                        paddingAngle={4}
                        dataKey="value"
                        activeIndex={activeIndex}
                        activeShape={renderActiveShape}
                        onMouseEnter={(_, index) => setActiveIndex(index)}
                        onMouseLeave={() => setActiveIndex(-1)}
                        animationBegin={0}
                        animationDuration={1000}
                        animationEasing="ease-out"
                      >
                        {dashboardData.patientStats.map((_, index) => (
                          <Cell
                            key={`cell-${index}`}
                            fill={COLOR_GRADIENTS[index]}
                            style={{
                              filter: activeIndex === index 
                                ? `drop-shadow(0 0 15px ${COLOR_GRADIENTS[index]})` 
                                : 'none',
                              transition: 'all 0.3s ease-in-out',
                              opacity: activeIndex === -1 || activeIndex === index ? 1 : 0.6,
                              border: `1px solid ${COLOR_GRADIENTS[index]}`,
                              boxShadow: `0 0 8px ${COLOR_GRADIENTS[index]}`,
                            }}
                          />
                        ))}
                      </Pie>
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </motion.div>
              </motion.div>
          }          </motion.div>
        </motion.div>
      </div>
    </motion.div>
  );
};

export default Dashboard;
