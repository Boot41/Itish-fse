// import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import ChatWindow from '../../components/chat/ChatWindow';

describe('ChatWindow Component', () => {
  test('renders ChatWindow', () => {
    render(<ChatWindow />);
    expect(screen.getByText(/chat/i)).toBeInTheDocument();
  });

  test('allows user to send a message', () => {
    render(<ChatWindow />);
    const input = screen.getByPlaceholderText(/type a message/i);
    const sendButton = screen.getByText(/send/i);

    fireEvent.change(input, { target: { value: 'Hello!' } });
    fireEvent.click(sendButton);

    expect(screen.getByText(/Hello!/i)).toBeInTheDocument();
  });
});
