import { DefaultPortModel, NodeModel } from "@projectstorm/react-diagrams";

/**
 * Example of a custom model using pure javascript
 */
export class IIPCustomNodeModel extends NodeModel {
  constructor(options = {}) {
    super({
      ...options,
      type: "js-custom-node"
    });
    this.color = options.color || { options: "red" };

    // setup an in and out port
    this.addPort(
      new DefaultPortModel({
        in: true,
        name: "in"
      })
    );
    this.addPort(
      new DefaultPortModel({
        in: false,
        name: "out"
      })
    );
  }

  serialize() {
    return {
      ...super.serialize(),
      color: this.options.color
    };
  }

  deserialize(ob, engine) {
    super.deserialize(ob, engine);
    this.color = ob.color;
  }
}

// import {
//   DefaultPortModel,
//   NodeModel,
//   PortModelAlignment
// } from "@projectstorm/react-diagrams";
// import * as _ from "lodash";

// export class IIPNodeModel extends NodeModel {
//   constructor(options = {}, color) {
//     if (typeof options === "string") {
//       options = {
//         name: options,
//         color: color
//       };
//     }
//     super({
//       type: "default",
//       name: "Untitled",
//       color: "rgb(0,192,255)",
//       ...options
//     });
//     this.portsOut = [];
//     this.portsIn = [];
//   }

//   doClone(lookupTable, clone) {
//     clone.portsIn = [];
//     clone.portsOut = [];
//     super.doClone(lookupTable, clone);
//   }

//   removePort(port) {
//     super.removePort(port);
//     if (port.getOptions().in) {
//       this.portsIn.splice(this.portsIn.indexOf(port));
//     } else {
//       this.portsOut.splice(this.portsOut.indexOf(port));
//     }
//   }

//   addPort(port) {
//     super.addPort(port);
//     if (port.getOptions().in) {
//       if (this.portsIn.indexOf(port) === -1) {
//         this.portsIn.push(port);
//       }
//     } else {
//       if (this.portsOut.indexOf(port) === -1) {
//         this.portsOut.push(port);
//       }
//     }
//     return port;
//   }

//   addInPort(label, after = true) {
//     const p = new DefaultPortModel({
//       in: true,
//       name: label,
//       label: label,
//       alignment: PortModelAlignment.LEFT
//     });
//     if (!after) {
//       this.portsIn.splice(0, 0, p);
//     }
//     return this.addPort(p);
//   }

//   addOutPort(label, after = true) {
//     const p = new DefaultPortModel({
//       in: false,
//       name: label,
//       label: label,
//       alignment: PortModelAlignment.RIGHT
//     });
//     if (!after) {
//       this.portsOut.splice(0, 0, p);
//     }
//     return this.addPort(p);
//   }

//   deserialize(event) {
//     super.deserialize(event);
//     this.options.name = event.data.name;
//     this.options.color = event.data.color;
//     this.portsIn = _.map(event.data.portsInOrder, id => {
//       return this.getPortFromID(id);
//     });
//     this.portsOut = _.map(event.data.portsOutOrder, id => {
//       return this.getPortFromID(id);
//     });
//   }

//   serialize() {
//     return {
//       ...super.serialize(),
//       name: this.options.name,
//       color: this.options.color,
//       portsInOrder: _.map(this.portsIn, port => {
//         return port.getID();
//       }),
//       portsOutOrder: _.map(this.portsOut, port => {
//         return port.getID();
//       })
//     };
//   }

//   getInPorts() {
//     return this.portsIn;
//   }

//   getOutPorts() {
//     return this.portsOut;
//   }
// }
