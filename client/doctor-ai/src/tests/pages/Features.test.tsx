// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import Features from '../../pages/Features';

describe('Features Page', () => {
  test('renders Features page', () => {
    render(<Features />);
    expect(screen.getByText(/features/i)).toBeInTheDocument(); // Adjust based on actual text
  });
});
