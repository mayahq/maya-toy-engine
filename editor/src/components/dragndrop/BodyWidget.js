import * as React from "react";
import * as _ from "lodash";
import { LeftBar } from "./LeftBar";
import { RightBar } from "./RightBar";
import { FlowTabs } from "./FlowTabs";
import { TrayItemWidget } from "./TrayItemWidget";
import { DefaultNodeModel } from "@projectstorm/react-diagrams";
import { CanvasWidget } from "@projectstorm/react-canvas-core";
import { DemoCanvasWidget } from "./DemoCanvasWidget";
import { IIPCustomNodeModel } from "./IIP/IIPNodeModel";
import { HeaderSection } from "./Header";
import library from "./library.json";
import styled from "@emotion/styled";

export const Body = styled.div`
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  min-height: 100%;
`;

export const Content = styled.div`
  display: flex;
  flex-grow: 1;
`;

export const Layer = styled.div`
  position: relative;
  flex-grow: 1;
`;

export const CenterSection = styled.div`
  position: relative;
  flex-grow: 1;
  flex-direction: row;
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
        <HeaderSection {...this.props} />
        <Content>
          <LeftBar>
            <TrayItemWidget
              model={{
                name: "core/iip"
              }}
              name="core/iip"
              color="rgb(0,192,255)"
            />
            {children}
          </LeftBar>
          <CenterSection>
            <FlowTabs />
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
          </CenterSection>
          <RightBar />
        </Content>
      </Body>
    );
  }
}
