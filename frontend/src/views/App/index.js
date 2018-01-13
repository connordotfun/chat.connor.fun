import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Switch, Route } from 'react-router-dom'

import Landing from '../Landing'
import Login from '../Login'
import Register from '../Register'

import Chat from '../Chat'
import Home from '../Home'

import Lock from '../../components/Lock/HOC'

@inject('commonStore')
@observer
class App extends Component {
    render() {
        if (this.props.commonStore.token) {
            return (
                <Switch>
                    <Route exact path="/" component={Lock(Home, "normalUser", this.props.commonStore.user)}/>
                    <Route path="/at/:room" component={Lock(Chat, "normalUser", this.props.commonStore.user)}/>
                    <Route component={Lock(Home, "normalUser", this.props.commonStore.user)}/>
                </Switch>
            )
        } else {
            return (
                <Switch>
                    <Route exact path="/" component={Landing} />
                    <Route exact path="/login" component={Login} />
                    <Route exact path="/register" component={Register} />
                    <Route component={Landing} />
                </Switch>
            )
        }
    }
}

export default App