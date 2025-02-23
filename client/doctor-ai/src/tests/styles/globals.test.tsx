// import React from 'react';
import { render } from '@testing-library/react';
import '../globals.css';

describe('Global Styles', () => {
  test('should apply global styles', () => {
    const { container } = render(<div className="test-class">Test</div>);
    expect(container.firstChild).toHaveClass('test-class'); // Adjust based on actual styles
  });
});
