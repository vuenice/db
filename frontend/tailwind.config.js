/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          900: '#004d40', // Deep teal, matching InstaWP
        },
        grayCust: {
          50: '#f9fafb',
          160: '#e5e7eb',
          500: '#6b7280',
          550: '#4b5563',
          640: '#374151',
          700: '#1f2937',
          800: '#111827',
          980: '#030712'
        },
        secondary: {
          800: '#2563eb'
        },
        emerald: {
          50: '#ecfdf5'
        }
      }
    },
  },
  plugins: [],
}
