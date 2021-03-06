import * as React from "react";
import styled from "@emotion/styled";
import {
  Tabs,
  TabList,
  TabPanels,
  TabPanel,
  Tab,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Icon,
  Box
} from "@chakra-ui/core";

export const Tray = styled.div`
  width: 100%;
  margin-top: 0.5rem;
  background: "#FFFFFF";
  border-bottom: solid 1px #e6e5e5;
  flex-grow: 0;
  flex-shrink: 0;
`;

export const FlowTabs = () => {
  const CustomTab = React.forwardRef((props, ref) => {
    // `isSelected` is passed to all children of `TabList`.
    return (
      <Tab
        ref={ref}
        isSelected={props.isSelected}
        padding="none"
        background="none"
        _focus={{ outline: "none" }}
        fontFamily="body"
        fontSize="0.7rem"
        border="none"
        color="maya_dark.500"
        // borderTop={props.isSelected ? "solid 1px #e6e5e5" : "white"}
        //borderBottom={props.isSelected ? null : "solid 1px #e6e5e5"}
        {...props}
      >
        <Box
          zIndex="2"
          width="100%"
          paddingX="1rem"
          paddingTop="0.3rem"
          paddingBottom="0.3rem"
          marginBottom="-0.1rem"
          color="maya_dark.500"
          bg={props.isSelected ? "#FCFCFC" : null}
          borderRadius="0.2rem 0.2rem 0 0"
          borderTop={props.isSelected ? "solid 1px #e6e5e5" : null}
          borderLeft={props.isSelected ? "solid 1px #e6e5e5" : null}
          borderRight={props.isSelected ? "solid 1px #e6e5e5" : null}
        >
          {props.children}
        </Box>
      </Tab>
    );
  });
  return (
    <Tray>
      <Tabs variant="enclosed" color="maya_dark.500">
        <TabList color="maya_dark.500">
          <CustomTab>Flow 1</CustomTab>
          <CustomTab>Flow 2</CustomTab>
        </TabList>
        <TabPanels>
          <TabPanel
            bg="#FCFCFC"
            height="1.8rem"
            borderTop="solid 1px #e6e5e5"
            boxShadow="0px 2px 4px rgba(0,0,0,0.05)"
          >
            <Box
              display="flex"
              flexDirection="columns"
              justifyItems="start"
              alignItems="center"
            >
              <Breadcrumb
                spacing="0.5rem"
                fontFamily="body"
                fontSize="0.8rem"
                color="maya_dark.500"
                marginTop="-0.5rem"
                separator={<Icon color="gray.300" name="chevron-right" />}
              >
                <BreadcrumbItem>
                  <BreadcrumbLink>Main</BreadcrumbLink>
                </BreadcrumbItem>

                <BreadcrumbItem>
                  <BreadcrumbLink>SubFlow</BreadcrumbLink>
                </BreadcrumbItem>

                <BreadcrumbItem isCurrentPage>
                  <BreadcrumbLink>SubSubFlow</BreadcrumbLink>
                </BreadcrumbItem>
              </Breadcrumb>
            </Box>
          </TabPanel>
          <TabPanel
            bg="#FCFCFC"
            height="1.8rem"
            borderTop="solid 1px #e6e5e5"
            boxShadow="0px 2px 4px rgba(0,0,0,0.05)"
          >
            <Box
              display="flex"
              flexDirection="columns"
              justifyItems="start"
              alignItems="center"
            >
              <Breadcrumb
                spacing="0.5rem"
                fontFamily="body"
                fontSize="0.8rem"
                color="maya_dark.500"
                marginTop="-0.5rem"
                separator={<Icon color="gray.300" name="chevron-right" />}
              >
                <BreadcrumbItem>
                  <BreadcrumbLink>Main2</BreadcrumbLink>
                </BreadcrumbItem>

                <BreadcrumbItem>
                  <BreadcrumbLink>SubFlow2</BreadcrumbLink>
                </BreadcrumbItem>

                <BreadcrumbItem isCurrentPage>
                  <BreadcrumbLink>SubSubFlow2</BreadcrumbLink>
                </BreadcrumbItem>
              </Breadcrumb>
            </Box>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </Tray>
  );
};
