import * as React from "react";
import { DefaultPortLabel } from "@projectstorm/react-diagrams";
import styled from "@emotion/styled";
import * as _ from "lodash";
import "./node.css";

export const Node = styled.div`
  background-color: ${p => p.background};
  border-radius: 5px;
  font-family: sans-serif;
  color: white;
  border: solid 2px black;
  overflow: visible;
  font-size: 11px;
  border: solid 2px ${p => (p.selected ? "rgb(0,192,255)" : "black")};
`;

const Ports = styled.div`
  display: flex;
  color: white;
  background-image: linear-gradient(rgba(0, 0, 0, 0.1), rgba(0, 0, 0, 0.2));
`;

const PortsContainer = styled.div`
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  &:first-of-type {
    margin-right: 10px;
  }
  &:only-child {
    margin-right: 0px;
  }
`;

const Title = styled.div`
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  white-space: nowrap;
  justify-items: center;
`;

const TitleName = styled.div`
  flex-grow: 1;
  padding: 5px 5px;
`;

export class IIPCustomNodeWidget extends React.Component {
  generatePort = port => {
    return (
      <DefaultPortLabel
        engine={this.props.engine}
        port={port}
        key={port.getID()}
      />
    );
  };
  render() {
    return (
      <Node
        data-default-node-name={this.props.node.getOptions().name}
        selected={this.props.node.isSelected()}
        background={this.props.node.getOptions().color}
      >
        <Title>
          <TitleName>{this.props.node.getOptions().name}</TitleName>
        </Title>

        <Ports>
          <PortsContainer>
            {_.map(this.props.node.getInPorts(), this.generatePort)}
          </PortsContainer>
          <PortsContainer>
            {_.map(this.props.node.getOutPorts(), this.generatePort)}
          </PortsContainer>
        </Ports>
      </Node>
    );
  }
}

// import * as React from "react";
// import * as _ from "lodash";
// import "./node.css";
// import { DefaultPortLabel } from "@projectstorm/react-diagrams";
// import styled from "@emotion/styled";

// export const Node = styled.div`
//   background-color: ${p => p.background};
//   border-radius: 5px;
//   font-family: sans-serif;
//   color: white;
//   border: solid 2px black;
//   overflow: visible;
//   font-size: 11px;
//   border: solid 2px ${p => (p.selected ? "rgb(0,192,255)" : "blue")};
// `;

// export class IIPCustomNodeWidget extends React.Component {
//   constructor(props) {
//     super(props);
//     console.log("IIP Widget loaded");
//     this.generatePort = this.generatePort.bind(this);
//   }
//   generatePort = port => {
//     return (
//       <DefaultPortLabel
//         engine={this.props.engine}
//         port={port}
//         key={port.getID()}
//       />
//     );
//   };
//   render() {
//     return (
//       <Node
//         data-default-node-name={this.props.node.getOptions().name}
//         selected={this.props.node.isSelected()}
//         background={this.props.node.getOptions().color}
//       >
//         <Title>
//           <TitleName>{this.props.node.getOptions().name}</TitleName>
//         </Title>

//         <Ports>
//           <PortsContainer>
//             {_.map(this.props.node.getInPorts(), this.generatePort)}
//           </PortsContainer>
//           <PortsContainer>
//             {_.map(this.props.node.getOutPorts(), this.generatePort)}
//           </PortsContainer>
//         </Ports>
//       </Node>
//     );
//   }
// }
