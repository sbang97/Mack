import React, { Component } from "react";
import {Button, FormControl, Col, Form, FormGroup, ControlLabel, Checkbox} from "react-bootstrap";
import NavLink from "../Layout/NavLink";
import ReactDOM from "react-dom";
import { browserHistory} from "react-router"
import "./style.css";


export default class Channel extends Component {
    constructor(props) {
        super(props);
        this.state = {name: "name", description:"description", private:false}
    }

    handleCheck(event) {
        if (this.state.private === false) { 
            this.setState({
                private: true
            })
        } else {
            this.setState({
                private: false
            })
        }
    }

    handleCancel() {
        browserHistory.push("/common");
    }

    handleSubmit(event) {
        event.preventDefault();
        console.log(this.state.private)
        if (this.refs.firstname === "" || this.refs.firstname === "") alert("name cannot be empty");
        else {
            fetch('https://api.sbang9.me/v1/channels', {
                method:'post',
                mode:'cors',
                contentType:'application/json',
                headers: {
                    'Authorization': localStorage.getItem("authorization"),
                },
                body: JSON.stringify({
                    name: ReactDOM.findDOMNode(this.refs.name).value,
                    description: ReactDOM.findDOMNode(this.refs.description).value,
                    private: this.state.private,
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

    render() {
        return (
        <div className="channel-div">
            <Form className="channel" horizontal onSubmit={event => this.handleSubmit(event)}>
                <FormGroup className="input-field" controlId="formHorizontalText">
                    <Col componentClass={ControlLabel} sm={2}>
                        <h4 className="form-header">Enter a name for the channel:</h4>
                    </Col>
                    <Col sm={12}>
                        <FormControl className="name" type="text" ref="name"/>
                    </Col>
                </FormGroup>
                <FormGroup className="input-field" controlId="formHorizontalText">
                    <Col componentClass={ControlLabel} sm={2}>
                        <h4 className="form-header">Enter a description for the channel:</h4>
                    </Col>
                    <Col sm={12}>
                        <FormControl className="description" type="text" ref="description"/>
                    </Col>
                    <Checkbox onChange={event => this.handleCheck(event)}>
                       Private
                    </Checkbox>
                </FormGroup>
                <Button onClick={this.handleCancel}>Cancel</Button>
                <Button className="submit" type="submit">Create</Button>
            </Form>
        </div>
        )
    }
}