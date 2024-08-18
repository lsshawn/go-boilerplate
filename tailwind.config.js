/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./views/**/*.{html,js,templ,go}",
    "./views/common/**/*.{html,js,templ,go}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      fontFamily: {
        mono: ["Courier Prime", "monospace"],
      },
    },
  },
  plugins: [require("daisyui"), require("@tailwindcss/forms")],
  corePlugins: {
    preflight: true,
  },
  daisyui: {
    themes: ["light", "dark"],
  },
};
