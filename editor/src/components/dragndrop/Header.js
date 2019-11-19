import * as React from "react";
import styled from "@emotion/styled";
import { Box, Button } from "@chakra-ui/core";

import _ from "lodash";

const Header = styled.div`
  display: flex;
  background: #414040;
  height: 4vh;
  z-index: 2000;
  border-bottom: solid 1px #e6e5e5;
  flex-grow: 0;
  flex-shrink: 0;
  color: white;
  font-family: Helvetica, Arial, sans-serif;
  padding: 10px;
  justify-content: space-between;
  align-items: center;
`;

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

export const HeaderSection = props => {
  return (
    <Header>
      <Box className="title" fontFamily="body">
        maya engine
      </Box>
      <Button
        paddingX="1rem"
        height="1.5rem"
        bg="#fcfcfc"
        color="#686767"
        size="xs"
        leftIcon="check"
        borderRadius="0.2rem"
        marginX="0.5rem"
        border="none"
        _hover={{ opacity: 0.8 }}
        _focus={{ outline: "none" }}
        onClick={() =>
          console.log(
            JsonToFBP(
              props.app
                .getDiagramEngine()
                .getModel()
                .serialize()
            )
          )
        }
      >
        Deploy
      </Button>
    </Header>
  );
};
