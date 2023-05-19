import "./App.css";

import { AppShell, Container, Header, Text } from "@mantine/core";
import { Outlet } from "react-router-dom";

function Root() {
  return (
    <AppShell
      padding="sm"
      header={
        <Header height={60} p="lg">
          <Text ta="left">
            <a href="/">Todo App</a>
          </Text>
        </Header>
      }
    >
      <Container>
        <Outlet />
      </Container>
    </AppShell>
  );
}
export default Root;
