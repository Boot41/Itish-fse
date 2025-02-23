import { format, subDays } from 'date-fns';

const patientNames = [
  'Rahul Sharma', 'Priya Patel', 'Amit Kumar', 'Neha Singh', 'Raj Malhotra',
  'Anita Desai', 'Suresh Verma', 'Meera Reddy', 'Vikram Mehta', 'Pooja Gupta',
  'Arun Joshi', 'Kavita Nair', 'Deepak Shah', 'Anjali Kapoor', 'Sanjay Chopra'
];

const conditions = [
  'Hypertension', 'Diabetes', 'Asthma', 'Arthritis', 'Migraine',
  'Anxiety', 'Depression', 'Back Pain', 'Allergies', 'Insomnia'
];

const transcriptTemplates = [
  'Patient complains of {condition} symptoms. Prescribed medication and advised follow-up in 2 weeks.',
  'Regular check-up for {condition}. Condition stable, continuing current treatment plan.',
  'Initial consultation for {condition}. Started on new treatment regimen.',
  'Follow-up visit for {condition}. Showing improvement, adjusting medication dosage.',
  'Emergency visit due to severe {condition} symptoms. Provided immediate treatment.'
];

const generateDates = (days: number) => {
  const dates = [];
  const today = new Date();
  for (let i = days - 1; i >= 0; i--) {
    const date = subDays(today, i);
    dates.push(format(date, 'yyyy-MM-dd'));
  }
  return dates;
};

const generatePatients = (count: number, maxDays: number) => {
  return Array.from({ length: count }, (_, i) => {
    const condition = conditions[Math.floor(Math.random() * conditions.length)];
    return {
      id: `patient-${i + 1}`,
      name: patientNames[i % patientNames.length],
      age: Math.floor(Math.random() * 50) + 20,
      gender: Math.random() > 0.5 ? 'Male' : 'Female',
      condition,
      created_at: format(subDays(new Date(), Math.floor(Math.random() * maxDays)), 'yyyy-MM-dd'),
    };
  });
};

const generateTranscripts = (count: number, maxDays: number) => {
  return Array.from({ length: count }, (_, i) => {
    const condition = conditions[Math.floor(Math.random() * conditions.length)];
    const template = transcriptTemplates[Math.floor(Math.random() * transcriptTemplates.length)];
    return {
      id: `transcript-${i + 1}`,
      text: template.replace('{condition}', condition),
      report: `Medical report for ${condition} treatment and follow-up care plan`,
      created_at: format(subDays(new Date(), Math.floor(Math.random() * maxDays)), 'yyyy-MM-dd'),
    };
  });
};

const generateRealisticCounts = (days: number, avgPerDay: number, variance: number = 0.3) => {
  // Generate a more realistic pattern with some randomness but following a trend
  return Array.from({ length: days }, () => {
    const base = avgPerDay * (1 + (Math.random() - 0.5) * variance);
    return Math.max(1, Math.floor(base));
  });
};

export const dummyData = {
  dailyStatistics: (days: number = 30) => {
    const dates = generateDates(days);
    const counts = generateRealisticCounts(days, 8); // Average 8 patients per day
    return dates.map((date, i) => ({
      date,
      count: counts[i],
    }));
  },

  monthlyStatistics: (months: number = 6) => {
    const dates = Array.from({ length: months }, (_, i) => {
      const date = new Date();
      date.setMonth(date.getMonth() - i);
      return format(date, 'yyyy-MM');
    }).reverse();
    const counts = generateRealisticCounts(months, 180); // Average 180 patients per month
    return dates.map((date, i) => ({
      date,
      count: counts[i],
    }));
  },

  busiestDays: (days: number = 30) => {
    const dates = generateDates(days);
    const counts = generateRealisticCounts(days, 12, 0.5); // Busier days
    return dates.map((date, i) => ({
      date,
      count: counts[i],
    }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5);
  },

  dashboardPatients: (days: number = 30) => {
    const patients = generatePatients(10, days);
    return {
      patients,
      totalPatients: 450,
      totalTranscriptions: 850,
      patientTypes: [
        { name: 'New Patients', value: 120, description: 'First-time consultations' },
        { name: 'Regular Patients', value: 200, description: 'Consistent follow-up visits' },
        { name: 'Returning Patients', value: 80, description: 'Returning after 6+ months' },
        { name: 'Emergency Cases', value: 50, description: 'Urgent care visits' }
      ],
    };
  },

  dashboardTranscripts: (days: number = 30) => {
    const transcripts = generateTranscripts(10, days);
    return {
      transcripts,
      totalTranscriptions: 850,
    };
  },
};
