/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "selector",
  content: ["./templates/**/*.{html,templ,js}"],
  theme: {
    extend: {},
  },
  plugins: [
	require('@tailwindcss/typography'),
  ],
}

