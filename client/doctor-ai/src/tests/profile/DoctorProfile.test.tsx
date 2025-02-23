// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import DoctorProfile from '../../components/profile/DoctorProfile';

describe('DoctorProfile Component', () => {
  test('renders DoctorProfile', () => {
    render(<DoctorProfile />);
    expect(screen.getByText(/doctor profile/i)).toBeInTheDocument(); // Adjust based on actual text
  });
});
