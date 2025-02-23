import { useRefreshStore } from '../../stores/refreshStore';

describe('Refresh Store', () => {
  test('should have expected properties', () => {
    const refreshStore = useRefreshStore();
    expect(refreshStore).toHaveProperty('triggerPatientRefresh');
    expect(refreshStore).toHaveProperty('triggerTranscriptRefresh');
  });
});
