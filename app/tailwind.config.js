/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: false, 
  theme: {
    extend: {
      colors: {
        primary: 'var(--primary)',
        secondary: 'var(--secondary)',
        dark: 'var(--dark)',
        half_dark: 'var(--half-dark)',
        primary_light: 'var(--primary_light)',
        full_dark: 'var(--full-dark)',
        half_gray: 'var(--half-gray)',
      }
    },
  },
  plugins: [],
}

