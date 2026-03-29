import { render, screen, fireEvent } from '@testing-library/react';
import CookieBanner from '@/components/layout/CookieBanner';
import { act } from 'react-dom/test-utils';

jest.mock('@/lib/dompurify', () => ({
  sanitize: jest.fn().mockImplementation((str) => str),
}));

describe('CookieBanner', () => {
  beforeEach(() => {
    localStorage.clear();
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('does not show when preferences are saved', () => {
    localStorage.setItem('cookiePreferences', JSON.stringify({
      necessary: true,
      analytics: true,
      marketing: true
    }));
    
    render(<CookieBanner />);
    
    expect(screen.queryByText('Мы заботимся о вашей конфиденциальности')).not.toBeInTheDocument();
  });

  it('shows when no preferences are saved', () => {
    render(<CookieBanner />);
    
    expect(screen.getByText('Мы заботимся о вашей конфиденциальности')).toBeInTheDocument();
  });

  it('saves preferences when "Сохранить настройки" is clicked', () => {
    render(<CookieBanner />);
    
    const analyticsCheckbox = screen.getByLabelText('Разрешить');
    fireEvent.click(analyticsCheckbox);
    
    const saveButton = screen.getByText('Сохранить настройки');
    fireEvent.click(saveButton);
    
    expect(localStorage.getItem('cookiePreferences')).toBe(
      JSON.stringify({ 
        necessary: true,
        analytics: true,
        marketing: false 
      })
    );
  });

  it('accepts all cookies when "Принять все" is clicked', () => {
    render(<CookieBanner />);
    
    const acceptAllButton = screen.getByText('Принять все');
    fireEvent.click(acceptAllButton);
    
    expect(localStorage.getItem('cookiePreferences')).toBe(
      JSON.stringify({ 
        necessary: true,
        analytics: true,
        marketing: true 
      })
    );
  });

  it('initializes analytics when analytics is enabled', () => {
    const initAnalyticsMock = jest.fn();
    jest.spyOn(console, 'log').mockImplementation((msg) => {
      if (msg === 'Analytics initialized') initAnalyticsMock();
    });
    
    render(<CookieBanner />);
    
    const analyticsCheckbox = screen.getByLabelText('Разрешить');
    fireEvent.click(analyticsCheckbox);
    
    const saveButton = screen.getByText('Сохранить настройки');
    fireEvent.click(saveButton);
    
    expect(initAnalyticsMock).toHaveBeenCalled();
  });
});