import React, { FunctionComponent } from "react";
import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import Container from 'react-bootstrap/Container';

export interface MenuItem {
  name: string;
  path: string;
}

interface NavigationProps {
  title: string;
  menu: MenuItem[];
}

const Navigation: FunctionComponent<NavigationProps> = ({ title, menu }) => (
  <Navbar bg="light" expand="sm">
    <Container>
      <Navbar.Brand href="/" className="px-3">
        {title}
      </Navbar.Brand>
      <Navbar.Toggle aria-controls="navbar-nav" />
      <Navbar.Collapse id="navbar-nav">
        <Nav className="me-auto">
          {menu.map((e) => {
            return <Nav.Link key={e.name} href={e.path}>{e.name}</Nav.Link>
          })}
        </Nav>
      </Navbar.Collapse>
    </Container>
  </Navbar>
);

export default Navigation;