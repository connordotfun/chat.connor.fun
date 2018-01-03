import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Switch, Route } from 'react-router-dom'

import Chat from '../Chat'
import Landing from '../Landing'
import Home from '../Home'

@inject('commonStore')
@observer
class App extends Component {
    render() {
        if (this.props.commonStore.token) {
            return (
                <Switch>
                    <Route exact path="/" component={Home}/>
                    <Route path="/at/:room" component={Chat}/>
                    <Route component={Home}/>
                </Switch>
            )
        } else {
            return (
                <Switch>
                    <Route exact path="/" component={Landing}/>
                    <Route component={Landing} />
                </Switch>
            )
        }
    }
}

export default App