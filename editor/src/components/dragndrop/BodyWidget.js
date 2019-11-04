import * as React from "react";
import * as _ from "lodash";
import { TrayWidget } from "./TrayWidget";
import { TrayItemWidget } from "./TrayItemWidget";
import { DefaultNodeModel } from "@projectstorm/react-diagrams";
import { CanvasWidget } from "@projectstorm/react-canvas-core";
import { DemoCanvasWidget } from "./DemoCanvasWidget";
import { IIPCustomNodeModel } from "./IIP/IIPNodeModel";
import { Button } from "@chakra-ui/core";
import library from "./library.json";
import styled from "@emotion/styled";

export const Body = styled.div`
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  min-height: 100%;
`;

export const Header = styled.div`
  display: flex;
  background: rgb(30, 30, 30);
  flex-grow: 0;
  flex-shrink: 0;
  color: white;
  font-family: Helvetica, Arial, sans-serif;
  padding: 10px;
  align-items: center;
`;

export const Content = styled.div`
  display: flex;
  flex-grow: 1;
`;

export const Layer = styled.div`
  position: relative;
  flex-grow: 1;
`;

const ProcessJson = array => {
  let newArray = [];
  for (var key in array) {
    if (array.hasOwnProperty(key)) {
      var val = array[key];

      newArray.push(val);
    }
  }
  return newArray;
};

const JsonToFBP = model => {
  console.log("Library: ", model);
  var processesObj = {};
  var connectionsObj = [];
  var ports = {};
  var processes = {};
  _.forEach(model.layers, item => {
    if (item.type === "diagram-nodes") {
      for (var key in item.models) {
        if (item.models.hasOwnProperty(key)) {
          var val = item.models[key];

          _.forEach(val.ports, item => {
            ports[item.id] = item;
          });
          processes[val.id] = val;
          if (val.type !== "core/iip") {
            processesObj[val.id] = { component: val.name };
          }
        }
      }
    }
  });
  _.forEach(model.layers, item => {
    console.log("Layers", model.layers);
    if (item.type === "diagram-links") {
      for (var key in item.models) {
        if (item.models.hasOwnProperty(key)) {
          var val = item.models[key];
          let tgt;
          if (processes[val.source].type !== "core/iip") {
            var src = {
              process: val.source,
              port: ports[val.sourcePort].label
            };
            tgt = {
              process: val.target,
              port: ports[val.targetPort].label
            };
            connectionsObj.push({ src: src, tgt: tgt });
          } else if (processes[val.source].type === "core/iip") {
            var data = processes[val.source].data;
            tgt = {
              process: val.target,
              port: ports[val.targetPort].label
            };
            connectionsObj.push({ data: data, tgt: tgt });
          }
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

const GenerateNode = data => {
  var node = null;
  if (data.name === "core/iip") {
    node = new IIPCustomNodeModel({
      name: "core/iip",
      color: "rgb(0,192,255)",
      data: "4s"
    });
    node.addOutPort("out");
  } else {
    node = new DefaultNodeModel(data.name, "rgb(192,255,0)");
    _.forEach(data.outports, port => {
      node.addOutPort(port.name);
    });
    _.forEach(data.inports, port => {
      node.addInPort(port.name);
    });
  }

  return node;
};

export class BodyWidget extends React.Component {
  render() {
    const children = ProcessJson(library.entries).map((entry, i) => {
      return (
        <TrayItemWidget
          key={i}
          model={{
            name: entry.name,
            inports: entry.inports,
            outports: entry.outports
          }}
          name={entry.name}
          color="rgb(192,255,0)"
        />
      );
    });

    return (
      <Body>
        <Header>
          <div className="title">Maya Engine Demo</div>
          <Button
            variantColor="teal"
            size="xs"
            marginX="1rem"
            border="none"
            onClick={() =>
              console.log(
                JsonToFBP(
                  this.props.app
                    .getDiagramEngine()
                    .getModel()
                    .serialize()
                )
              )
            }
          >
            Serialize
          </Button>
        </Header>
        <Content>
          <TrayWidget>
            <TrayItemWidget
              model={{
                name: "core/iip"
              }}
              name="core/iip"
              color="rgb(0,192,255)"
            />
            {children}
          </TrayWidget>
          <Layer
            onDrop={event => {
              var data = JSON.parse(
                event.dataTransfer.getData("storm-diagram-node")
              );

              let node = GenerateNode(data);

              var point = this.props.app
                .getDiagramEngine()
                .getRelativeMousePoint(event);
              node.setPosition(point);
              this.props.app
                .getDiagramEngine()
                .getModel()
                .addNode(node);
              this.forceUpdate();
            }}
            onDragOver={event => {
              event.preventDefault();
            }}
          >
            <DemoCanvasWidget>
              <CanvasWidget engine={this.props.app.getDiagramEngine()} />
            </DemoCanvasWidget>
          </Layer>
        </Content>
      </Body>
    );
  }
}
