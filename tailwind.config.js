/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{gohtml,html,js}"],
  theme: {
    fontFamily: {
      sans: ["Roboto Mono", "sans-serif"]
    },
    extend: {},
  },
  // darkMode: ['class', '[data-theme="dark"]'],
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        dark: {
          ...require("daisyui/src/theming/themes")["[data-theme=dark]"],
          "primary": "orange",
          "secondary": "white"
        },
      },
      {
        light: {
          ...require("daisyui/src/theming/themes")["[data-theme=light]"],
          "primary": "orange",
          "secondary": "lightgray",
        },
      },
    ],
    darkTheme: "dark",
    base: true,
    styled: true,
    utils: true,
  }
}

