import React, { Component } from 'react'
import {observable, action} from 'mobx'
import { observer, inject } from 'mobx-react'
import Header from '../../components/Header'
import Messages from '../../components/Messages'
import Input from '../../components/ChatInput'
import './index.css'


@inject('commonStore')
@inject('socketStore')
@inject('authStore')
@observer
class Chat extends Component {
    @observable _messages = []

    @action
    componentWillMount() {
        this.props.socketStore.joinRoom(this.props.match.params.room)
        this.props.socketStore.addListener(this.onMessage.bind(this))
    }

    @action
    componentWillUnmount() {
        this.props.socketStore.resetListeners()
        this.props.socketStore.leaveRoom()
    }

    render() {
        return (
            <div className="Chat convex">
                <Header room={this.props.match.params.room} handleLeave={this.props.authStore.logout}/>
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