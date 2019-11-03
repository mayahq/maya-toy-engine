import * as React from "react";
import { IIPCustomNodeModel } from "./IIPNodeModel";
import { IIPCustomNodeWidget } from "./IIPNodeWidget";
import { AbstractReactFactory } from "@projectstorm/react-canvas-core";

export class IIPCustomNodeFactory extends AbstractReactFactory {
  constructor() {
    super("js-custom-node");
  }

  generateModel(event) {
    console.log("Generating model", event);
    return new IIPCustomNodeModel();
  }

  generateReactWidget(event) {
    console.log("Generating widget", event);
    return <IIPCustomNodeWidget engine={this.engine} node={event.model} />;
  }
}

// import * as React from "react";
// import { IIPNodeModel } from "./IIPNodeModel";
// import { IIPNodeWidget } from "./IIPNodeWidget";
// import { AbstractReactFactory } from "@projectstorm/react-canvas-core";

// export class IIPNodeFactory extends AbstractReactFactory {
//   constructor() {
//     super("core/iip");
//   }

//   generateModel(event) {
//     console.log("Generating", event);
//     return new IIPNodeModel();
//   }

//   generateReactWidget(event) {
//     console.log("Generating", event.model);
//     return <IIPNodeWidget engine={this.engine} node={event.model} />;
//   }
// }
