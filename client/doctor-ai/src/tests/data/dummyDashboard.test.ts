import { dummyData } from '../../data/dummyDashboard';

describe('Dummy Dashboard Data', () => {
  test('should have expected methods', () => {
    expect(dummyData).toHaveProperty('dashboardPatients');
    expect(dummyData).toHaveProperty('dashboardTranscripts');
    expect(dummyData).toHaveProperty('dailyStatistics');
    expect(dummyData).toHaveProperty('monthlyStatistics');
    expect(dummyData).toHaveProperty('busiestDays');
  });

  test('methods should return expected data structures', () => {
    const dashboardData = dummyData.dashboardPatients(7);
    const transcriptData = dummyData.dashboardTranscripts(7);
    
    // Check dashboard data structure
    expect(dashboardData).toHaveProperty('patients');
    expect(dashboardData).toHaveProperty('totalPatients');
    expect(dashboardData).toHaveProperty('totalTranscriptions');
    expect(dashboardData).toHaveProperty('patientTypes');
    
    // Check patients array
    expect(Array.isArray(dashboardData.patients)).toBe(true);
    if (dashboardData.patients.length > 0) {
      expect(dashboardData.patients[0]).toHaveProperty('name');
      expect(dashboardData.patients[0]).toHaveProperty('condition');
      expect(dashboardData.patients[0]).toHaveProperty('created_at');
    }
    
    // Check transcripts data structure
    expect(transcriptData).toHaveProperty('transcripts');
    expect(transcriptData).toHaveProperty('totalTranscriptions');
    
    // Check transcripts array
    expect(Array.isArray(transcriptData.transcripts)).toBe(true);
    if (transcriptData.transcripts.length > 0) {
      expect(transcriptData.transcripts[0]).toHaveProperty('text');
      expect(transcriptData.transcripts[0]).toHaveProperty('created_at');
    }
  });
});
