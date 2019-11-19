import React from "react";
import EditorDemo from "./components/dragndrop/index";
import customTheme from "./theme";
import { ThemeProvider } from "@chakra-ui/core";
import "./App.css";

import WebFont from "webfontloader";

WebFont.load({
  google: {
    families: ["Lato:300,400,700,900", "sans-serif"]
  }
});

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
