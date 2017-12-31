import React from 'react';
import ReactDOM from 'react-dom';
import Chat from './';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<Chat />, div);
});
