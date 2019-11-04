import React from "react";
import EditorDemo from "./components/dragndrop/index";

import { theme, ThemeProvider } from "@chakra-ui/core";

import "./App.css";
// Let's say you want to add custom colors
const customTheme = {
  ...theme,
  colors: {
    ...theme.colors,
    brand: {
      900: "#1a365d",
      800: "#153e75",
      700: "#2a69ac"
    }
  }
};

function App() {
  return (
    <div className="App">
      <ThemeProvider theme={customTheme}>
        <EditorDemo />
      </ThemeProvider>
    </div>
  );
}

export default App;
