import React, {Component} from "react";
import ReactDOM from "react-dom";
import {Router, browserHistory} from "react-router"
import { Col, Button, FormGroup, FormControl, ControlLabel, Form} from 'react-bootstrap';
import NavLink from "../Layout/NavLink";
import "./style.css";

export default class Home extends Component {

	handleSubmit(event) {
		event.preventDefault();
		fetch('https://api.sbang9.me/v1/sessions', {
			method :'post',
			mode: 'cors',
			contentType: 'application/json',
			body: JSON.stringify({
				email: 	ReactDOM.findDOMNode(this.refs.email).value,
				password: ReactDOM.findDOMNode(this.refs.password).value
			})
		})
		.then(function(resp) {
			var header = resp.headers.get("Authorization");
			if (header !== null) {
				localStorage.setItem("authorization", header);
				browserHistory.push('/common');
			} else {
				alert("Incorrect Email or Password");
			}
		})
		.catch(function(error) {
			alert(error);
		})
	}

	render() {
		return (
			<div className ="sign-in-div">
				<Form className="sign-in" horizontal onSubmit={event => this.handleSubmit(event)}>
					<h1 className="sign-in-title">Sign In</h1>
					<FormGroup className="input-field" controlId="formHorizontalEmail">
					<Col componentClass={ControlLabel} sm={2}>
						<h2 className="form-header">Email:</h2>
					</Col>
					<Col sm={12}>
						<FormControl className="email" type="email" placeholder="Email" ref="email"/>
					</Col>
					</FormGroup>
					<FormGroup controlId="formHorizontalPassword">
					<Col componentClass={ControlLabel} sm={2}>
						<h2 className="form-header">Password:</h2>
					</Col>
					<Col sm={12}>	
						<FormControl className="password" type="password" placeholder="Password"  ref="password"/>
					</Col>
					</FormGroup>
					<FormGroup>
					<Col smOffset={2} sm={12}>
						<Button className="submit" type="submit">
							Sign In
						</Button>
					</Col>
					</FormGroup>
					<p className="signup">Don't have an account yet? <NavLink to="/signup">Sign Up</NavLink></p>
				</Form>
				<nav>
				</nav>
			</div>
		);
		
	}
}
