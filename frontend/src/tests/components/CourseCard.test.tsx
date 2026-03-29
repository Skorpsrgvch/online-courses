import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import CourseCard from '@/components/landing/CourseCard';

const mockCourse = {
  id: '1',
  title: 'Основы восстановления тазового дна',
  description: 'Подробный курс по упражнениям для восстановления после родов',
  price: 4990,
  image: '/course1.jpg',
};

describe('CourseCard', () => {
  it('renders course information correctly', () => {
    render(
      <MemoryRouter>
        <CourseCard {...mockCourse} />
      </MemoryRouter>
    );
    
    expect(screen.getByText('Основы восстановления тазового дна')).toBeInTheDocument();
    expect(screen.getByText('Подробный курс по упражнениям для восстановления после родов')).toBeInTheDocument();
    expect(screen.getByText('4 990 ₽')).toBeInTheDocument();
    expect(screen.getByRole('img')).toHaveAttribute('src', '/course1.jpg');
  });

  it('displays "Бесплатно" badge for free course', () => {
    render(
      <MemoryRouter>
        <CourseCard {...mockCourse} price={0} />
      </MemoryRouter>
    );
    
    expect(screen.getByText('Бесплатно')).toBeInTheDocument();
  });

  it('shows progress bar when progress prop is provided', () => {
    render(
      <MemoryRouter>
        <CourseCard {...mockCourse} progress={45} />
      </MemoryRouter>
    );
    
    const progressBar = screen.getByRole('progressbar');
    expect(progressBar).toHaveStyle('width: 45%');
    expect(screen.getByText('45% завершено')).toBeInTheDocument();
  });

  it('displays correct button text based on progress', () => {
    render(
      <MemoryRouter>
        <>
          <CourseCard {...mockCourse} />
          <CourseCard {...mockCourse} progress={45} />
        </>
      </MemoryRouter>
    );
    
    expect(screen.getAllByText('Подробнее').length).toBe(1);
    expect(screen.getAllByText('Продолжить обучение').length).toBe(1);
  });

  it('navigates to course preview when clicked for non-purchased course', () => {
    render(
      <MemoryRouter>
        <CourseCard {...mockCourse} />
      </MemoryRouter>
    );
    
    const link = screen.getByRole('link', { name: /подробнее/i });
    fireEvent.click(link);
    
    expect(window.location.pathname).toBe('/course/1/preview');
  });

  it('navigates to course when clicked for purchased course', () => {
    render(
      <MemoryRouter>
        <CourseCard {...mockCourse} progress={45} />
      </MemoryRouter>
    );
    
    const link = screen.getByRole('link', { name: /продолжить обучение/i });
    fireEvent.click(link);
    
    expect(window.location.pathname).toBe('/course/1');
  });
});