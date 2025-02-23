// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Patients from '../../components/dashboard/Patients';

describe('Patients Component', () => {
  test('renders Patients', () => {
    render(<Patients />);
    expect(screen.getByText(/patients/i)).toBeInTheDocument();
  });

  test('displays patient list', () => {
    render(<Patients />);
    // Assuming there is a mock patient data or a way to test the list
    expect(screen.getByText(/john doe/i)).toBeInTheDocument(); // Adjust based on actual patient names
  });
});
