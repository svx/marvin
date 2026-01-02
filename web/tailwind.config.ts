import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#f5f5ff',
          100: '#ebebff',
          200: '#d6d6ff',
          300: '#b8b8ff',
          400: '#9999ff',
          500: '#635BFF',
          600: '#5851e8',
          700: '#4d47d1',
          800: '#423dba',
          900: '#3733a3',
        },
        success: {
          50: '#e6fff0',
          100: '#b3ffe0',
          200: '#80ffd1',
          300: '#4dffc1',
          400: '#1affb2',
          500: '#00D924',
          600: '#00c020',
          700: '#00a71c',
          800: '#008e18',
          900: '#007514',
        },
        warning: {
          50: '#fff8e6',
          100: '#ffecb3',
          200: '#ffe080',
          300: '#ffd44d',
          400: '#ffc81a',
          500: '#FFA500',
          600: '#e69500',
          700: '#cc8500',
          800: '#b37500',
          900: '#996500',
        },
        error: {
          50: '#ffe6ec',
          100: '#ffb3c7',
          200: '#ff80a2',
          300: '#ff4d7d',
          400: '#ff1a58',
          500: '#DF1B41',
          600: '#c8183a',
          700: '#b11533',
          800: '#9a122c',
          900: '#830f25',
        },
        background: '#ffffff',
        foreground: '#0a2540',
        border: '#e3e8ee',
        muted: '#697386',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['Fira Code', 'monospace'],
      },
      borderRadius: {
        stripe: '6px',
      },
      boxShadow: {
        stripe: '0 2px 4px rgba(0,0,0,0.06)',
        'stripe-lg': '0 4px 12px rgba(0,0,0,0.08)',
      },
    },
  },
  plugins: [],
}

export default config
