import React, { Component } from 'react'
import {observable, action} from 'mobx'
import { observer } from 'mobx-react'
import Header from '../../components/Header'
import Messages from '../../components/Messages'
import Input from '../../components/Input'
import './index.css';

@observer
class Chat extends Component {
    @observable _messages
    @observable _socket
    constructor(props) {
        super(props)
        this._token = ""
        this._messages = []
    }

    @action
    componentWillMount() {
        this._socket = new WebSocket("ws://localhost:4000/api/v1/rooms/" + this.props.match.params.room +  "/messages/ws", this._token);
        this._socket.onmessage = this.onMessage.bind(this)
    }

    render() {
        return (
            <div className="Chat convex">
                <Header room={this.props.match.params.room}/>
                <Messages messages={this._messages} />
                <Input />
            </div>
        )
    }

    @action
    onMessage(ev) {
        let message = JSON.parse(ev.data)[0]
        console.log(message.sender.username)
        this._messages.push(message)
    }
}

export default Chat