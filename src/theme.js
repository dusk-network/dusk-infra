import red from "@material-ui/core/colors/red";
import { createMuiTheme } from "@material-ui/core/styles";

// A custom theme for this app
const theme = createMuiTheme({
  /* src/ui/theme/theme.json */

  palette: {
    tonalOffset: 0.2,
    background: { paper: "#fff", default: "#fafafa" },
    contrastThreshold: 3,
    grey: {
      "50": "#fafafa",
      "100": "#f5f5f5",
      "200": "#eeeeee",
      "300": "#e0e0e0",
      "400": "#bdbdbd",
      "500": "#9e9e9e",
      "600": "#757575",
      "700": "#616161",
      "800": "#424242",
      "900": "#212121",
      A700: "#616161",
      A100: "#d5d5d5",
      A400: "#303030",
      A200: "#aaaaaa"
    },
    text: {
      primary: "rgba(0, 0, 0, 0.87)",
      secondary: "rgba(0, 0, 0, 0.54)",
      disabled: "rgba(0, 0, 0, 0.38)",
      hint: "rgba(0, 0, 0, 0.38)"
    },
    divider: "rgba(0, 0, 0, 0.12)",
    secondary: {
      main: "#f37235",
      light: "rgb(245, 142, 93)",
      dark: "rgb(170, 79, 37)",
      contrastText: "#ffffff"
    },
    common: { black: "#000", white: "#fff" },
    error: {
      light: "#e57373",
      main: "#f44336",
      dark: "#d32f2f",
      contrastText: "#fff"
    },
    type: "light",
    action: {
      hoverOpacity: 0.08,
      hover: "rgba(0, 0, 0, 0.08)",
      selected: "rgba(0, 0, 0, 0.14)",
      disabledBackground: "rgba(0, 0, 0, 0.12)",
      disabled: "rgba(0, 0, 0, 0.26)",
      active: "rgba(0, 0, 0, 0.54)"
    },
    primary: {
      main: "#523b97",
      light: "rgb(116, 98, 171)",
      dark: "rgb(57, 41, 105)",
      contrastText: "#fff"
    }
  },
  themeName: "Dusk"

  // palette: {
  //   primary: {
  //     main: "#523b97"
  //   },
  //   secondary: {
  //     main: "#19857b"
  //   },
  //   error: {
  //     main: red.A400
  //   },
  //   background: {
  //     default: "#fff"
  //   }
  // }
});

export default theme;
