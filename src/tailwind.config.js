/** @type {import( 'tailwindess').Config} */
module.exports = {
  content: ["./internal/view/**/*.templ}", "./**/*.templ"],
  safelist: [],
  plugins: [
    require("daisyui"),
    require("@tailwindcss/typography"),
    require('@tailwindcss/forms'),
  ],
  daisyui: {
    themes: ["light"]
  }
}
