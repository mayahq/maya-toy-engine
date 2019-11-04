import * as SRD from "@projectstorm/react-diagrams";
import createEngine from "@projectstorm/react-diagrams";
import { InputType } from "@projectstorm/react-canvas-core";
import { IIPCustomNodeFactory } from "./IIP/IIPNodeFactory";

export class Application {
  constructor() {
    this.diagramEngine = createEngine();
    this.newModel = this.newModel.bind(this);
    this.getActiveDiagram = this.getActiveDiagram.bind(this);
    this.getDiagramEngine = this.getDiagramEngine.bind(this);

    this.newModel();
  }

  newModel() {
    this.activeModel = new SRD.DiagramModel();
    this.diagramEngine
      .getNodeFactories()
      .registerFactory(new IIPCustomNodeFactory());

    // Disable delete action
    // const actions = this.diagramEngine
    //   .getActionEventBus()
    //   .getActionsForType(InputType.KEY_DOWN);
    // this.diagramEngine.getActionEventBus().deregisterAction(actions[0]);

    //3-A) create a default node
    var node1 = new SRD.DefaultNodeModel("Node 1", "rgb(0,192,255)");
    let port = node1.addOutPort("Out");
    node1.setPosition(100, 100);

    //3-B) create another default node
    var node2 = new SRD.DefaultNodeModel("Node 2", "rgb(192,255,0)");
    let port2 = node2.addInPort("In");
    node2.setPosition(400, 100);

    // const node3 = new IIPCustomNodeModel({
    //   name: "core/iip",
    //   color: "rgb(0,192,255)"
    // });

    // node3.addInPort("In");

    // link the ports
    let link1 = port.link(port2);

    this.activeModel.addAll(node1, node2, link1);
    this.diagramEngine.setModel(this.activeModel);
  }

  getActiveDiagram() {
    return this.activeModel;
  }

  getDiagramEngine() {
    return this.diagramEngine;
  }
}
