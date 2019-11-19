import * as React from "react";
import {
  Button,
  useDisclosure,
  Drawer,
  DrawerBody,
  DrawerOverlay,
  DrawerContent,
  DrawerHeader,
  DrawerFooter,
  Stack,
  Box,
  FormLabel,
  Input,
  DrawerCloseButton,
  Select,
  Textarea
} from "@chakra-ui/core";

export const DrawerExample = props => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const firstField = React.useRef();
  const btnRef = React.useRef();

  return (
    <>
      <Box ref={btnRef} onDoubleClick={onOpen}>
        {props.children}
      </Box>
      <Drawer
        isOpen={isOpen}
        placement="right"
        initialFocusRef={firstField}
        finalFocusRef={btnRef}
        onClose={onClose}
        right="12rem"
      >
        <DrawerOverlay opacity="0" />
        <DrawerContent
          height="94vh"
          top="6.5vh"
          right="12rem"
          boxShadow="none"
          borderLeft="solid 1px #e6e5e5"
          transition="opacity 0.5s ease-in-out"
        >
          <DrawerCloseButton />
          <DrawerHeader borderBottomWidth="1px">
            Create a new account
          </DrawerHeader>

          <DrawerBody>
            <Stack spacing="24px">
              <Box>
                <FormLabel htmlFor="username">Name</FormLabel>
                <Input
                  ref={firstField}
                  id="username"
                  placeholder="Please enter username"
                  width
                />
              </Box>

              <Box>
                <FormLabel htmlFor="owner">Select Owner</FormLabel>
                <Select id="owner" defaultValue="segun">
                  <option value="segun">Segun Adebayo</option>
                  <option value="kola">Kola Tioluwani</option>
                </Select>
              </Box>

              <Box>
                <FormLabel htmlFor="desc">Description</FormLabel>
                <Textarea id="desc" width="15rem" />
              </Box>
            </Stack>
          </DrawerBody>

          <DrawerFooter borderTopWidth="1px">
            <Button variant="outline" mr={3} onClick={onClose}>
              Cancel
            </Button>
            <Button variantColor="blue">Submit</Button>
          </DrawerFooter>
        </DrawerContent>
      </Drawer>
    </>
  );
};
