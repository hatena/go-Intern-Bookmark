import React from "react";
import {NavLink} from "react-router-dom";

export const GlobalHeader = () => (
  <header className="GlobalHeader">
    <h1>Bookmark</h1>
    <nav>
      <ul>
        <li><NavLink to="/">Top</NavLink></li>
        <li><NavLink to="/me">Me</NavLink></li>
      </ul>
    </nav>
  </header>
);