/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        surface: '#F9F9F6',
        navy: {
          DEFAULT: '#1E2A38',
          light:   '#2C3E54',
          dark:    '#141D27',
        },
        terra: {
          DEFAULT: '#CC6649',
          light:   '#D97F65',
          dark:    '#A8503A',
        },
        'warm-gray': {
          50:  '#FAF9F7',
          100: '#F2F1EE',
          200: '#E8E7E4',
          300: '#D4D2CC',
          400: '#B8B5AE',
          500: '#8C8980',
          600: '#5A5750',
          700: '#3D3B36',
          800: '#2A2926',
          900: '#1A1917',
        }
      },
      fontFamily: {
        sans: ['"Inter"', 'system-ui', '-apple-system', 'sans-serif']
      },
      boxShadow: {
        soft:  '0 2px 12px 0 rgba(30,42,56,0.07)',
        card:  '0 4px 24px 0 rgba(30,42,56,0.08)',
        modal: '0 8px 40px 0 rgba(30,42,56,0.18)',
      },
      borderRadius: {
        '2xl': '1rem',
        '3xl': '1.5rem',
      },
      animation: {
        'spin-slow': 'spin 1.5s linear infinite',
      }
    }
  },
  plugins: [
    require('@tailwindcss/typography'),
  ]
}
