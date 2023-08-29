/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.{gohtml,html,js}"],
  theme: {
    fontFamily: {
      sans: ["Roboto Mono", "sans-serif"]
    },
    extend: {},
  },
  plugins: [require("daisyui")],
  //darkMode: "class",
  darkMode: ['class', '[data-theme="dark"]'],
  daisyui: {
    themes: [
      {
        dark: {
          ...require("daisyui/src/theming/themes")["[data-theme=dark]"],
          "primary": "red",
          "secondary": "orange"
        },
      },
      "light"
    ],
    darkTheme: 'dark',
    base: true,
    styled: true,
    utils: true,
  }
}

