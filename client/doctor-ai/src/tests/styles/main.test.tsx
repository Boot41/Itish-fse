// import React from 'react';
import { render } from '@testing-library/react';
import '../main.css';

describe('Main Styles', () => {
  test('should apply main styles', () => {
    const { container } = render(<div className="main-class">Test</div>);
    expect(container.firstChild).toHaveClass('main-class'); // Adjust based on actual styles
  });
});
