/**
 * Валидаторы для форм.
 */

export interface ValidationResult {
  isValid: boolean;
  error: string;
}

/**
 * Валидация email
 */
export function validateEmail(email: string): ValidationResult {
  if (!email) return { isValid: false, error: 'Email обязателен' };
  const regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!regex.test(email)) return { isValid: false, error: 'Некорректный email' };
  return { isValid: true, error: '' };
}

/**
 * Валидация пароля (минимум 6 символов)
 */
export function validatePassword(password: string): ValidationResult {
  if (!password) return { isValid: false, error: 'Пароль обязателен' };
  if (password.length < 6) return { isValid: false, error: 'Минимум 6 символов' };
  return { isValid: true, error: '' };
}

/**
 * Валидация имени (минимум 2 символа)
 */
export function validateName(name: string): ValidationResult {
  if (!name) return { isValid: false, error: 'Имя обязательно' };
  if (name.length < 2) return { isValid: false, error: 'Минимум 2 символа' };
  return { isValid: true, error: '' };
}

/**
 * Валидация совпадения паролей
 */
export function validatePasswordMatch(password: string, confirmPassword: string): ValidationResult {
  if (!confirmPassword) return { isValid: false, error: 'Подтвердите пароль' };
  if (password !== confirmPassword) return { isValid: false, error: 'Пароли не совпадают' };
  return { isValid: true, error: '' };
}

/**
 * Валидация рейтинга (1-5)
 */
export function validateRating(rating: number): ValidationResult {
  if (rating < 1 || rating > 5) return { isValid: false, error: 'Рейтинг должен быть от 1 до 5' };
  return { isValid: true, error: '' };
}

/**
 * Валидация текста отзыва (минимум 10 символов)
 */
export function validateReviewText(text: string): ValidationResult {
  if (!text) return { isValid: false, error: 'Текст отзыва обязателен' };
  if (text.length < 10) return { isValid: false, error: 'Минимум 10 символов' };
  return { isValid: true, error: '' };
}
