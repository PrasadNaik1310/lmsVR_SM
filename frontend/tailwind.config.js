/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './index.html',
    './src/**/*.{js,jsx,ts,tsx,vue,svelte}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#0f766e',
      },
      borderRadius: {
        xl: '16px',
      },
    },
  },
  plugins: [require('@tailwindcss/forms')],
};
