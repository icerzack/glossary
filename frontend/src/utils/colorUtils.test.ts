import { stringToColor } from './colorUtils';

describe('stringToColor', () => {
  it('should generate consistent colors for the same string', () => {
    const color1 = stringToColor('Test Category');
    const color2 = stringToColor('Test Category');
    expect(color1).toBe(color2);
  });

  it('should generate different colors for different strings', () => {
    const color1 = stringToColor('Category A');
    const color2 = stringToColor('Category B');
    expect(color1).not.toBe(color2);
  });

  it('should return HSL format', () => {
    const color = stringToColor('Test');
    expect(color).toMatch(/^hsl\(\d+,\s*\d+%,\s*\d+%\)$/);
  });

  it('should handle empty string', () => {
    const color = stringToColor('');
    expect(color).toMatch(/^hsl\(\d+,\s*\d+%,\s*\d+%\)$/);
  });

  it('should handle special characters', () => {
    const color = stringToColor('Test!@#$%^&*()');
    expect(color).toMatch(/^hsl\(\d+,\s*\d+%,\s*\d+%\)$/);
  });
});

