import defaultTheme from 'tailwindcss/defaultTheme';

export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx,css}'],
  darkMode: 'class', // Use class-based dark mode toggling
  theme: {
    extend: {
      fontFamily: {
        inter: ['Inter', ...defaultTheme.fontFamily.sans],
      },
      colors: {
        primary: {
          DEFAULT: '#3ECF8E',
          50: '#E6FFF4',
          100: '#B3FFE4',
          200: '#80FFD4',
          300: '#4DFFC4',
          400: '#1AFFB4',
          500: '#00E69D',
          600: '#00B37A',
          700: '#008057',
          800: '#004D34',
          900: '#001A11',
        },
        dark: {
          DEFAULT: '#1C1C1C',
          50: '#4D4D4D',
          100: '#404040',
          200: '#333333',
          300: '#262626',
          400: '#1A1A1A',
          500: '#0D0D0D',
          600: '#000000',
        },
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-out forwards',
        'fade-in-up': 'fadeInUp 0.5s ease-out forwards',
        'slide-in': 'slideIn 0.5s ease-out forwards',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'bounce-slow': 'bounce 2s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        fadeInUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideIn: {
          '0%': { transform: 'translateX(-100%)' },
          '100%': { transform: 'translateX(0)' },
        },
      },
      boxShadow: {
        'subtle-primary': '0 10px 25px rgba(0, 255, 0, 0.1)',
      },
    },
    extend: {
      colors: {
        primary: {
          DEFAULT: 'rgb(var(--color-primary) / <alpha-value>)',
          50: 'rgb(var(--color-primary-50) / <alpha-value>)',
          100: 'rgb(var(--color-primary-100) / <alpha-value>)',
          // Add more scales if you've defined them
          // 200, 300, 400, etc.
        },
        darks: { // using "darks" to avoid conflict with Tailwind's dark mode
          DEFAULT: 'rgb(var(--color-darks) / <alpha-value>)',
          // You can also define additional scales if needed.
        },
      },
    },
  },
  plugins: [
    function({ addUtilities }) {
      addUtilities({
        '.scrollbar-hide': {
          /* Hide scrollbar for Chrome, Safari and Opera */
          '&::-webkit-scrollbar': {
            display: 'none'
          },
          /* Hide scrollbar for IE, Edge and Firefox */
          '-ms-overflow-style': 'none',  /* IE and Edge */
          'scrollbar-width': 'none'  /* Firefox */
        }
      })
    },
    require('framer-motion/plugin')
  ]
};
