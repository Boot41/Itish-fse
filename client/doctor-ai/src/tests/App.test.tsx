import { describe, test, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from '../App';

describe('App Component', () => {
  test('renders App component', () => {
    render(<App />);
    expect(screen.getByText(/welcome to the app/i)).toBeInTheDocument(); // Adjust based on actual text
  });
});
