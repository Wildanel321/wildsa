/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: '#090d16',
        card: '#121827',
        border: '#1f293d',
        accent: '#3b82f6',
        emerald: '#10b981',
        amber: '#f59e0b',
        rose: '#f43f5e',
      },
    },
  },
  plugins: [],
}
