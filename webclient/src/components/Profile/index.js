import React, { Component } from "react";
import {Button, FormControl, Col, Form, FormGroup, ControlLabel} from "react-bootstrap";
import NavLink from "../Layout/NavLink";
import ReactDOM from "react-dom";
import { browserHistory} from "react-router"
import "./style.css";

export default class Profile extends Component {
    constructor(props) {
        super(props);
        this.state = {firstname: "firstname", lastname:"lastname", photoURL:""}
    }

    componentWillMount() {
        var authHeader = localStorage.getItem("authorization");
        if (authHeader == "null" || (authHeader != null && authHeader <= 0) || authHeader == undefined) {
            browserHistory.push("/")
            return;
        }
        fetch("https://api.sbang9.me/v1/users/me", {
            method:'get',
            mode:'cors',
            contentType: 'application/json',
            headers: {
                'Authorization': localStorage.getItem("authorization"),
            }
        })
        .then(response => {
            if (response.ok) {
                return response.json()
            }
            if (response.statusCode === 403) {
                browserHistory.push("/");
                return;
            }
            return Promise.reject().then(() => response.text())
        })
        .then(data => {
            this.setState({
                username: data["userName"],
                firstname: data["firstName"],
                lastname: data["lastName"],
                photoURL: data["photoURL"]
            })
        })
        .catch(error => {
            alert(error);
        })
    }

    handleReturn() {
        browserHistory.push("/common");
    }


    handleSubmit(event) {
        event.preventDefault();
        if (this.refs.firstname === "" || this.refs.firstname === "") alert("name cannot be empty");
        else {
            fetch('https://api.sbang9.me/v1/users/me', {
                method:'PATCH',
                mode:'cors',
                contentType:'application/json',
                headers: {
                    'Authorization': localStorage.getItem("authorization"),
                },
                body: JSON.stringify({
                    firstName: ReactDOM.findDOMNode(this.refs.firstname).value,
                    lastName: ReactDOM.findDOMNode(this.refs.lastname).value
                })
            }).then(resp => {
                if (resp.ok) {
                    browserHistory.push("/common");
                }
            }).catch(error => {
                alert(error);
            })
        }
    }

    firstNameHandler() {
        return this.state.firstname;
    }
    render() {
        var userName;
        var firstName;
        var lastName;
        var imgURL;    
        if (this.state.username) {
            userName = this.state.username;
        }
        if (this.state.firstname) {
            firstName = this.state.firstname;
        }
        if (this.state.lastname) {
            lastName = this.state.lastname;
        }
        if(this.state.photoURL) {
            imgURL = this.state.photoURL;
        }
        return (
            <div className="profile-div">
                <Form className="profile" horizontal onSubmit={event => this.handleSubmit(event)}>
                        <div>
                            <img src={imgURL} alt="profile photo"/>
                            <h1>{userName}</h1>
                            <h3>{this.state.firstname} {lastName}</h3>
                        </div>
                        <FormGroup className="input-field" controlId="formHorizontalText">
                        <Col componentClass={ControlLabel} sm={2}>
                            <h4 className="form-header">Enter a new First Name:</h4>
                        </Col>
                        <Col sm={12}>
                            <FormControl className="firstname" type="text" ref="firstname"/>
                        </Col>
                        </FormGroup>
                        <FormGroup controlId="formHorizontalText">
                        <Col componentClass={ControlLabel} sm={2}>
                            <h4 className="form-header">Enter a new Last Name:</h4>
                        </Col>
                        <Col sm={12}>	
                            <FormControl className="lastname" type="text" ref="lastname"/>
                        </Col>
                        </FormGroup>
                        <Button onClick={this.handleReturn}>Cancel</Button>
                        <Button className="submit" type="submit">Update</Button>
                    </Form>
            </div>
        )
    }
}