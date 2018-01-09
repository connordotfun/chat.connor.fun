import React, { Component } from 'react'
import {action} from 'mobx'
import { observer, inject } from 'mobx-react'
import Header from '../../components/Header'
import Messages from '../../components/Messages'
import Input from '../../components/ChatInput'
import './index.css'


@inject('commonStore', 'socketStore', 'authStore', 'roomStore')
@observer
class Chat extends Component {
    @action
    componentWillMount() {
        this.props.socketStore.joinRoom(this.props.match.params.room)
    }

    @action
    componentWillUnmount() {
        this.props.socketStore.resetListeners()
        this.props.socketStore.leaveRoom()
    }

    render() {
        return (
            <div className="Chat convex">
                <Header room={this.props.match.params.room} handleLeave={this.onLeave.bind(this)} handleExit={this.props.authStore.logout} />
                <Messages messages={this.props.roomStore.messages} />
                <Input />
            </div>
        )
    }

    @action
    onLeave() {
        this.props.socketStore.leaveRoom()
        this.props.history.push('/')
    }
}

export default Chat