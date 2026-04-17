import { useState, useCallback } from 'react';

interface ValidationRule {
  required?: string;
  pattern?: RegExp;
  patternMessage?: string;
  minLength?: number;
  minLengthMessage?: string;
  validate?: (value: any) => string | true;
}

export const useFormValidation = <T extends Record<string, any>>(initialValues: T) => {
  const [values, setValues] = useState<T>(initialValues);
  const [errors, setErrors] = useState<Partial<Record<keyof T, string>>>({});
  const [touched, setTouched] = useState<Partial<Record<keyof T, boolean>>>({});

  const validateField = useCallback((name: keyof T, value: any, rules: ValidationRule) => {
    let error = '';

    if (rules.required && !value) {
      error = rules.required;
    } else if (rules.minLength && typeof value === 'string' && value.length < rules.minLength) {
      error = rules.minLengthMessage || `Минимальная длина ${rules.minLength}`;
    } else if (rules.pattern && !rules.pattern.test(value)) {
      error = rules.patternMessage || 'Неверный формат';
    } else if (rules.validate) {
      const result = rules.validate(value);
      if (typeof result === 'string') error = result;
    }

    setErrors((prev) => ({ ...prev, [name]: error }));
    return !error;
  }, []);

  const handleChange = useCallback((name: keyof T, value: any, rules?: ValidationRule) => {
    setValues((prev) => ({ ...prev, [name]: value }));
    if (rules) {
      validateField(name, value, rules);
    }
  }, [validateField]);

  const handleBlur = useCallback((name: keyof T) => {
    setTouched((prev) => ({ ...prev, [name]: true }));
  }, []);

  const reset = useCallback(() => {
    setValues(initialValues);
    setErrors({});
    setTouched({});
  }, [initialValues]);

  return { values, errors, touched, handleChange, handleBlur, reset, setValues, setErrors };
};