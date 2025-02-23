// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import About from '../../pages/About';

describe('About Page', () => {
  test('renders About page', () => {
    render(<About />);
    expect(screen.getByText(/about us/i)).toBeInTheDocument(); // Adjust based on actual text
  });
});
