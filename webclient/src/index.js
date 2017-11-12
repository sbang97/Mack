import React from "react";
import ReactDOM from "react-dom";
import { Router, browserHistory } from "react-router";
import Routes from "./routes";

import "./index.css";

const outlet = document.getElementById("root");
ReactDOM.render(<Router history={browserHistory} routes={Routes} />, outlet);
