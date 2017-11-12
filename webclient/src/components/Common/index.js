import React, { Component } from "react";
import { Dropdown, Form, Col, Glyphicon, Button, FormControl, Nav, NavItem, MenuItem, Navbar, FormGroup} from "react-bootstrap";
import NavLink from "../Layout/NavLink";
import ReactDOM from "react-dom";
import {Router, browserHistory} from "react-router";
import "./style.css";
const BASE_URL = "https://api.sbang9.me/v1/";
const generalChID = "59151b9b6b2d600001f328e8";
export default class Common extends Component {
    constructor(props) {
        super(props);
        this.state = {username: "username"}
    }

    componentDidMount() {
        var websock = new WebSocket("wss://api.sbang9.me/v1/websocket?auth="+localStorage.getItem("authorization"));
        websock.onmessage = (event) => {
            var evt = JSON.parse(event.data);
            switch(evt.type) {
            case "new message":
                var msgs = this.state.messages;
                msgs.push(evt.message);
                this.setState({
                    messages:msgs
                });
                break;
            case "new channel":
                var channels = this.state.channels;
                channels.push(evt.channel);
                this.setState({
                    channels: channels
                });
                break;
            case "new user":
                var users = this.state.users;
                users.push(evt.user);
                this.setState({
                    users: users
                });
                break;
            }
        };
    }

    componentWillMount() {
        var authHeader = localStorage.getItem("authorization");
        if (authHeader == "null" || (authHeader != null && authHeader <= 0)) {
            browserHistory.push("/")
            return;
        }
        fetch(BASE_URL+ "users/me", {
            method:'get',
            mode:'cors',
            contentType: 'application/json',
            headers: {
                'Authorization': authHeader,
            }
        })
        .then(response => {
            if (response.ok) {
                return response.json()
            } else {
                browserHistory.push("/");
                 return Promise.reject().then(() => response.text())
            }
        })
        .then(data => {
            this.setState({
                email: data["email"],
                username: data["userName"],
                firstname: data["firstName"],
                lastname: data["lastName"]
            })
        })
       
        fetch(BASE_URL + "channels",{
            method:'get',
            contentType: 'application/json',
            headers: {
                'Authorization': authHeader,
            }
        }).then(response => {
            if (response.ok) {
                return response.json()
            }
            return Promise.reject().then(() => response.text())
        }).then(data => {
            this.setState ({
                channels: data
            })
        })
        fetch(BASE_URL+"channels/" + generalChID, {
            method:`get`,
            contentType: 'application/json',
            headers: {
                'Authorization': authHeader,
            }
        }).then(response => { 
            if(response.ok) {
                return response.json()
            }
            return Promise.reject().then(() => response.text())
        }).then(data => {
            this.setState({
                currentChannel: "General",
                currentID: generalChID,
                messages: data
            })
        })
        .catch(err => {
            alert(err);
        }) 
    }

    handleUpdate() {
        browserHistory.push("/profile");
    }
    handleNewChannel() {
        browserHistory.push("/channel");
    }

    handleSignOut() {
        fetch(BASE_URL+"sessions/mine", {
            method:'delete',
            mode:'cors',
            headers: {
                'Authorization': localStorage.getItem("authorization"),
            }
        })
        .then(response => {
            if (response.ok) {
                localStorage.removeItem("authorization");
                browserHistory.push("/");
            }
        }).catch(error => {
            alert(error);
        })
    }

    handleChannelChange(event, chID) {
        fetch(BASE_URL+"channels/"+chID, {
            method:`get`,
            headers: {
                'Authorization': localStorage.getItem("authorization"),
            }
        }).then(response => {
            return response.json()
        }).then(data => {
            this.setState({
                currentID: chID,
                messages: data
            })
        }).catch(error =>{
            alert(error);
        })
    }

    handleMessage(event) {
        event.preventDefault();
        if (ReactDOM.findDOMNode(this.refs.msg).value === "") {
            alert("message field cannot be left empty");
        } else {
            fetch(BASE_URL+"messages", {
            method:`post`,
            headers: {
                'Authorization': localStorage.getItem("authorization"),
            },
            body: JSON.stringify({
                body: ReactDOM.findDOMNode(this.refs.msg).value,
                channelid: this.state.currentID
            })
        })
        }
        ReactDOM.findDOMNode(this.refs.msg).value = "";
    }

    render() {
        var userName;
        var firstName;
        var lastName;
        var channels;
        var messages;
        var currentChannel;
        var users;
        userName = this.state.username;
        firstName = this.state.firstname;
        lastName = this.state.lastname;
        if (this.state.channels) {
            channels = this.state.channels.map(c => 
                {
                    return (
                        <NavItem key={c.id} onClick={event=> this.handleChannelChange(event, c.id)}>
                    <h3>{c.name}</h3>
                </NavItem>
                    )
                }
            );
        }
        if (this.state.users) {
            users = this.state.users.map(u => 
                <li key={u.id}>{u.username}</li>
            );
        }
        if (this.state.messages) {
            messages = this.state.messages.map(m =>
                <div className="message" key={m.id}>
                    <img className="userphoto" src={m.photoURL} alt=""/>
                    <span className="username">
                        {m.username}
                    </span>
                    <span className="created">
                        {m.createdAt}
                    </span>
                    <p className="messagebody">
                        {m.body}
                    </p>
                </div>
            );
        }
        return (
            <div>
            <div className="side-nav">
                <Navbar>
                    <Navbar.Header>
                        <Nav bsStyle="pills" stacked>
                            <h1>Channels</h1>
                            {channels}
                        </Nav>
                    </Navbar.Header>
                    <div>
                        <p>{userName}</p>
                        <p>{firstName} {lastName}</p>
                        <Button onClick={this.handleUpdate}>Edit Profile</Button>
                        <Button onClick={this.handleSignOut}>Sign Out</Button>
                        <Button onClick={this.handleNewChannel}>Create a New Channel</Button>
                    </div>
                </Navbar>
            </div>
                <div className="messages">
                {messages}
                 <Form className="msg-form" horizontal onSubmit={event => this.handleMessage(event)}>
                    <FormGroup>
                        <Col sm={12}>	
                            <FormControl id="msg" className="msg" type="text" ref="msg" placeholder="type a message..."/>
                        </Col>
                        <Button className="send" type="submit">Send</Button>
                    </FormGroup>
                </Form>
                </div>
            </div>
        )
    }
}