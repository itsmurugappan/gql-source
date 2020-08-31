import React, { useContext } from 'react';
import './Foooddiieee.css';
import { AppContext } from './AppContext'


function TopNav() {
  const [isManager, setIsManager] = useContext(AppContext); 

  const toggleManagerAccess = e => {
    setIsManager(prevValue => !prevValue)
  };
  return (
    <nav>
      <ul className="nav-flex-row">
        <li className="nav-item">
          <a href="#about">About</a>
        </li>
        <li className="nav-item">
          <a href="#reservation">Reservation</a>
        </li>
        <li className="nav-item">
          <a href="#menu">Menu</a>
        </li>
        <li className="nav-item">
          <a href="#shop">Shop</a>
        </li>
        <li className="nav-item">
          <a href="#" onClick={toggleManagerAccess}>Manager Access</a>
        </li>        
      </ul>
    </nav>
  );
}

export default TopNav;
