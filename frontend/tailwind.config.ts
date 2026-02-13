import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './lib/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        },
        category: {
          entertainment: '#FF6B6B',
          productivity: '#4ECDC4',
          cloud: '#45B7D1',
          ai: '#96CEB4',
          shopping: '#FFEAA7',
          etc: '#DFE6E9',
        },
      },
    },
  },
  plugins: [],
};

export default config;
