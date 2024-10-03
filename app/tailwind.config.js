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
        calm: 'var(--calm)',
      },
      animation: {
        wiggle: 'wiggle 20s infinite ease-in',
      },
      keyframes: {
        wiggle: {
          '0%': { transform: 'translate(0, 0)' },
          '10%': { transform: 'translate(-3px, 3px)' },
          '20%': { transform: 'translate(3px, -3px)' },
          '30%': { transform: 'translate(-3px, -3px)' },
          '40%': { transform: 'translate(3px, 3px)' },
          '50%': { transform: 'translate(0, 0)' },
          '60%': { transform: 'translate(3px, -3px)' },
          '70%': { transform: 'translate(-3px, 3px)' },
          '80%': { transform: 'translate(-3px, -3px)' },
          '90%': { transform: 'translate(3px, 3px)' },
          '100%': { transform: 'translate(0, 0)' },
        },
      },
    },
  },
  plugins: [],
}

