// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import MainLayout from '../../components/layout/MainLayout';

describe('MainLayout Component', () => {
  test('renders MainLayout', () => {
    render(<MainLayout><div>Test Content</div></MainLayout>);
    expect(screen.getByText(/header text/i)).toBeInTheDocument(); // Adjust based on actual header text
    expect(screen.getByText(/footer text/i)).toBeInTheDocument(); // Adjust based on actual footer text
  });
});
