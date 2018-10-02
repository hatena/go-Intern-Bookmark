import "./spa.scss";

import React from "react";
import ReactDOM from "react-dom";

import {App} from "./components/app";

const container = document.getElementById('container') as HTMLDivElement;

ReactDOM.render(<App />, container);
