import * as React from "react";
import { Icon, Box } from "@chakra-ui/core";
import styled from "@emotion/styled";

export const Tray = styled.div`
  display: flex;
  flex-direction: columns;
  align-items: center;
  color: "#645f5f";
  background-color: "white";
  margin: 0.5rem 0.5rem;
  border: solid 1.5px ${p => p.color};
  border-radius: 0.5rem;
  cursor: pointer;
`;

export class TrayItemWidget extends React.Component {
  render() {
    return (
      <Tray
        color={this.props.color}
        draggable={true}
        onDragStart={event => {
          event.dataTransfer.setData(
            "storm-diagram-node",
            JSON.stringify(this.props.model)
          );
        }}
      >
        <Box
          bg={this.props.color}
          width="1rem"
          height="100%"
          borderRadius="0.4rem 0 0 0.4rem"
          paddingX="0.5rem"
          paddingTop="0.2rem"
          paddingBottom="0.4rem"
        >
          <Icon name="view" color="white" marginTop="0" />
        </Box>
        <Box
          paddingX="0.5rem"
          fontFamily="body"
          fontWeight="700"
          fontSize="0.8rem"
          color="maya_dark.500"
        >
          {this.props.name}
        </Box>
      </Tray>
    );
  }
}
