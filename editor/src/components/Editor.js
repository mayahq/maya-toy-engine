import React from "react";
import createEngine, {
  DefaultLinkModel,
  DiagramModel
} from "@projectstorm/react-diagrams";
import { Box, Button } from "@chakra-ui/core";

import { JSCustomNodeFactory } from "./custom/JSCustomNodeFactory";
import { JSCustomNodeModel } from "./custom/JSCustomNodeModel";
import { CanvasWidget } from "@projectstorm/react-canvas-core";
import _ from "lodash";

// create an instance of the engine with all the defaults

const JsonToFBP = model => {
  var processesObj = {};
  var connectionsObj = [];
  var ports = {};
  _.forEach(model.layers, item => {
    if (item.type === "diagram-nodes") {
      for (var key in item.models) {
        if (item.models.hasOwnProperty(key)) {
          var val = item.models[key];

          _.forEach(val.ports, item => {
            ports[item.id] = item;
          });
          processesObj[val.id] = { component: val.type };
        }
      }
    }
  });
  _.forEach(model.layers, item => {
    if (item.type === "diagram-links") {
      for (var key in item.models) {
        if (item.models.hasOwnProperty(key)) {
          var val = item.models[key];
          var src = {
            process: val.source,
            port: ports[val.sourcePort].label
          };
          var tgt = {
            process: val.target,
            port: ports[val.targetPort].label
          };
          connectionsObj.push({ src: src, tgt: tgt });
        }
      }
    }
  });
  var newArray = {
    properties: { name: "This is a test model" },
    processes: processesObj,
    connections: connectionsObj
  };

  return newArray;
};

export const Editor = () => {
  const engine = createEngine();

  engine.getNodeFactories().registerFactory(new JSCustomNodeFactory());

  // node 1
  const node1 = new JSCustomNodeModel({
    name: "core/console",
    color: "rgb(0,192,255)"
  });
  node1.setPosition(50, 50);

  // node 2
  const node2 = new JSCustomNodeModel({
    name: "core/exec",
    color: "rgb(0,192,255)"
  });
  node2.setPosition(200, 100);

  // node 3
  const node3 = new JSCustomNodeModel({
    name: "core/exec",
    color: "rgb(0,192,255)"
  });
  node3.setPosition(300, 200);

  // link them and add a label to the link
  const link1 = new DefaultLinkModel();
  link1.setSourcePort(node1.getPort("out"));
  link1.setTargetPort(node2.getPort("in"));

  const link2 = new DefaultLinkModel();
  link2.setSourcePort(node2.getPort("out"));
  link2.setTargetPort(node3.getPort("in"));

  const model = new DiagramModel();
  model.addAll(node1, node2, node3, link1, link2);
  engine.setModel(model);

  return (
    <Box>
      <Button
        variantColor="teal"
        variant="solid"
        size="xs"
        onClick={() => console.log(JsonToFBP(model.serialize()))}
      >
        Serialize
      </Button>
      <CanvasWidget className="srd-demo-canvas" engine={engine} />
    </Box>
  );
};
