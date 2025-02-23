import { describe, test, expect } from 'vitest';
import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from '../App';

describe('Root Element', () => {
  test('root element exists', () => {
    // Create a mock root element
    const root = document.createElement('div');
    root.id = 'root';
    document.body.appendChild(root);

    // Render the App
    render(<App />, { container: root });

    // Check if root element exists and contains the app
    expect(document.getElementById('root')).toBeInTheDocument();
    expect(root.children.length).toBeGreaterThan(0);

    // Cleanup
    document.body.removeChild(root);
  });
});
