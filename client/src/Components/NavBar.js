import React from 'react';
import PropTypes from 'prop-types';
import {
  Navbar,
  NavbarBrand,
  Nav,
  NavItem,
  NavLink
} from 'reactstrap';


class NavBar extends React.Component {
  constructor(props){
    super(props);
    this.state = {active: null};
  }

  render () {
    return (
      <div>
        <Navbar color="dark" dark expand="md">
          <NavbarBrand href="#">Sublime-CTF-Platform</NavbarBrand>
          <Nav className="ml-auto" navbar>
            <NavItem>
              <NavLink href="#">About</NavLink>
            </NavItem>
          </Nav>
        </Navbar>
      </div>
    );
  }
}

NavBar.propTypes = {
  path: PropTypes.string.isRequired
}

export default NavBar;
