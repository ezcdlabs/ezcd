import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],

  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        sm: "1100px",
      },
    },

    extend: {
      colors: {
        "cornflower-blue": {
          "50": "#eef5ff",
          "100": "#dae7ff",
          "200": "#bdd6ff",
          "300": "#90bdff",
          "400": "#6aa2ff",
          "500": "#3573fc",
          "600": "#1f52f1",
          "700": "#173dde",
          "800": "#1933b4",
          "900": "#1a308e",
          "950": "#151f56",
        },
        saffron: {
          "50": "#fefbe8",
          "100": "#fdf6c4",
          "200": "#fceb8c",
          "300": "#fad74a",
          "400": "#f7c421",
          "500": "#e7a90b",
          "600": "#c78207",
          "700": "#9f5c09",
          "800": "#834910",
          "900": "#703c13",
          "950": "#411e07",
        },

        white: {
          primary: "#BDBDBD",
          secondary: "#797979",
        },
        navy: {
          primary: "#071221",
        },
        cyan: {
          primary: "#53ACD3",
        },
        red: {
          primary: "#BF5852",
          secondary: "#4e0e0e",
        },
        yellow: {
          primary: "#E1E13A",
        },
        green: {
          primary: "#4EAC6F",
          secondary: "#3A7F5D",
        },

        grey: {
          primary: "#171717",
        },

        ansiBlack: "#000000",
        ansiBlue: "#2472c8",
        ansiBrightBlack: "#666666",
        ansiBrightBlue: "#3b8eea",
        ansiBrightCyan: "#29b8db",
        ansiBrightGreen: "#23d18b",
        ansiBrightMagenta: "#d670d6",
        ansiBrightRed: "#f14c4c",
        ansiBrightWhite: "#e5e5e5",
        ansiBrightYellow: "#f5f543",
        ansiCyan: "#11a8cd",
        ansiGreen: "#0dbc79",
        ansiMagenta: "#bc3fbc",
        ansiRed: "#cd3131",
        ansiWhite: "#e5e5e5",
        ansiYellow: "#e5e510",
      },
      container: {
        center: true,
      },
      // colors: {
      //   "azure-radiance": {
      //     "50": "#eff6ff",
      //     "100": "#dbebfe",
      //     "200": "#bfddfe",
      //     "300": "#92c8fe",
      //     "400": "#5fa9fb",
      //     "500": "#3282f7",
      //     "600": "#2467ec",
      //     "700": "#1c52d9",
      //     "800": "#1d43b0",
      //     "900": "#1d3d8b",
      //     "950": "#172754",
      //   },
      //   "black-pearl": {
      //     "1000": "hsl(217 76% 12% / 1);",
      //     "1050": "hsl(217 76% 8% / 1);",
      //     "1100": "hsl(217 73% 6% / 1);",
      //   },
      //   "deep-red": {
      //     "1000": "hsl(0 96% 12% / 1);",
      //     "1050": "hsl(0 96% 8% / 1);",
      //     "1100": "hsl(0 73% 6% / 1);",
      //   },
      // },
    },
  },
  plugins: [],
} satisfies Config;
