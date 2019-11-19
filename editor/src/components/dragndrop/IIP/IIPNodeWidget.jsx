import * as React from "react";

import { Box, Icon, PseudoBox, Tooltip } from "@chakra-ui/core";
import { PortWidget } from "@projectstorm/react-diagrams";
import { DrawerExample } from "../PropsDrawer";
import styled from "@emotion/styled";
import * as _ from "lodash";
import "./node.css";

export const Node = styled.div`
  display: flex;
  flex-direction: columns;
  border-radius: 5px;
  font-family: sans-serif;
  color: white;
  height: 2rem;
  overflow: visible;
  font-size: 11px;
`;

const PortsContainer = styled.div`
  flex-grow: 1;
  display: flex;
  flex-direction: columns;
  align-items: center;
`;

const NodeButton = props => {
  return (
    <DrawerExample>
      <Box
        display="flex"
        flex-direction="columns"
        alignItems="center"
        bg={props.node.isSelected() ? "#FCFCFC" : "white"}
        borderRadius="0.6rem"
        border="solid 1.5px"
        borderColor={props.node.getOptions().color}
      >
        <Box
          display="flex"
          bg={props.node.getOptions().color}
          width="1rem"
          height="2rem"
          borderRadius="0.4rem 0 0 0.4rem"
          paddingX="0.6rem"
          alignItems="center"
        >
          <Icon name="view" color="white" size="1rem" />
        </Box>
        <Box
          paddingX="1rem"
          paddingY="0.5rem"
          fontFamily="body"
          fontWeight="700"
          fontSize="0.8rem"
          color="maya_dark.500"
        >
          {props.node.getOptions().name}
        </Box>
      </Box>
    </DrawerExample>
  );
};
export class IIPCustomNodeWidget extends React.Component {
  generatePortLeft = port => {
    return (
      <PortWidget engine={this.props.engine} port={port} key={port.getID()}>
        <Tooltip hasArrow label={port.getName()} placement="left">
          <PseudoBox
            width="10px"
            height="10px"
            borderRadius="5px"
            zIndex="9999"
            marginRight="-5px"
            cursor="pointer"
            bg={this.props.node.getOptions().color}
            _hover={{
              boxShadow: "inset 0 0 10px 10px rgba(255, 255, 255, 0.4)"
            }}
          />
        </Tooltip>
      </PortWidget>
    );
  };
  generatePortRight = port => {
    return (
      <PortWidget engine={this.props.engine} port={port} key={port.getID()}>
        <Tooltip hasArrow label={port.getName()} placement="right">
          <PseudoBox
            width="10px"
            height="10px"
            borderRadius="5px"
            zIndex="9999"
            marginLeft="-6px"
            bg={this.props.node.getOptions().color}
            cursor="pointer"
            _hover={{
              boxShadow: "inset 0 0 10px 10px rgba(255, 255, 255, 0.4)"
            }}
          />
        </Tooltip>
      </PortWidget>
    );
  };
  render() {
    return (
      <Node
        data-default-node-name={this.props.node.getOptions().name}
        selected={this.props.node.isSelected()}
        background={this.props.node.getOptions().color}
      >
        <PortsContainer>
          {_.map(this.props.node.getInPorts(), this.generatePortLeft)}
        </PortsContainer>
        <NodeButton {...this.props} />
        <PortsContainer>
          {_.map(this.props.node.getOutPorts(), this.generatePortRight)}
        </PortsContainer>
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
