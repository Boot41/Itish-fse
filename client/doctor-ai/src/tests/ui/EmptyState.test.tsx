// import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import {EmptyState} from '../../components/ui/EmptyState';
// import {Spinner} from '../../components/ui/Spinner';

describe('EmptyState Component', () => {
  test('renders EmptyState', () => {
    render(<EmptyState title="No Data" description="There is no data available at this time." />);
    expect(screen.getByText(/no data available/i)).toBeInTheDocument(); // Adjust based on actual text
  });

  // test('renders Spinner', () => {
  //   render(<Spinner />);
  //   expect(screen.getByRole('progressbar')).toBeInTheDocument();
  // });
});
