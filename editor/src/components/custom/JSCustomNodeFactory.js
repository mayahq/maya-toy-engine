import * as React from "react";
import { JSCustomNodeModel } from "./JSCustomNodeModel";
import { JSCustomNodeWidget } from "./JSCustomNodeWidget";
import { AbstractReactFactory } from "@projectstorm/react-canvas-core";

export class JSCustomNodeFactory extends AbstractReactFactory {
  constructor() {
    super("js-custom-node");
  }

  generateModel(event) {
    console.log("Generating model", event);
    return new JSCustomNodeModel();
  }

  generateReactWidget(event) {
    console.log("Generating widget", event);
    return <JSCustomNodeWidget engine={this.engine} node={event.model} />;
  }
}
