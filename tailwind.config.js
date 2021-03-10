module.exports = {
  purge: [
    "app/ui/html/*.html",
    "app/ui/html/*.js"
  ],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        background: "#262b36"
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
