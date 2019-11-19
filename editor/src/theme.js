import { theme } from "@chakra-ui/core";

const customTheme = {
  ...theme,
  colors: {
    ...theme.colors,
    maya_light: {
      400: "-webkit-linear-gradient(180deg, #F6F6F6 0%, rgba(234, 234, 234, 0.97) 100%)",
      200: "#F4F4F4",
      100: "#FAFAFA",
      50: "#FCFCFC"
    },
    maya_dark: {
      500: "#686767",
      400: "linear-gradient(to right, red, rgba(255, 82, 66, 0) 52%)",
      300: "#3D3D3D",
      200: "#3F3E3E",
      100: "#474545",
      50: "#323334",
      25: "#ABA7A7",
      15: "#CAC6C6"
    },
    brown: {
      200: "#787373",
      100: "#D4D4D4",
      50: "#E5E3E3"
    }
  },
  fonts: {
    body: "Lato, sans-serif"
  },
  icons: {
    ...theme.icons
  }
};

export default customTheme;
