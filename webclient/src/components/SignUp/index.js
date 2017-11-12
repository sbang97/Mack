import React, { Component } from "react";
import ReactDOM from "react-dom";
import {Router, browserHistory} from "react-router";
import {Form, Col, Button, FormGroup, FormControl, ControlLabel} from "react-bootstrap";
import NavLink from "../Layout/NavLink";
import "./style.css";

export default class About extends Component {

	handleSubmit(event) {
		event.preventDefault();
		if (ReactDOM.findDOMNode(this.refs.password).value === "" || ReactDOM.findDOMNode(this.refs.conf).value === "" ||
			ReactDOM.findDOMNode(this.refs.username).value === "" || ReactDOM.findDOMNode(this.refs.firstname).value === "" 
			|| ReactDOM.findDOMNode(this.refs.lastname).value === "" ) alert("fields cannot be left empty");
		else if (ReactDOM.findDOMNode(this.refs.password).value.length < 6) alert("Password must be at least 6 characters long");
		else if (ReactDOM.findDOMNode(this.refs.password).value!== ReactDOM.findDOMNode(this.refs.conf).value) alert("Passwords do not match");
		else {
			fetch('https://api.sbang9.me/v1/users', {
				method :'post',
				mode: 'cors',
				contentType: 'application/json',
				body: JSON.stringify({
					email: ReactDOM.findDOMNode(this.refs.email).value,
					password: ReactDOM.findDOMNode(this.refs.password).value,
					passwordConf: ReactDOM.findDOMNode(this.refs.conf).value,
					username: ReactDOM.findDOMNode(this.refs.username).value,
					firstname: ReactDOM.findDOMNode(this.refs.firstname).value,
					lastname: ReactDOM.findDOMNode(this.refs.lastname).value
				})
			}).then(function(resp) {
				if (resp.ok) {
					browserHistory.push('/');
				} else {
					alert("user already exists for that username or email");
				}
			}).catch(function(error) {
				alert(error);
			});
		}
	}

	render() {
		return (
			<div className ="sign-up-div">
				<Form className="sign-up" horizontal onSubmit={event => this.handleSubmit(event)}>
					<h1 className="sign-up-title">Sign Up</h1>
					<FormGroup className="input-field" controlId="formHorizontalEmail">
						<Col componentClass={ControlLabel} sm={2}>
							<h3 className="form-header">Email:</h3>
						</Col>
						<Col sm={12}>
							<FormControl className="email" type="email" placeholder="Email" ref="email" />
						</Col>
					</FormGroup>
					<FormGroup validationState='success' controlId="formHorizontalPassword">
						<Col componentClass={ControlLabel} sm={2}>
							<h3 className="form-header">Password:</h3>
						</Col>
						<Col sm={12}>	
							<FormControl className="password" type="password" placeholder="Password" ref="password"/>
						</Col>
					</FormGroup>
					<FormGroup controlId="formBasicText">
						<Col componentClass={ControlLabel} sm={2}>
							<h3 className="form-header">Password Confirmation:</h3>
						</Col>
						<Col smOffset={2} sm={12}>
							<FormControl type='password' className='conf' placeholder='Confirm Password' ref="conf"/>
						</Col>
					</FormGroup>
					<FormGroup controlId="formBasicText">
						<Col componentClass={ControlLabel} sm={2}>
							<h3 className="form-header">Username:</h3>
						</Col>
						<Col smOffset={2} sm={12}>
							<FormControl type='username' className='username' placeholder='Username' ref="username"/>
						</Col>
					</FormGroup>
					<FormGroup controlId="formBasicText">
						<Col componentClass={ControlLabel} sm={2}>
							<h3 className="form-header">First Name:</h3>
						</Col>
						<Col smOffset={2} sm={12}>
							<FormControl type='username' className='firstname' placeholder='First Name' ref="firstname"/>
						</Col>
					</FormGroup>
					<FormGroup controlId="formBasicText">
						<Col componentClass={ControlLabel} sm={2}>
							<h3 className="form-header">Last Name:</h3>
						</Col>
						<Col smOffset={2} sm={12}>
							<FormControl type='username' className='lastname' placeholder='Last Name' ref="lastname"/>
						</Col>
					</FormGroup>
					<FormGroup>
					<Col smOffset={2} sm={12}>
						<Button className="submit" type="submit">
							Register
						</Button>
					</Col>
					</FormGroup>
					<p className="signup">Already have an account? <NavLink to="/">Sign In</NavLink></p>
				</Form>
				<nav>
				</nav>
			</div>
		);
	}
}
