import * as React from "react";

import { BodyWidget } from "./BodyWidget";
import { Application } from "./Application";

export default () => {
  var app = new Application();
  return <BodyWidget app={app} />;
};
