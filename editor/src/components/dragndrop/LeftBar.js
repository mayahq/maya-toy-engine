import * as React from "react";
import styled from "@emotion/styled";
import { Box, Icon, Input } from "@chakra-ui/core";

export const Tray = styled.div`
  min-width: 10rem;
  background: "#FFFFFF";
  border-right: solid 1px #e6e5e5;
  flex-grow: 0;
  flex-shrink: 0;
`;

export class LeftBar extends React.Component {
  render() {
    return (
      <Tray>
        <Box
          display="flex"
          flexDirection="columns"
          borderBottom="solid 1px"
          borderColor="#e6e5e5"
        >
          <Icon
            name="search"
            color="maya_dark.500"
            padding="0.6rem"
            size="0.8rem"
          />
          <Input
            placeholder="filter nodes..."
            fontSize="0.8rem"
            width="6rem"
            paddingX="0.2rem"
            paddingY="0.5rem"
            fontFamily="body"
            border="none"
            outline="none"
            height="1rem"
            color="maya_dark.500"
            _focus={{
              outline: "none"
            }}
          ></Input>
        </Box>
        <Box
          display="flex"
          flexDirection="columns"
          alignItems="center"
          bg="#FCFCFC"
          boxShadow="0px 2px 4px rgba(0, 0, 0, 0.05)"
        >
          <Icon name="chevron-down" paddingX="0.5rem" color="maya_dark.500" />
          <Box
            paddingBottom="0.5rem"
            paddingTop="0.4rem"
            fontFamily="body"
            fontWeight="400"
            fontSize="0.8rem"
            color="maya_dark.500"
          >
            core
          </Box>
        </Box>
        {this.props.children}
      </Tray>
    );
  }
}
