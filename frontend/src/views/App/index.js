import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Switch, Route } from 'react-router-dom'

import Landing from '../Landing'
import Login from '../Login'
import Register from '../Register'

import VerifyLanding from '../VerifyLanding'
import VerifyAccount from '../VerifyAccount'

import Chat from '../Chat'
import Home from '../Home'

import LockRoute from '../../components/Lock/HOC'
import Lock from '../../components/Lock/Semantic'

@inject('commonStore')
@observer
class App extends Component {
    render() {
        if (this.props.commonStore.token && this.props.commonStore.user) {
            return (
                [
                    <Lock minRole="unverifiedUser" user={this.props.commonStore.user}>
                        <Switch>
                            <Route path="/verify/account/:code" component={VerifyAccount} />
                            <Route component={VerifyLanding}/>
                        </Switch>
                    </Lock>,
                    <Lock minRole="normalUser" user={this.props.commonStore.user}>
                        <Switch>
                            <Route exact path="/" component={LockRoute(Home, "normalUser", this.props.commonStore.user)}/>
                            <Route path="/at/:room" component={LockRoute(Chat, "normalUser", this.props.commonStore.user)}/>
                            <Route component={LockRoute(Home, "normalUser", this.props.commonStore.user)}/>
                        </Switch>
                    </Lock>
                ]
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