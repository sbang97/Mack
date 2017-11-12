import React, { Component } from "react";
import classnames from "classnames";

import { IndexLink } from "react-router";
import NavLink from "../NavLink";
import "./style.css"

export default class Header extends Component {
	render() {
		return(
			<header className={classnames("Header", this.props.className)}>
				<h1 className="title">Mack</h1>
			</header>
		);
	}
}

/*
			
*/
