import * as React from "react";
import styled from "@emotion/styled";
import { Box, Icon } from "@chakra-ui/core";

export const Tray = styled.div`
  width: 12rem;
  flex-grow: 0;
  flex-shrink: 0;
  background-color: white;
  border-left: solid 1px #e6e5e5;
  flex-grow: 0;
  z-index: 2000;
  flex-shrink: 0;
`;

export class RightBar extends React.Component {
  render() {
    return (
      <Tray>
        <Box height="2rem" fontFamily="body" />
        <Box
          display="flex"
          flexDirection="columns"
          justifyContent="space-between"
          alignItems="center"
          bg="#FCFCFC"
          borderTop="solid 1px #e6e5e5"
        >
          <Box display="flex" flexDirection="columns" alignItems="center">
            <Icon
              name="settings"
              paddingX="0.5rem"
              color="maya_dark.500"
              size="0.7rem"
            />
            <Box
              paddingBottom="0.4rem"
              paddingTop="0.4rem"
              fontFamily="body"
              fontWeight="400"
              fontSize="0.8rem"
              color="maya_dark.500"
            >
              Properties
            </Box>
          </Box>
          <Icon
            name="chevron-down"
            paddingX="0.5rem"
            color="maya_dark.500"
            size="1rem"
          />
        </Box>
        <Box height="5rem" fontFamily="body">
          <Box display="flex" flexDirection="columns" height="1.5rem">
            <Box
              display="flex"
              alignItems="center"
              width="3rem"
              bg="#FCFCFC"
              paddingX="1rem"
              fontSize="0.8rem"
              color="maya_dark.500"
              border="solid 1px #e6e5e5"
              borderLeft="none"
            >
              name
            </Box>
            <Box
              display="flex"
              alignItems="center"
              width="6rem"
              overflowWrap="break-word"
              paddingX="0.5rem"
              paddingBottom="0.2rem"
              fontSize="0.8rem"
              alignContent="center"
              border="solid 1px #e6e5e5"
              borderLeft="none"
              color="red.700"
            >
              "67483.8df9"
            </Box>
          </Box>
          <Box
            display="flex"
            flexDirection="columns"
            height="1.5rem"
            marginTop="-1px"
          >
            <Box
              display="flex"
              alignItems="center"
              width="3rem"
              bg="#FCFCFC"
              paddingX="1rem"
              fontSize="0.8rem"
              color="maya_dark.500"
              border="solid 1px #e6e5e5"
              borderLeft="none"
            >
              type
            </Box>
            <Box
              display="flex"
              alignItems="center"
              width="6rem"
              overflowWrap="break-word"
              paddingX="0.5rem"
              paddingBottom="0.2rem"
              fontSize="0.8rem"
              alignContent="center"
              border="solid 1px #e6e5e5"
              borderLeft="none"
              color="red.700"
            >
              "core/iip"
            </Box>
          </Box>
        </Box>
        <Box
          display="flex"
          flexDirection="columns"
          justifyContent="space-between"
          alignItems="center"
          bg="#FCFCFC"
          borderTop="solid 1px #e6e5e5"
        >
          <Box display="flex" flexDirection="columns" alignItems="center">
            <Icon
              name="info"
              paddingX="0.5rem"
              color="maya_dark.500"
              size="0.7rem"
            />
            <Box
              paddingBottom="0.4rem"
              paddingTop="0.4rem"
              fontFamily="body"
              fontWeight="400"
              fontSize="0.8rem"
              color="maya_dark.500"
            >
              Description
            </Box>
          </Box>
          <Icon
            name="chevron-down"
            paddingX="0.5rem"
            color="maya_dark.500"
            size="1rem"
          />
        </Box>
      </Tray>
    );
  }
}
