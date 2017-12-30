import React from 'react';
import ReactDOM from 'react-dom';
import promiseFinally from 'promise.prototype.finally';
import { useStrict } from 'mobx';
import { Provider } from 'mobx-react';

import { BrowserRouter, Switch, Route } from 'react-router-dom'
import './index.css';
import Chat from './views/Chat';
import registerServiceWorker from './registerServiceWorker';

import authStore from './stores/authStore';
import commonStore from './stores/commonStore';

promiseFinally.shim();
useStrict(true);

const stores = {
    authStore,
    commonStore
}

ReactDOM.render(
    <Provider {...stores}>
        <BrowserRouter>
            <Switch>
            {/* <Route exact path="/" component={Home} /> */}
            <Route path="/at/:room" component={Chat} />
            </Switch>
        </BrowserRouter>
    </Provider>,
    document.getElementById('root'));
registerServiceWorker();