import 'bootstrap/dist/css/bootstrap.min.css';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App';

const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Не удалось найти корневой элемент #root в DOM');
}

createRoot(rootElement).render(
  <StrictMode>
    <App />
  </StrictMode>
);