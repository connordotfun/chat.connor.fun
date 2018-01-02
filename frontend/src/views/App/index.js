import React, { Component } from 'react'
import { observer, inject } from 'mobx-react'
import { Switch, Route } from 'react-router-dom'

import Chat from '../Chat'
import Landing from '../Landing'

@inject('commonStore')
@observer
class App extends Component {
    render() {
        if (this.props.commonStore.token) {
            return (
                <Switch>
                    <Route path="/at/:room" component={Chat}/>
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