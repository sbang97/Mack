import React, { Component } from "react";
import { Link } from "react-router";
import classnames from "classnames";

export default class NavLink extends Component {
	render() {
		return <Link className={classnames("NavLink", this.props.className)} {...this.props} activeClassName="active"/>;
	}
}
