import * as React from "react";
import styled from "@emotion/styled";

export const Container = styled.div`
  height: 100%;
  background-color: #fcfcfc;
  background-size: 40px 40px;
  display: flex;
  height: 100vh;
  > * {
    height: 100%;
    min-height: 100%;
    width: 100%;
  }
  background-image: linear-gradient(
      0deg,
      transparent 24%,
      ${p => p.color} 25%,
      ${p => p.color} 26%,
      transparent 27%,
      transparent 74%,
      ${p => p.color} 75%,
      ${p => p.color} 76%,
      transparent 77%,
      transparent
    ),
    linear-gradient(
      90deg,
      transparent 24%,
      ${p => p.color} 25%,
      ${p => p.color} 26%,
      transparent 27%,
      transparent 74%,
      ${p => p.color} 75%,
      ${p => p.color} 76%,
      transparent 77%,
      transparent
    );
`;

export class DemoCanvasWidget extends React.Component {
  render() {
    return (
      <Container
        background={this.props.background || "rgb(60, 60, 60)"}
        color={this.props.color || "#E2E2E2"}
      >
        {this.props.children}
      </Container>
    );
  }
}
