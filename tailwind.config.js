/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{html,templ,js}"],
  theme: {
    extend: {},
  },
  plugins: [
	require('@tailwindcss/typography'),
  ],
  safelist: [
	{
		pattern: /bg-blue-+/,
	}
  ], 
}

