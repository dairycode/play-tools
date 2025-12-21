/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        dark: {
          bg: '#0a1929',
          card: 'rgba(255, 255, 255, 0.05)',
          'card-hover': 'rgba(255, 255, 255, 0.08)',
        },
        glass: {
          light: 'rgba(255, 255, 255, 0.1)',
          lighter: 'rgba(255, 255, 255, 0.15)',
          border: 'rgba(255, 255, 255, 0.2)',
          'border-light': 'rgba(255, 255, 255, 0.1)',
        },
        primary: {
          from: '#4fc3f7',
          to: '#2196f3',
        },
        success: {
          from: '#66bb6a',
          to: '#4caf50',
        },
        danger: {
          from: '#f44336',
          to: '#d32f2f',
        },
      },
      backdropBlur: {
        xs: '2px',
      },
      boxShadow: {
        'glass': '0 8px 32px rgba(0, 0, 0, 0.3)',
        'glow-blue': '0 4px 12px rgba(33, 150, 243, 0.3)',
        'glow-blue-lg': '0 6px 16px rgba(33, 150, 243, 0.4)',
        'glow-green': '0 4px 12px rgba(76, 175, 80, 0.3)',
        'glow-green-lg': '0 6px 16px rgba(76, 175, 80, 0.4)',
        'glow-red': '0 4px 12px rgba(244, 67, 54, 0.3)',
        'glow-red-lg': '0 6px 16px rgba(244, 67, 54, 0.4)',
      },
      animation: {
        'pulse-scale': 'pulseScale 2s ease-in-out infinite',
        'fade-in': 'fadeIn 0.3s ease-out',
      },
      keyframes: {
        pulseScale: {
          '0%, 100%': { transform: 'scale(1)', opacity: '1' },
          '50%': { transform: 'scale(1.05)', opacity: '0.8' },
        },
        fadeIn: {
          from: { opacity: '0', transform: 'translateY(10px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
};
